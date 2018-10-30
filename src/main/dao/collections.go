package dao

import (
	// "main/dao/connection"
	// "gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

// type CocCollection struct {
// 	c *mgo.Collection
// }

func CocCollection() *mgo.Collection{

	session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds143953.mlab.com:43953/tracified-gateway")
	if err != nil {
		panic(err)
	}
	// object:=Connection{}
	// session := object.GetSession()
	c := session.DB("tracified-gateway").C("COC")
	return c
}
