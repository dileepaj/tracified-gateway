package commons

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var mgoSession mongo.Session
var mongoConnectionUrl string

func GetMongoSession() (mongo.Session, error) {
	if mgoSession == nil {
		mongoConnectionUrl = GoDotEnvVariable("DB_URI")
		var err error
		mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoConnectionUrl))
		if err != nil {
			return nil, err
		}
		mgoSession, err = mongoClient.StartSession()
		if err != nil {
			return nil, err
		}
	}

	return mgoSession, nil
}
