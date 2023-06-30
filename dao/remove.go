package dao

import (
	"context"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

/*RemoveFromOrphanage Remove a single Transaction Object from the OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) RemoveFromOrphanage(Identifier string) error {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Orphan")
	c.DeleteOne(context.TODO(), bson.M{"identifier": Identifier})

	if err != nil {
		logger.LogWriter("Error while remove from Orphan "+err.Error(), constants.ERROR)
	}
	return err
}

func (cd *Connection) RemoveFromTempOrphanList(Publickey string, SequenceNo int64) error {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("TempOrphan")
	c.DeleteOne(context.TODO(), bson.M{"publickey": Publickey, "sequenceno": SequenceNo})

	if err != nil {
		logger.LogWriter("Error while remove from TempOrphan "+err.Error(), constants.ERROR)
	}
	return err
}

//remove proof presentation protocol by proof name
func (cd *Connection) DeleteProofPresentationProtocolByProofName(proofName string) error {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("ProofProtocols")
	c.DeleteOne(context.TODO(), bson.M{"proofname": proofName})
	if err != nil {
		logger.LogWriter("Error while remove from Protocol "+err.Error(), constants.ERROR)
	}
	return err
}
