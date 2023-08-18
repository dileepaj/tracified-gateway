package adminDAO

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"go.mongodb.org/mongo-driver/mongo"
)

//Get db name from .env file
var adminDBName = commons.GoDotEnvVariable("ADMIN_BE_DB_NAME")

/*Connection The Mgo Connection
@author - Azeem Ashraf
*/
type Connection struct {
}

func (cd *Connection) adminConnect() (mongo.Session, error) {
	return GetAdminMongoSession()
}
