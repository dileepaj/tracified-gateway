package dao

import (
	"context"
	"fmt"
	"log"

	"github.com/chebyrash/promise"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (cd *Connection) InsertRecords(Coc model.Testing) error {
	fmt.Println("--------------------------- InsertRecords ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("Records")
	_, err = c.InsertOne(context.TODO(), Coc)
	if err != nil {
		fmt.Println("Error while inserting to Records " + err.Error())
	}
	return err
}

func (cd *Connection) GetRecordsByID(id string) *promise.Promise {
	result := model.Testing{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Records")
		err1 := c.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}
	})

	return p

}

func (cd *Connection) UpdateRecordsById(selector model.Testing, update model.Testing) error {
	fmt.Println("--------------------------- Update Records ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.Testing{
		Id:            selector.Id,
		Name:          update.Name,
		Address:       update.Address,
		Designation:   update.Designation,
		Specification: update.Specification,
	}

	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}

	var filter bson.M
	err = bson.Unmarshal(pByte, &filter)
	if err != nil {
		return err
	}

	pByte, err = bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("Records")
	_, err = c.UpdateOne(context.TODO(), bson.M{"id": selector.Id}, bson.D{{Key: "$set", Value: updateNew}})

	if err != nil {
		fmt.Println("Error while updating proof protocols " + err.Error())
	}
	return err
}

func (cd *Connection) RemoveFromRecords(id string) error {
	fmt.Println("--------------------------- Remove Records ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Records")
	c.DeleteOne(context.TODO(), bson.M{"id": id})

	if err != nil {
		fmt.Println("Error while remove from Records " + err.Error())
	}
	return err
}

// func (cd *Connection) GetAllRecords(Coc model.Testing) error {
// 	fmt.Println("--------------------------- Get All Records ------------------------")
// 	session, err := cd.connect()
// 	if err != nil {
// 		fmt.Println("Error while getting session " + err.Error())
// 	}
// 	defer session.EndSession(context.TODO())
// 	c := session.Client().Database(dbName).Collection("Records")
// 	_, err = c.Find(context.TODO(), Coc)

// 	if err != nil {
// 		fmt.Println("Error while getting all Records " + err.Error())
// 	}
// 	fmt.Println(err)
// 	return err
// }

func (cd *Connection) GetAllRecords(Coc model.Testing) *promise.Promise {
	//result := model.Testing{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Records")
		cursor, err1 := c.Find(context.TODO(), bson.D{})

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			fmt.Println("---------------Reading 1 -----------------------")
			var results []bson.M
			if err = cursor.All(context.TODO(), &results); err != nil {
				log.Fatal(err)
			}
			fmt.Println("---------------Reading 2 -----------------------")
			for _, result := range results {
				logrus.Error("---------------Reading 3 -----------------------", result)
				resolve(result)
			}
		}
	})

	return p

}

func (cd *Connection) DeleteAllRecords(Coc model.Testing) error {
	fmt.Println("--------------------------- Delete All Records ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("Records")
	_, err = c.DeleteMany(context.TODO(), Coc)

	if err != nil {
		fmt.Println("Error while Deleting all Records " + err.Error())
	}
	return err
}
