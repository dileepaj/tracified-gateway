package connections

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get db name from .env file
var dbName = commons.GoDotEnvVariable("DBNAME")

type DemoConnection struct {
}

func (dcd *DemoConnection) DbConnection() (mongo.Session, error) {
	return commons.GetMongoSession()
}
