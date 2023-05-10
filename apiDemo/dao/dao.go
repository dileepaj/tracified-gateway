package dao

import (
	"context"

	"github.com/chebyrash/promise"
	"github.com/dileepaj/tracified-gateway/apiDemo/dao/connections"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DbName = commons.GoDotEnvVariable("DBNAME")
)

type IndexType interface {
	[]model.RSAPublickey
}

func Index[inedxObj IndexType](collection string, searchMap map[string]any, objectType inedxObj) *promise.Promise {

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		var returnIndex = objectType
		db := connections.DemoConnection{}
		customLogger := utilities.NewCustomLogger()
		session, err := db.DbConnection()

		bsonDoc := bson.D{}
		for key, value := range searchMap {
			bsonDoc = append(bsonDoc, bson.E{Key: key, Value: value})
		}

		if err != nil {
			customLogger.LogWriter("Error when connecting to DB "+err.Error(), 3)
			reject(err)
		}
		defer session.EndSession(context.TODO())

		rst := session.Client().Database(DbName).Collection(collection)
		cursor, err := rst.Find(context.TODO(), bsonDoc)
		if err != nil {
			reject(err)
		} else {
			err2 := cursor.All(context.TODO(), &returnIndex)

			if err2 != nil || len(returnIndex) == 0 {
				customLogger.LogWriter("Error while getting organizations from db "+err2.Error(), 3)
				reject(err2)
			} else {
				resolve(returnIndex)
			}
		}
	})
	return p
}

type CreateType interface {
	model.TransactionCollectionBody
}

func Create[T CreateType](model T, collection string) (string, error) {
	db := connections.DemoConnection{}
	customerLogger := utilities.NewCustomLogger()
	session, err := db.DbConnection()
	if err != nil {
		customerLogger.LogWriter("Error in DB connection : "+err.Error(), 3)
	}
	defer session.EndSession(context.TODO())
	c, err := session.Client().Database(DbName).Collection(collection).InsertOne(context.TODO(), model)
	if err != nil {
		customerLogger.LogWriter("Error when getting the DB session : "+err.Error(), 3)
	}
	id := c.InsertedID.(primitive.ObjectID)
	return id.Hex(), err

}

type ShowType interface {
}

func Show[T ShowType](idName string, id T, collection string, searchMap map[string]any, object ShowType) *promise.Promise {
	result := object

	bsonDoc := bson.D{}
	for key, value := range searchMap {
		bsonDoc = append(bsonDoc, bson.E{Key: key, Value: value})
	}

	promise := promise.New(func(resolve func(interface{}), reject func(error)) {
		db := connections.DemoConnection{}
		session, err := db.DbConnection()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		dbclient := session.Client().Database(DbName).Collection(collection)
		err1 := dbclient.FindOne(context.TODO(), bsonDoc).Decode(&result)
		if err1 != nil {
			reject(err)
		} else {
			resolve(result)
		}
	})
	return promise
}

func Update(findBy string, value string, update primitive.M, projectionData primitive.M, collection string) *mongo.SingleResult {
	//TODO :  Need to add Update logic
	return nil
}

func Remove(idName string, id, collection string) (int64, error) {
	db := connections.DemoConnection{}
	customerLogger := utilities.NewCustomLogger()
	session, err := db.DbConnection()
	if err != nil {
		customerLogger.LogWriter("Error while getting session "+err.Error(), 3)
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(DbName).Collection(collection)
	rst, err := c.DeleteOne(context.TODO(), bson.M{idName: id})

	if err != nil {
		customerLogger.LogWriter("Error while remove from Orphan "+err.Error(), 3)
	}
	return rst.DeletedCount, err
}
