package adminDAO

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

/*GetCOCbyReceiver Retrieve All COC Object from COCCollection in DB by Receiver PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetPublicKeysOfFO() []string {
	log.Debug("----------------- GetPublicKeysOfFO ------------------")
	var strArray []string
	var result []FOPK
	session, err := cd.adminConnect()
	if err != nil{
		log.Error(err.Error())
	}
	if session != nil {
		defer session.Close()
		c := session.DB("admin-db").C("userkeys")
		err1 := c.Find(bson.M{"accounts.FO": true}).All(&result)
		if err1 != nil {
			// fmt.Println(err1)
			log.Error(err1.Error())
		}
		for _, e := range result {
			for _, s := range e.Accounts {
				if len(s.Pk)==56 {
					strArray = append(strArray, s.Pk)
				}
			}
		}
	}
	return strArray
}
