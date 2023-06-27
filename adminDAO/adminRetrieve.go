package adminDAO

import (
	"context"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

/*
GetCOCbyReceiver Retrieve All COC Object from COCCollection in DB by Receiver PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetPublicKeysOfFO() ([]string, error) {
	//log.Debug("----------------- GetPublicKeysOfFO ------------------")
	var strArray []string
	var result []FOPK
	session, err := cd.adminConnect()
	if err != nil {
		log.Error(err.Error())
	}
	if session != nil {
		defer session.EndSession(context.TODO())
		dbName := commons.GoDotEnvVariable("ADMINDBNAME")
		c := session.Client().Database(dbName).Collection("userkeys")

		findCursor, err1 := c.Find(context.TODO(), bson.M{"accounts.FO": true})
		if err1 != nil {
			log.Error(err1.Error())
		}

		if findErr := findCursor.All(context.TODO(), &result); findErr != nil {
			logger := utilities.NewCustomLogger()
			logger.LogWriter("Error GetPublicKeysOfFO : "+err.Error(), constants.ERROR)
			return strArray, findErr
		}

		for _, e := range result {
			for _, s := range e.Accounts {
				if len(s.Pk) == 56 {
					strArray = append(strArray, s.Pk)
				}
			}
		}
	}
	return strArray, nil
}
