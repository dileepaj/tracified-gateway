package dao

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

/*RemoveFromOrphanage Remove a single Transaction Object from the OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) RemoveFromOrphanage(Identifier string) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Orphan")
	c.DeleteOne(context.TODO(), bson.M{"identifier": Identifier})

	if err != nil {
		fmt.Println("Error while remove from Orphan " + err.Error())
	}
	return err
}
