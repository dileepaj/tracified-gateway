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

var adminMgoSession mongo.Session
var adminMgoConnectionUrl string

func GetAdminMongoSession() (mongo.Session, error) {
	//log.Println("-------------------------Get admin mongo session---------------")
	if adminMgoSession == nil {
		var err error
		adminMgoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://admin-user:TAqa123@ap-cluster-0.wrst2.mongodb.net/tracified-admin-backend-qa?retryWrites=true&w=majority"))
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

func ConstructAdminConnectionPool() {
	username := commons.GoDotEnvVariable("ADMINDBUSERNAME")
	password := commons.GoDotEnvVariable("ADMINDBPASSWORD")
	dbName := commons.GoDotEnvVariable("ADMINDBNAME")
	host := commons.GoDotEnvVariable("ADMINDBHOST")
	port := commons.GoDotEnvVariable("ADMINDBPORT")

	branch := commons.GoDotEnvVariable("BRANCH_NAME")

	if branch == "production" {
		adminMgoConnectionUrl = "mongodb://" + username + ":" + password + "@" + host + ":" + port + "/?authSource=" + dbName
	} else {
		adminMgoConnectionUrl = commons.GoDotEnvVariable("ADMIN_ATLAS")
	}
}
