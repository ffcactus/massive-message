package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"massive-message/notification/sdk"
	"time"
)

var (
	connection *gorm.DB
	tables     = []TableInfo{
		{Name: "Notification", Info: new(Notification)},
	}
)

// TableInfo The tables in DB.
type TableInfo struct {
	Name string
	Info interface{}
}

// Notification is the notification entity.
type Notification struct {
	ID          string    `gorm:"column:ID;primary_key"`
	CreatedAt   time.Time `gorm:"column:CreatedAt"`
	UpdatedAt   time.Time `gorm:"column:UpdatedAt"`
	Key         string
	VersusKey   string
	URL         string
	Type        string
	ReceivedAt  time.Time
	Severity    string
	Description string
}

// TableName will set the table name.
func (Notification) TableName() string {
	return "Notification"
}

func newEntity(o *sdk.Notification) *Notification {
	ret := Notification{}
	ret.ID = uuid.New().String()
	ret.Key = o.Key
	ret.URL = o.URL
	ret.Type = o.Type
	ret.ReceivedAt = o.ReceivedAt
	ret.Severity = o.Severity
	ret.Description = o.Description
	return &ret
}

// Init perform the initialization work.
func Init() error {
	if connection == nil {
		log.Info("[Event] Init DB connection.")
		args := fmt.Sprintf("host=postgres port=5432 user=postgres dbname=notification sslmode=disable password=iforgot")
		db, err := gorm.Open("postgres", args)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("[Event] DB open failed.")
			return err
		}
		// db.LogMode(true)
		db.SingularTable(true)
		connection = db
	} else {
		log.Info("[Event] DB connection exist.")
	}
	return nil
}

// CreateTables creates the tables.
func CreateTables() error {
	for i := range tables {
		if err := connection.CreateTable(tables[i].Info).Error; err != nil {
			log.WithFields(log.Fields{"Table": tables[i].Name, "error": err}).Error("[Event] create table failed.")
			return err
		}
	}
	return nil
}

// DropTablesIfExist drops tables if they are exist.
func DropTablesIfExist() error {
	for i := range tables {
		if err := connection.DropTableIfExists(tables[i].Info).Error; err != nil {
			log.WithFields(log.Fields{"Table": tables[i].Name, "error": err}).Error("[Event] remove table failed.")
			return err
		}
	}
	return nil
}

// ProcessNotification will first save this notification.
// If this notification is a alert, it will try to remove the versus notification.
// To find out the versus one, it will search by using the URL and Key.
//
func ProcessNotification(o *sdk.Notification) error {
	entity := newEntity(o)
	if err := connection.Create(entity).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("[Event] Save event failed.")
		return err
	}
	if entity.Type != "Alert" {
		return nil
	}
	record := Notification{}
	if connection.Where("\"URL\" = ? and \"VersusKey\" = ?", entity.URL, entity.VersusKey).First(&record).RecordNotFound() {
		return nil
	}

	return nil
}
