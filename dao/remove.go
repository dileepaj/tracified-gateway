package dao

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"

)
/*RemoveFromOrphanage Remove a single Transaction Object from the OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) RemoveFromOrphanage(Identifier string) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session "+err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Orphan")
	err= c.Remove(bson.M{"identifier": Identifier})
	if err != nil {
		fmt.Println("Error while remove from Orphan "+err.Error())
	}
	return err
}
