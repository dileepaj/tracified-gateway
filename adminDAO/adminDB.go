package adminDAO

import (
	"context"

	"github.com/dileepaj/tracified-gateway/commons"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminStore struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var (
	adminMgoSession       mongo.Session
	adminMgoConnectionUrl string
)

func GetAdminMongoSession() (mongo.Session, error) {
	if adminMgoSession == nil {
		adminMgoConnectionUrl = commons.GoDotEnvVariable("ADMIN_BE_MONGODB_URI")
		var err error
		adminMgoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(adminMgoConnectionUrl))
		if err != nil {
			return nil, err
		}
		adminMgoSession, err = adminMgoClient.StartSession()
		if err != nil {
			return nil, err
		}
	}
	return adminMgoSession, nil
}
