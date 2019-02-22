package dao

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"

)

func (cd *Connection) RemoveFromBuffer(Identifier string) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("Buffer")
	err1 := c.Remove(bson.M{"identifier": Identifier})
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}
