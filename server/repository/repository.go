package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"massive-message/server/sdk"
	"os"
)

var (
	connection *gorm.DB
	tables     = []TableInfo{
		{Name: "Server", Info: new(Server)},
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
		log.Info("[Server-Repository] Init DB connection.")
		args := fmt.Sprintf("host=postgres port=5432 user=postgres dbname=server sslmode=disable password=iforgot")
		connection, err = gorm.Open("postgres", args)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("[Server-Repository] DB open failed.")
			return err
		}
		// connection.LogMode(true)
		connection.SingularTable(true)
	} else {
		log.Info("[Server-Repository] DB connection exist.")
	}

	v, found := os.LookupEnv("build_mock_servers")
	if found && v == "yes" {
		DropTablesIfExist()
		CreateTables()
		log.Info("[Server-Repository] Prepare mock servers start.")
		prepareMockServers()
		log.Info("[Server-Repository] Prepare mock servers done.")
	}
	return nil
}

// CreateTables creates the tables.
func CreateTables() error {
	for _, v := range tables {
		if err := connection.CreateTable(v.Info).Error; err != nil {
			log.WithFields(log.Fields{"Table": v.Name, "error": err}).Error("[Server-Repository] create table failed.")
			return err
		}
	}
	return nil
}

// DropTablesIfExist drops tables if they are exist.
func DropTablesIfExist() error {
	for _, v := range tables {
		if err := connection.DropTableIfExists(v.Info).Error; err != nil {
			log.WithFields(log.Fields{"Table": v.Name, "error": err}).Error("[Server-Repository] remove table failed.")
			return err
		}
	}
	return nil
}

// SaveServer saves the notification into the database.
func SaveServer(entity *Server) error {
	if err := connection.Create(entity).Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("[Server-Repository] Save server failed.")
		return err
	}
	return nil
}

// UpdateServerHealth update server's health.
func UpdateServerHealth(id string, warnings, criticals int) error {
	var entity Server
	entity.ID = id
	if err := connection.Model(&entity).Where("id = ?", id).Updates(
		map[string]interface{}{
			"warnings":  warnings,
			"criticals": criticals,
		}).Error; err != nil {
		log.WithFields(log.Fields{"id": id, "error": err}).Error("[Server-Repository] Update server health failed.")
		return err
	}
	return nil
}

// GetServers return a server collection specified by start, count and orderby.
func GetServers(start, count int64, orderby string) (*sdk.ServerCollection, error) {
	var err error

	servers := []Server{}
	if orderby == "Name" {
		err = connection.Limit(count).Offset(start).Order("name").Find(&servers).Error
	} else {
		err = connection.Limit(count).Offset(start).Order("criticals").Order("warnings").Find(&servers).Error
	}

	if err != nil {
		log.WithFields(log.Fields{"start": start, "count": count, "orderby": orderby, "error": err}).Error("[Server-Repository] Get servers failed.")
		return nil, err
	}
	ret := sdk.ServerCollection{}
	for _, server := range servers {
		ret.Member = append(ret.Member, sdk.Server{
			ID:           server.ID,
			URL:          server.URL,
			Name:         server.Name,
			SerialNumber: server.SerialNumber,
			Warnings:     server.Warnings,
			Criticals:    server.Criticals,
		})
	}
	return &ret, nil
}

// GetServerByID gets server by ID, if server does not exsit return nil, nil.
func GetServerByID(id string) (*sdk.Server, error) {
	var (
		record Server
	)

	// Here we try to find the SSG before preload all the values, hoping this would be faster.
	notFound := connection.Where("id = ?", id).First(&record).RecordNotFound()
	if err := connection.Error; err != nil {
		log.WithFields(log.Fields{"id": id}).Warn("[Server-Repository] Get server by ID failed, find first record failed.")
		return nil, err
	}
	if notFound {
		log.WithFields(log.Fields{"id": id}).Warn("[Server-Repository] Get server by ID failed, server does not exist.")
		return nil, nil
	}
	ret := sdk.Server{
		ID:           record.ID,
		URL:          record.URL,
		Name:         record.Name,
		SerialNumber: record.SerialNumber,
		Warnings:     record.Warnings,
		Criticals:    record.Criticals,
	}
	return &ret, nil
}

func prepareMockServers() {
	for i := 0; i < 10000; i++ {
		sn := fmt.Sprintf("sn-huawei-%d", i)
		server := Server{}
		server.Name = fmt.Sprintf("Huawei Server %05d", i)
		server.ID = sn
		server.URL = "/api/v1/servers/" + sn
		server.SerialNumber = sn
		SaveServer(&server)
	}
	for i := 0; i < 10000; i++ {
		sn := fmt.Sprintf("sn-hpe-%d", i)
		server := Server{}
		server.Name = fmt.Sprintf("HPE Server %05d", i)
		server.ID = sn
		server.URL = "/api/v1/servers/" + sn
		server.SerialNumber = sn
		SaveServer(&server)
	}
	for i := 0; i < 10000; i++ {
		sn := fmt.Sprintf("sn-dell-%d", i)
		server := Server{}
		server.Name = fmt.Sprintf("Dell Server %05d", i)
		server.ID = sn
		server.URL = "/api/v1/servers/" + sn
		server.SerialNumber = sn
		SaveServer(&server)
	}
	for i := 0; i < 10000; i++ {
		sn := fmt.Sprintf("sn-ibm-%d", i)
		server := Server{}
		server.Name = fmt.Sprintf("IBM Server %05d", i)
		server.ID = sn
		server.URL = "/api/v1/servers/" + sn
		server.SerialNumber = sn
		SaveServer(&server)
	}
	for i := 0; i < 10000; i++ {
		sn := fmt.Sprintf("sn-lenovo-%d", i)
		server := Server{}
		server.Name = fmt.Sprintf("Lenovo Server %05d", i)
		server.ID = sn
		server.URL = "/api/v1/servers/" + sn
		server.SerialNumber = sn
		SaveServer(&server)
	}
}
