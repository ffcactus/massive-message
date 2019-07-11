package repository

import (
	"github.com/google/uuid"
	"massive-message/server/sdk"
)

// Server represents the server table in DB.
type Server struct {
	ID           string `gorm:"primary_key"`
	URL          string
	Name         string `gorm:"index"`
	SerialNumber string
	Warnings     int `gorm:"index"`
	Criticals    int `gorm:"index"`
}

func newServer(o *sdk.Server) *Server {
	ret := Server{}
	ret.ID = uuid.New().String()
	ret.URL = "/api/v1/servers/" + ret.ID
	ret.Name = o.Name
	ret.SerialNumber = o.SerialNumber
	return &ret
}
