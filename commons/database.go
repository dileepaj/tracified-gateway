package commons

import (
	"fmt"
	"github.com/astaxie/beego/core/config"
	"gopkg.in/mgo.v2"
)

type Store struct {
	Host          string
	Port          string
	Username      string
	Password      string
	DBName        string
}

var mgoSession *mgo.Session
var mongoConnectionUrl string

func GetMongoSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mongoConnectionUrl)
		if err != nil {
			return nil, err
		}
	}
	return mgoSession.Clone(), nil
}

func ConstructConnectionPool(conf config.Configer) {
	username,err :=     conf.String("datastore" + "::username")
	password,err :=     conf.String("datastore" + "::password")
	dbName,err :=     conf.String("datastore" + "::dbName")
	host,err :=     conf.String("datastore" + "::host")
	port,err :=     conf.String("datastore" + "::port")
	if err != nil {
		fmt.Println(err.Error())
	}
	mongoConnectionUrl = "mongodb://"+username+":"+password+"@"+host+":"+port+"/"+dbName
}

