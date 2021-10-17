package adminDAO

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"gopkg.in/mgo.v2"
)

type AdminStore struct {
	Host          string
	Port          string
	Username      string
	Password      string
	DBName        string
}

var adminMgoSession *mgo.Session
var adminMgoConnectionUrl string

func GetAdminMongoSession() (*mgo.Session, error) {
	if adminMgoSession == nil {
		var err error
		adminMgoSession, err = mgo.Dial(adminMgoConnectionUrl)
		if err != nil {
			return nil, err
		}
	}
	return adminMgoSession.Clone(), nil
}

func ConstructAdminConnectionPool() {
	username:=     commons.GoDotEnvVariable("ADMINDBUSERNAME")
	password:=     commons.GoDotEnvVariable("ADMINDBPASSWORD")
	dbName:=     commons.GoDotEnvVariable("ADMINDBNAME")
	host:=     commons.GoDotEnvVariable("ADMINDBHOST")
	port:=     commons.GoDotEnvVariable("ADMINDBPORT")
	adminMgoConnectionUrl = "mongodb://"+username+":"+password+"@"+host+":"+port+"/?authSource="+dbName
}

