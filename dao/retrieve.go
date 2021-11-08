package dao

import (
	"context"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "fmt"

	"github.com/chebyrash/promise"
)

/*GetCOCbySender Retrieve All COC Object from COCCollection in DB by Sender PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbySender(sender string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		cursor, err1 := c.Find(context.TODO(), bson.M{"sender": sender})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				reject(err2)
			} else {
				resolve(result)
			}
		}

	})

	return p

}

/*GetLastCOCbySubAccount Retrieve the Last COC Object from COCCollection in DB by SubAccount PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCOCbySubAccount(subAccount string) *promise.Promise {
	result := model.COCCollectionBody{}
	result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		count, er := c.CountDocuments(context.TODO(), bson.M{"subaccount": subAccount})

		if er != nil {
			// fmt.Println(er)
			reject(er)
		}

		options := options.FindOne()
		options.SetSkip(count - 1)

		err1 := c.FindOne(context.TODO(), bson.M{"subaccount": subAccount}, options).Decode(&result)

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		}
		result2.Receiver = result.Receiver
		bumpSeq, err := strconv.Atoi(result.SequenceNo)
		if err == nil {
			fmt.Println(bumpSeq)
			// bumpSeq = bumpSeq
			fmt.Println(bumpSeq)
		}
		result2.SequenceNo = strconv.Itoa(bumpSeq)
		result2.SubAccount = result.SubAccount
		if result.Status == model.Pending.String() {
			result2.Available = false
			result2.Operation = "COC"
		} else {
			result2.Available = true
			result2.Operation = "COC"
		}
		resolve(result2)

	})

	return p

}

/*GetCOCbyReceiver Retrieve All COC Object from COCCollection in DB by Receiver PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyReceiver(receiver string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		cursor, err1 := c.Find(context.TODO(), bson.M{"receiver": receiver})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				reject(err2)
			} else {
				resolve(result)
			}
		}

	})

	return p

}

/*GetCOCbyAcceptTxn Retrieve a COC Object from COCCollection in DB by Accept TXN
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyAcceptTxn(accepttxn string) *promise.Promise {
	result := model.COCCollectionBody{}
	// p := promise.NewPromise()
	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		err1 := c.FindOne(context.TODO(), bson.M{"accepttxn": accepttxn}).Decode(&result)

		if err1 != nil {
			log.Error("Error while getting COC from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

/*GetCOCbyRejectTxn Retrieve a COC Object from COCCollection in DB by Reject TXN
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyRejectTxn(rejecttxn string) *promise.Promise {
	result := model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		err1 := c.FindOne(context.TODO(), bson.M{"rejecttxn": rejecttxn}).Decode(&result)

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}
	})

	return p

}

/*GetCOCbyStatus Retrieve All COC Object from COCCollection in DB by Status
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyStatus(status string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		cursor, err1 := c.Find(context.TODO(), bson.M{"status": status})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				reject(err2)
			} else {
				resolve(result)
			}
		}

	})

	return p

}

/*GetLastCOCbyIdentifier Retrieve Last COC Object from COCCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCOCbyIdentifier(identifier string) *promise.Promise {
	result := model.COCCollectionBody{}
	// result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		count, er := c.CountDocuments(context.TODO(), bson.M{"identifier": identifier})

		if er != nil {
			// fmt.Println(er)
			reject(er)
		}

		options := options.FindOne()
		options.SetSkip(count - 1)
		err1 := c.FindOne(context.TODO(), bson.M{"identifier": identifier}, options).Decode(&result)

		if err1 != nil {
			// fmt.Println(err1)

			reject(err1)
		}

		resolve(result)

	})

	return p

}

/*GetCOCByTxn Retrieve COC Object from COCCollection in DB by Txn
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCByTxn(txnHash string) *promise.Promise {
	result := model.COCCollectionBody{}
	// result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		er := c.FindOne(context.TODO(), bson.M{"txnhash": txnHash}).Decode(&result)

		if er != nil {
			// fmt.Println(er)
			reject(er)
		}

		resolve(result)

	})

	return p

}

/*GetLastTransactionbyIdentifier Retrieve Last Transaction Object from TransactionCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastTransactionbyIdentifier(identifier string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"identifier": identifier})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				reject(err2)
			} else {
				resolve(result[len(result)-1])
			}
		}

	})

	return p

}

/*GetFirstTransactionbyIdentifier Retrieve First Transaction Object from TransactionCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetFirstTransactionbyIdentifier(identifier string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		err1 := c.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&result)

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

/*GetTransactionsbyIdentifier Retrieve All Transaction Objects from TransactionCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetTransactionsbyIdentifier(identifier string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"identifier": identifier})
		err2 := cursor.All(context.TODO(), &result)

		if err1 != nil || err2 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}

	})

	return p

}

/*GetTransactionForTdpId Retrieve a Transaction Object from TransactionCollection in DB by TDPID
@author - Azeem Ashraf
*/
func (cd *Connection) GetTransactionForTdpId(TdpId string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()
	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while getting db connection " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		err1 := c.FindOne(context.TODO(), bson.M{"tdpid": TdpId}).Decode(&result)

		if err1 != nil {
			log.Error("Error while retrieving Transactions by tdpid " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetPreviousTransactions(limit int) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()
	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while getting db connection " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")

		count, er := c.CountDocuments(context.TODO(), bson.M{})

		if er != nil {
			log.Error("Error while get f.count " + err.Error())
			reject(er)
		}
		if count > int64(limit) {
			options := options.Find()
			options.SetSkip(count - int64(limit))
			cursor, err1 := c.Find(context.TODO(), bson.M{}, options)
			err2 := cursor.All(context.TODO(), &result)

			if err1 != nil || err2 != nil || len(result) == 0 {
				log.Error("Error while f.skip " + err1.Error())
				reject(err1)
			} else {
				resolve(result)
			}
		}

		cursor, err1 := c.Find(context.TODO(), bson.M{})
		err2 := cursor.All(context.TODO(), &result)

		if err1 != nil || err2 != nil || len(result) == 0 {
			log.Error("Error while f.All " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetPogTransaction(Identifer string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")

		err1 := c.FindOne(context.TODO(), bson.M{"identifier": Identifer, "tdpid": ""}).Decode(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}

	})

	return p

}

/*GetTransactionForTdpId Retrieve a Transaction Object from TransactionCollection in DB by TDPID
@author - Azeem Ashraf
*/
func (cd *Connection) GetAllTransactionForTdpId(TdpId string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()
	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to the db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")

		cursor, err1 := c.Find(context.TODO(), bson.M{"tdpid": TdpId})
		err2 := cursor.All(context.TODO(), &result)

		if err1 != nil || err2 != nil || len(result) == 0 {
			log.Error("Error while getting transactions " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

/*GetTdpIdForTransaction Retrieve a Transaction Object from TransactionCollection in DB by TXNID
@author - Azeem Ashraf
*/
func (cd *Connection) GetTdpIdForTransaction(Txn string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")

		err1 := c.FindOne(context.TODO(), bson.M{"txnhash": Txn}).Decode(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}

	})

	return p

}

/*GetOrphanbyIdentifier Retrieve a Transaction Object from OrphanCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetOrphanbyIdentifier(identifier string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Orphan")

		err1 := c.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

/*GetProfilebyIdentifier Retrieve a Profile Object from ProfileCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetProfilebyIdentifier(identifier string) *promise.Promise {
	result := model.ProfileCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Profiles")

		err1 := c.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

/*GetProfilebyProfileID Retrieve a Profile Object from ProfileCollection in DB by ProfileID
@author - Azeem Ashraf
*/
func (cd *Connection) GetProfilebyProfileID(ProfileID string) *promise.Promise {
	result := model.ProfileCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Profiles")

		err1 := c.FindOne(context.TODO(), bson.M{"profileid": ProfileID}).Decode(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

/*GetLastCertificatebyPublicKey Retrieve a Certificate Object from CertificateCollection in DB by PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCertificatebyPublicKey(PublicKey string) *promise.Promise {
	result := []model.CertificateCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Certificates")

		cursor, err1 := c.Find(context.TODO(), bson.M{"publickey": PublicKey})
		err2 := cursor.All(context.TODO(), &result)

		if err1 != nil || err2 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result[len(result)-1])

		}

	})

	return p

}

/*GetLastCertificatebyCertificateID Retrieve Last Certificate Object from CertificateCollection in DB by CertificateID
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCertificatebyCertificateID(CertificateID string) *promise.Promise {
	result := []model.CertificateCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Certificates")

		cursor, err1 := c.Find(context.TODO(), bson.M{"certificateid": CertificateID})

		err2 := cursor.All(context.TODO(), &result)

		if err1 != nil || err2 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result[len(result)-1])
		}
	})
	return p
}

/*GetAllCertificatebyPublicKey Retrieve All Certificate Objects from CertificateCollection in DB by PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetAllCertificatebyPublicKey(PublicKey string) *promise.Promise {
	result := []model.CertificateCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Certificates")

		cursor, err1 := c.Find(context.TODO(), bson.M{"publickey": PublicKey})

		err2 := cursor.All(context.TODO(), &result)

		if err1 != nil || err2 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}

	})

	return p

}

func (cd *Connection) GetTransactionId(tdpid string) *promise.Promise {
	result := []model.TransactionId{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"tdpid": tdpid})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				reject(err2)
			} else {
				resolve(result)
			}
		}

	})

	return p

}

func (cd *Connection) GetTransactionForTdpIdIdentifier(TdpId string, identifer string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		err1 := c.FindOne(context.TODO(), bson.M{"tdpid": TdpId, "identifer": identifer}).Decode(&result)

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}

	})

	return p

}

func (cd *Connection) GetAllTransactionForPK(Publickey string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"publickey": Publickey})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				reject(err2)
			} else {
				resolve(result)
			}
		}

	})

	return p

}

func (cd *Connection) GetAllTransactionForTxId(Txnhash string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"txnhash": Txnhash})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				reject(err2)
			} else {
				resolve(result)
			}
		}

	})

	return p

}

//GetSpecialForPkAndSeq ...
func (cd *Connection) GetSpecialForPkAndSeq(Publickey string, SequenceNo int64) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TempOrphan")
		err1 := c.FindOne(context.TODO(), bson.M{"publickey": Publickey, "sequenceno": SequenceNo}).Decode(&result)

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}

	})

	return p

}

//GetTransactionByTxnhash ..
func (cd *Connection) GetTransactionByTxnhash(Txnhash string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()
	session, err := cd.connect()
	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.

		if err != nil {
			log.Error("Error while fetching data from get db connection " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		err = c.FindOne(context.TODO(), bson.M{"txnhash": Txnhash}).Decode(&result)

		if err != nil {
			log.Error("Error while fetching data from Transactions " + err.Error())
			reject(err)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetAllApprovedOrganizations() *promise.Promise {
	result := []model.TestimonialOrganization{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
<<<<<<< HEAD
		defer session.Close()
		c := session.DB(dbName).C("Organizations")
		err1 := c.Find(bson.M{"status": model.Approved.String()}).All(&result)
=======
		defer session.EndSession(context.TODO())
>>>>>>> aa39307546625fc940c129b4e3bd3ccef1596e02

		c := session.Client().Database(dbName).Collection("Organizations")
		cursor, err1 := c.Find(context.TODO(), bson.M{"status": "Approved"})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				log.Error("Error while getting organizations from db " + err.Error())
				reject(err2)
			} else {
				resolve(result)
			}
		}
	})
	return p
}

func (cd *Connection) GetOrganizationByAuthor(publickey string) *promise.Promise {
	result := model.TestimonialOrganization{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Organizations")
		err1 := c.FindOne(context.TODO(), bson.M{"author": publickey}).Decode(&result)
		if err1 != nil {
			log.Error("Error while getting organization from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetOrganizationByAcceptTxn(acceptTxn string) *promise.Promise {
	result := model.TestimonialOrganization{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Organizations")
		err1 := c.FindOne(context.TODO(), bson.M{"accepttxn": acceptTxn}).Decode(&result)

		if err1 != nil {
			log.Error("Error while getting organization from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetOrganizationByRejectTxn(rejectTxn string) *promise.Promise {
	result := model.TestimonialOrganization{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Organizations")
		err1 := c.FindOne(context.TODO(), bson.M{"rejecttxn": rejectTxn}).Decode(&result)

		if err1 != nil {
			log.Error("Error while getting organization from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetTestimonialBySenderPublickey(publickey string) *promise.Promise {
	result := []model.Testimonial{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Testimonials")
		cursor, err1 := c.Find(context.TODO(), bson.M{"sender": publickey})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				log.Error("Error while getting Testimonial from db " + err.Error())
				reject(err2)
			} else {
				resolve(result)
			}
		}
	})
	return p
}

func (cd *Connection) GetTestimonialByRecieverPublickey(publickey string) *promise.Promise {
	result := []model.Testimonial{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Testimonials")
		cursor, err1 := c.Find(context.TODO(), bson.M{"reciever": publickey})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || len(result) == 0 {
				log.Error("Error while getting Testimonial from db " + err.Error())
				reject(err2)
			} else {
				resolve(result)
			}
		}
	})
	return p
}

func (cd *Connection) GetTestimonialByAcceptTxn(acceptTxn string) *promise.Promise {
	result := model.Testimonial{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Testimonials")
		err1 := c.FindOne(context.TODO(), bson.M{"accepttxn": acceptTxn}).Decode(&result)

		if err1 != nil {
			log.Error("Error while getting Testimonials from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetTestimonialByRejectTxn(rejectTxn string) *promise.Promise {
	result := model.Testimonial{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Testimonials")
		err1 := c.FindOne(context.TODO(), bson.M{"rejecttxn": rejectTxn}).Decode(&result)

		if err1 != nil {
			log.Error("Error while getting organization from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetLastOrganizationbySubAccount(subAccount string) *promise.Promise {
	result := model.TestimonialOrganization{}
	result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB(dbName).C("Organizations")

		count, er := c.Find(bson.M{"subaccount": subAccount}).Count()
		if er != nil {
			// fmt.Println(er)
			reject(er)
		}

		err1 := c.Find(bson.M{"subaccount": subAccount}).Skip(count - 1).One(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		}
		result2.Receiver = result.ApprovedBy
		bumpSeq, err := strconv.Atoi(result.SequenceNo)
		if err == nil {
			fmt.Println(bumpSeq)
			bumpSeq = bumpSeq
			fmt.Println(bumpSeq)
		}
		result2.SequenceNo = strconv.Itoa(bumpSeq)
		result2.SubAccount = result.SubAccount
		if result.Status == model.Pending.String() {
			result2.Available = false
			result2.Operation = "Organization"
		} else {
			result2.Available = true
			result2.Operation = "Organization"
		}
		resolve(result2)

	})

	return p

}

func (cd *Connection) GetLastTestimonialbySubAccount(subAccount string) *promise.Promise {
	result := model.Testimonial{}
	result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			// fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB(dbName).C("Testimonials")

		count, er := c.Find(bson.M{"subaccount": subAccount}).Count()
		if er != nil {
			// fmt.Println(er)
			reject(er)
		}

		err1 := c.Find(bson.M{"subaccount": subAccount}).Skip(count - 1).One(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		}
		result2.Receiver = result.Reciever
		bumpSeq, err := strconv.Atoi(result.SequenceNo)
		if err == nil {
			fmt.Println(bumpSeq)
			bumpSeq = bumpSeq
			fmt.Println(bumpSeq)
		}
		result2.SequenceNo = strconv.Itoa(bumpSeq)
		result2.SubAccount = result.Subaccount
		if result.Status == model.Pending.String() {
			result2.Available = false
			result2.Operation = "Testimonial"
		} else {
			result2.Available = true
			result2.Operation = "Testimonial"
		}
		resolve(result2)

	})

	return p

}

func (cd *Connection) GetPendingAndRejectedOrganizations() *promise.Promise {
	result := []model.TestimonialOrganization{}

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.Close()
		c := session.DB(dbName).C("Organizations")
		err1 := c.Find(bson.M{"status": bson.M{"$in": []string{model.Pending.String(), model.Rejected.String()}}}).All(&result)

		if err1 != nil || len(result) == 0 {
			log.Error("Error while getting organizations from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}

	})
	return p
}
