package dao

import (
	"main/model"
	// "gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
)

func InsertCoc(Coc model.COCCollectionBody) error {
	
	session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds143953.mlab.com:43953/tracified-gateway")
	if err != nil {
		panic(err)
	}
	c := session.DB("tracified-gateway").C("COC")
	err1 := c.Insert(Coc)
	if err1 != nil {
		log.Fatal(err1)
	}
	return err
}
