package dao

import (
	"github.com/dileepaj/tracified-gateway/model"
)

func (cd *Connection) InsertTempOrphan(tempOrphan model.TransactionCollectionBody) (error, int) {
	// fmt.Println("--------------------------- InsertSpecialToTempOrphan ------------------------")
	// session, err := cd.connect()
	// if err != nil {
	// 	fmt.Println("Error while getting session " + err.Error())
	// }
	// defer session.EndSession(context.TODO())

	// c := session.Client().Database(dbName).Collection("TempOrphan")
	// _, err = c.InsertOne(context.TODO(), Coc)

	// if err != nil {
	// 	fmt.Println("Error while inserting to TempOrphan " + err.Error())
	// }
	// return err

	// call common insert
	Create(tempOrphan, "TempOrphan")
	return nil, 200
}
