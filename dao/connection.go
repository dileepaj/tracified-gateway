package dao

import (
	"fmt"
	"gopkg.in/mgo.v2"
)


type Connection struct {
}

func (cd *Connection) connect()(*mgo.Session,error) {
	session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds143953.mlab.com:43953/tracified-gateway")
	if err != nil {
		fmt.Println(err)
	}
	return session,err

}




//  Connection:=connect()
