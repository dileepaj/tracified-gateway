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
		fmt.Println(err)
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("Orphan")
	err1 := c.Remove(bson.M{"identifier": Identifier})
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}
