package dao

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
)

func (cd *Connection) InsertIdentifier(identifierModel apiModel.IdentifierModel) (error,int) {
	// session, err := cd.connect()
	// if err != nil {
	// 	log.Println("Error while getting session " + err.Error())
	// }
	// defer session.EndSession(context.TODO())

	// c := session.Client().Database(dbName).Collection("IdentifierMap")
	// _, err = c.InsertOne(context.TODO(), id)

	// if err != nil {
	// 	log.Println("Error while inserting to TempOrphan " + err.Error())
	// }
	// return err

	// call common insert
	return nil,200
}