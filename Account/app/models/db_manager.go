package models

/*
*TODO:
*	Should use revel's plug-in mechanism here.
*/

import (
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
)

const (
	DbSection = "db"
	Ip = "ip"
)

type DbManager struct {
	session *mgo.Session
}

func NewDbManager() (*DbManager, error) {
	revel.Config.SetSection(DbSection)
	ip, found := revel.Config.String(Ip)
	if !found {
		revel.ERROR.Fatal("Cannot load database ip from app.conf")
	}

	session, err := mgo.Dial(ip)
	if err != nil {
		return nil, err
	}

	return &DbManager{session}, nil
}

func (manager *DbManager) Close() {
	manager.session.Close()
}
