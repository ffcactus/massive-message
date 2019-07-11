package controller

import (
	"github.com/astaxie/beego"
	log "github.com/sirupsen/logrus"
	"massive-message/server/repository"
	"net/http"
	"strconv"
)

// Server controller process all the operation related to server.
type Server struct {
	beego.Controller
}

// GetServers retrieve and filter servers. Note that only certain properties can be used as filter key.
func (c *Server) GetServers() {
	var (
		start, count, orderby string = c.GetString("start"), c.GetString("count"), c.GetString("orderby")
		startInt, countInt    int64  = 0, -1
		parameterError        bool
	)
	if start != "" {
		_startInt, err := strconv.ParseInt(start, 10, 64)
		if err != nil || _startInt < 0 {
			parameterError = true
		} else {
			startInt = _startInt
		}
	}
	if count != "" {
		_countInt, err := strconv.ParseInt(count, 10, 64)
		// -1 means all.
		if err != nil || _countInt < -1 {
			parameterError = true
		} else {
			countInt = _countInt
		}
	}

	if parameterError {
		log.Warn("[Server-Controller] Get servers failed, parameter error.")
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	log.WithFields(log.Fields{"start": start, "count": count, "orderby": orderby}).Debug("[Server-Controller] Get servers.")
	servers, err := repository.GetServers(startInt, countInt, orderby)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("[Server-Controller] Get servers failed, repository operation failed.")
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}

	c.Data["json"] = servers
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.ServeJSON()
}

// GetServerByID handle the request that retrieve a server by it's ID.
func (c *Server) GetServerByID() {
	var (
		id = c.Ctx.Input.Param(":id")
	)

	server, err := repository.GetServerByID(id)
	if err != nil {
		log.WithFields(log.Fields{"id": id, "err": err}).Warn("[Controller] Get server by ID failed. repository operation failed.")
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}
	if server == nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.ServeJSON()
		return
	}
	c.Data["json"] = server
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.ServeJSON()
}
