package repository

import (
	"github.com/google/uuid"
	"massive-message/server/sdk"
)

// Server represents the server table in DB.
type Server struct {
	ID           string `gorm:"column:ID;primary_key"`
	URL          string `gorm:"column:URL"`
	Name         string `gorm:"column:Name"`
	SerialNumber string `gorm:"column:SerialNumber"`
	Warnings     int    `gorm:"column:Warnings"`
	Criticals    int    `gorm:"column:Criticals"`
}

// TableName will set the table name.
func (Server) TableName() string {
	return "Server"
}

func newServer(o *sdk.Server) *Server {
	ret := Server{}
	ret.ID = uuid.New().String()
	ret.URL = "/api/v1/servers/" + ret.ID
	ret.Name = o.Name
	ret.SerialNumber = o.SerialNumber
	return &ret
}
