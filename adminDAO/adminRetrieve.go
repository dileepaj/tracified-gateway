package adminDAO

import (
	"context"
	"fmt"

	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

/*GetCOCbyReceiver Retrieve All COC Object from COCCollection in DB by Receiver PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetPublicKeysOfFO() []string {
	//log.Debug("----------------- GetPublicKeysOfFO ------------------")
	fmt.Println("Get public keys of FO")
	var strArray []string
	var result []FOPK
	session, err := cd.adminConnect()
	if err != nil {
		//fmt.Println("An error occured 1")
		log.Error(err.Error())
	}
	if session != nil {
		//fmt.Println("All good with the session")
		defer session.EndSession(context.TODO())
		dbName := commons.GoDotEnvVariable("ADMINDBNAME")
		c := session.Client().Database(dbName).Collection("userkeys")

		findCursor, err1 := c.Find(context.TODO(), bson.M{"accounts.FO": true})
		if err1 != nil {
			// fmt.Println(err1)
			log.Error(err1.Error())
		}

		//fmt.Println("Cursor------------------", findCursor)
		if findErr := findCursor.All(context.TODO(), &result); findErr != nil {
			//fmt.Println("An error occured 2")
			panic(findErr)
		}

		for _, e := range result {
			for _, s := range e.Accounts {
				if len(s.Pk) == 56 {
					//fmt.Println("PKS  ", s.Pk)
					strArray = append(strArray, s.Pk)
				}
			}
		}
	}
	//fmt.Println("String Array in DB", strArray)
	return strArray
}
