package dao

import (
	// "fmt"

	// "fmt"
	// "log"
	// "main/api/routes"
	// "net/http"

	// "github.com/gorilla/handlers"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

var Session *mgo.Session

type Connection struct {
}

func (cd *Connection) connect() {
	session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds143953.mlab.com:43953/tracified-gateway")
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	Session = session

}

func (cd *Connection) GetSession() *mgo.Session {
	if Session == nil {
		cd.connect()
	}
	return Session
}

//  Connection:=connect()
