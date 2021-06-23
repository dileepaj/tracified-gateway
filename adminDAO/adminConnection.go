package adminDAO

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"gopkg.in/mgo.v2"
)
//Get db name from .env file
var adminDBName = commons.GoDotEnvVariable("ADMINDBNAME")
/*Connection The Mgo Connection
@author - Azeem Ashraf
*/
type Connection struct {
}

func (cd *Connection) adminConnect()(*mgo.Session,error) {
	return GetAdminMongoSession()
}
