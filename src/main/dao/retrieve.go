
package dao

import(
	"main/model"
	"gopkg.in/mgo.v2/bson"
	// "gopkg.in/mgo.v2"
	"log"

)

func GetAllCOC(){

}

func GetCOCbySender(sender string)(model.COCCollectionBody,error){
	result :=model.COCCollectionBody{}
	c := CocCollection()
	err := c.Find(bson.M{"sender": sender}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result,err
}

func GetCOCbyReceiver(receiver string)(model.COCCollectionBody,error){
	result :=model.COCCollectionBody{}
	c := CocCollection()
	err := c.Find(bson.M{"receiver": receiver}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result,err
}
	