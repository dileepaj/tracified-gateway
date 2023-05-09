package dao

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get db name from .env file
var dbName = commons.GoDotEnvVariable("DBNAME")

type Connection struct{}

func (cd *Connection) connect() (mongo.Session, error) {
	return commons.GetMongoSession()
}
