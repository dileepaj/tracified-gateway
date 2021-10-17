package commons

import (
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

func ConstructConnectionPool() {
	username:=     GoDotEnvVariable("DBUSERNAME")
	password:=     GoDotEnvVariable("DBPASSWORD")
	dbName:=     GoDotEnvVariable("DBNAME")
	host:=     GoDotEnvVariable("DBHOST")
	port:=     GoDotEnvVariable("DBPORT")
	mongoConnectionUrl = "mongodb://"+username+":"+password+"@"+host+":"+port+"/"+dbName
}

