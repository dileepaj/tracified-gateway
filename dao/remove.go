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

func (cd *Connection) RemoveFromTempOrphanList(Publickey string, SequenceNo int64) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("TESTTempOrphan")
	c.DeleteOne(context.TODO(), bson.M{"publickey": Publickey, "sequenceno": SequenceNo})

	if err != nil {
		fmt.Println("Error while remove from TempOrphan " + err.Error())
	}
	return err
}

//remove proof presentation protocol by proof name
func (cd *Connection) DeleteProofPresentationProtocolByProofName(proofName string) error{
	session, err := cd.connect()
	if err != nil{
		fmt.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("ProofProtocols")
	c.DeleteOne(context.TODO(), bson.M{"proofname" : proofName})
	if err != nil{
		fmt.Println("Error when removing the Protocol " + err.Error())
	}
	return err
}
