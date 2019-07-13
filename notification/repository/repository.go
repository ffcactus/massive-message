package repository

import (
	"container/list"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"massive-message/notification/sdk"
)

var (
	connection *gorm.DB
	tables     = []TableInfo{
		{Name: "event", Info: new(Event)},
		{Name: "alert", Info: new(Alert)},
		{Name: "url_health", Info: new(URLHealth)},
	}
)

// TableInfo The tables in DB.
type TableInfo struct {
	Name string
	Info interface{}
}

// Init perform the initialization work.
func Init() error {
	var err error
	if connection == nil {
		log.Info("[Notification-Repository] Init DB connection.")
		args := fmt.Sprintf("host=postgres port=5432 user=postgres dbname=notification sslmode=disable password=iforgot")
		connection, err = gorm.Open("postgres", args)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("[Notification-Repository] DB open failed.")
			return err
		}
		// connection.LogMode(true)
		connection.SingularTable(true)
	} else {
		log.Info("[Notification-Repository] DB connection exist.")
	}
	return nil
}

// CreateTables creates the tables.
func CreateTables() error {
	for _, v := range tables {
		if err := connection.CreateTable(v.Info).Error; err != nil {
			log.WithFields(log.Fields{"Table": v.Name, "error": err}).Error("[Notification-Repository] create table failed.")
			return err
		}
	}
	return nil
}

// DropTablesIfExist drops tables if they are exist.
func DropTablesIfExist() error {
	for _, v := range tables {
		if err := connection.DropTableIfExists(v.Info).Error; err != nil {
			log.WithFields(log.Fields{"Table": v.Name, "error": err}).Error("[Notification-Repository] remove table failed.")
			return err
		}
	}
	return nil
}

// SaveNotification saves the notification into the database.
func SaveNotification(o *sdk.Notification) error {
	var entity interface{}
	if o.Type == "Alert" {
		entity = newAlert(o)
	} else {
		entity = newEvent(o)
	}
	if err := connection.Create(entity).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("[Notification-Repository] Save notification failed.")
		return err
	}
	return nil
}

// GetTargetsHaveAlert returns all the targets that have alerts.
// On error, return nil.
func GetTargetsHaveAlert() ([]string, error) {
	var sqlResult []Alert
	if err := connection.Select("DISTINCT(url)").Find(&sqlResult).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("[Notification-Repository] Get targets that having alerts failed.")
		return nil, err
	}
	var ret []string
	for _, v := range sqlResult {
		ret = append(ret, v.URL)
	}
	return ret, nil
}

// CombineAlertsByURL finds all the alerts that matches the url, and remove the ones that can be removed.
// The ones that can be removed can be found like this:
// 1. Sorts the alerts by GeneratedAt.
// 2. Literates the sorted alerts, from the head,
func CombineAlertsByURL(url string) (*sdk.HealthChangeNotification, error) {
	var alerts []Alert
	if err := connection.Order("generated_at desc").Where("url = ?", url).Find(&alerts).Error; err != nil {
		log.WithFields(log.Fields{"url": url, "error": err}).Warn("[Notification-Repository] Combine alerts failed, get alerts by URL failed.")
		return nil, err
	}
	var removeFromDB []Alert
	// Save the record to the list for fast remove operation.
	l := list.New()
	for _, alert := range alerts {
		l.PushBack(alert)
	}

	// Pick up a element from the head. Literates to the end and remove the elements that the Key matahes key or Key matches VersusKey.
	for e := l.Front(); e != nil; e = e.Next() {
		alert := e.Value.(Alert)
		// We are not going to show any alerts with severity OK.
		if alert.Severity == "OK" {
			removeFromDB = append(removeFromDB, alert)
		}

		var removeFromList []*list.Element
		for r := e.Next(); r != nil; r = r.Next() {
			check := r.Value.(Alert)
			if check.Key == alert.VersusKey {
				removeFromList = append(removeFromList, r)
				removeFromDB = append(removeFromDB, check)
				continue
			}
			if check.Key == alert.Key {
				removeFromList = append(removeFromList, r)
				removeFromDB = append(removeFromDB, check)
				continue
			}
		}
		for _, toRemove := range removeFromList {
			l.Remove(toRemove)
		}
	}

	// Remove the records from DB.
	for _, toRemove := range removeFromDB {
		// Ignore errors here.
		// Errors may raise here, for example, another routine is doing the same work.
		// However, it seems OK. (Someone please help me to prove it)
		log.WithFields(log.Fields{"url": toRemove.URL, "key": toRemove.Key}).Info("[Notification-Repository] Remove de-active alert.")
		connection.Unscoped().Delete(&toRemove)
	}
	notification := sdk.HealthChangeNotification{}
	notification.URL = url
	for e := l.Front(); e != nil; e = e.Next() {
		alert := e.Value.(Alert)
		if alert.Severity == "Warning" {
			notification.Warnings++
		} else if alert.Severity == "Critical" {
			notification.Criticals++
		}
	}
	return &notification, nil
}

// CheckAndUpdateURLHealth checks if URL's health state changed, if yes it save the now health state.
// Note that for any error, we take it as should not update.
func CheckAndUpdateURLHealth(notification *sdk.HealthChangeNotification) (bool, error) {
	record := URLHealth{
		URL:       notification.URL,
		Warnings:  notification.Warnings,
		Criticals: notification.Criticals,
	}
	notFound := connection.Where("url = ?", notification.URL).First(&record).RecordNotFound()
	if err := connection.Error; err != nil {
		log.WithFields(log.Fields{"url": notification.URL, "error": err}).Warn("[Notification-Repository] Update URL health failed, get health failed.")
		return false, err
	}
	if notFound {
		if err := connection.Save(&record).Error; err != nil {
			log.WithFields(log.Fields{"url": record.URL, "Warnings": record.Warnings, "Criticals": record.Criticals, "error": err}).Info("[Notification-Repository] Update URL health failed, create health record failed.")
		}
		log.WithFields(log.Fields{"url": record.URL, "Warnings": record.Warnings, "Criticals": record.Criticals}).Info("[Notification-Repository] Found a new URL, create health state.")
		return true, nil
	}
	if record.Warnings == notification.Warnings && record.Criticals == notification.Criticals {
		return false, nil
	}

	log.WithFields(log.Fields{"url": notification.URL, "fromWarnings": record.Warnings, "toWarningss": notification.Warnings, "fromCriticals": record.Criticals, "toCriticals": notification.Criticals}).Info("[Notification-Repository] Update URL health.")
	record.Warnings = notification.Warnings
	record.Criticals = notification.Criticals

	if err := connection.Save(&record).Error; err != nil {
		log.WithFields(log.Fields{"url": notification.URL, "warnings": notification.Warnings, "criticals": notification.Criticals, "error": err}).Warn("[Notification-Repository] Update URL health failed, update URL health failed.")
		return false, err
	}
	return true, nil
}
