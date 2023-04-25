package dao

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	notificationhandler "github.com/dileepaj/tracified-gateway/services/notificationHandler.go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chebyrash/promise"
)

/*
GetCOCbySender Retrieve All COC Object from COCCollection in DB by Sender PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbySender(sender string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetLastCOCbySubAccount Retrieve the Last COC Object from COCCollection in DB by SubAccount PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCOCbySubAccount(subAccount string) *promise.Promise {
	result := model.COCCollectionBody{}
	result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

		if result.Status == model.Expired.String() {
			result2.Expiration = true
		} else {
			result2.Expiration = false
		}

		resolve(result2)
	})

	return p
}

/*
GetCOCbyReceiver Retrieve All COC Object from COCCollection in DB by Receiver PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyReceiver(receiver string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetCOCbyAcceptTxn Retrieve a COC Object from COCCollection in DB by Accept TXN
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyAcceptTxn(accepttxn string) *promise.Promise {
	result := model.COCCollectionBody{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("COC")
		err1 := c.FindOne(context.TODO(), bson.M{"accepttxn": accepttxn}).Decode(&result)

		if err1 != nil {
			logrus.Info("Error while getting COC from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

/*
GetCOCbyRejectTxn Retrieve a COC Object from COCCollection in DB by Reject TXN
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyRejectTxn(rejecttxn string) *promise.Promise {
	result := model.COCCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetCOCbyStatus Retrieve All COC Object from COCCollection in DB by Status
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCbyStatus(status string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetLastCOCbyIdentifier Retrieve Last COC Object from COCCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCOCbyIdentifier(identifier string) *promise.Promise {
	result := model.COCCollectionBody{}
	// result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetCOCByTxn Retrieve COC Object from COCCollection in DB by Txn
@author - Azeem Ashraf
*/
func (cd *Connection) GetCOCByTxn(txnHash string) *promise.Promise {
	result := model.COCCollectionBody{}
	// result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetLastTransactionbyIdentifier Retrieve Last Transaction Object from TransactionCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastTransactionbyIdentifier(identifier string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

func (cd *Connection) GetLastTransactionbyIdentifierNotSplitParent(identifier, tenantId string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"identifier": identifier, "txntype": bson.M{"$ne": "5"}, "tenantid": tenantId})

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

/*
GetFirstTransactionbyIdentifier Retrieve First Transaction Object from TransactionCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetFirstTransactionbyIdentifier(identifier string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetTransactionsbyIdentifier Retrieve All Transaction Objects from TransactionCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetTransactionsbyIdentifier(identifier string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetTransactionForTdpId Retrieve a Transaction Object from TransactionCollection in DB by TDPID
@author - Azeem Ashraf
*/
func (cd *Connection) GetTransactionForTdpId(TdpId string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*Get total transaction count in a collection
 */
func (cd *Connection) GetTransactionCount() *promise.Promise {
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while getting db connection " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		count, er := c.CountDocuments(context.TODO(), bson.M{})
		if er != nil {
			log.Error("Error while retrieving Transactions by tdpid " + err.Error())
			reject(er)
		} else {
			resolve(count)
		}
	})
	return p
}

func (cd *Connection) GetPreviousTransactions(perPage int, page int, NoPage int) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while getting db connection " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")

		count, er := c.CountDocuments(context.TODO(), bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "7", "10"}}})
		// count only genesis, TDP, splitparent, splitchild and COC transactions
		if er != nil {
			log.Error("Error while get f.count " + err.Error())
			reject(er)
		}
		if count < int64(perPage) {
			cursor, err1 := c.Find(context.TODO(), bson.M{})
			err2 := cursor.All(context.TODO(), &result)

			if err1 != nil || err2 != nil || len(result) == 0 {
				log.Error("Error while f.All " + err1.Error())
				reject(err1)
			} else {
				resolve(result)
			}

		}

		opt := options.Find().SetSkip(count - int64(page)*int64(perPage)).SetLimit(int64(perPage))
		cursor, err1 := c.Find(context.TODO(), bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "7", "10"}}}, opt)
		// retrieve only genesis, TDP, splitparent, splitchild and COC transactions
		err2 := cursor.All(context.TODO(), &result)

		if err1 != nil || err2 != nil || len(result) == 0 {
			log.Error("Error while f.skip " + err1.Error())
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetTransactionForTdpId Retrieve a Transaction Object from TransactionCollection in DB by TDPID
@author - Azeem Ashraf
*/
func (cd *Connection) GetAllTransactionForTdpId(TdpId string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetTdpIdForTransaction Retrieve a Transaction Object from TransactionCollection in DB by TXNID
@author - Azeem Ashraf
*/
func (cd *Connection) GetTdpIdForTransaction(Txn string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetOrphanbyIdentifier Retrieve a Transaction Object from OrphanCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetOrphanbyIdentifier(identifier string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetProfilebyIdentifier Retrieve a Profile Object from ProfileCollection in DB by Identifier
@author - Azeem Ashraf
*/
func (cd *Connection) GetProfilebyIdentifier(identifier string) *promise.Promise {
	result := model.ProfileCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetProfilebyProfileID Retrieve a Profile Object from ProfileCollection in DB by ProfileID
@author - Azeem Ashraf
*/
func (cd *Connection) GetProfilebyProfileID(ProfileID string) *promise.Promise {
	result := model.ProfileCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetLastCertificatebyPublicKey Retrieve a Certificate Object from CertificateCollection in DB by PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCertificatebyPublicKey(PublicKey string) *promise.Promise {
	result := []model.CertificateCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetLastCertificatebyCertificateID Retrieve Last Certificate Object from CertificateCollection in DB by CertificateID
@author - Azeem Ashraf
*/
func (cd *Connection) GetLastCertificatebyCertificateID(CertificateID string) *promise.Promise {
	result := []model.CertificateCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

/*
GetAllCertificatebyPublicKey Retrieve All Certificate Objects from CertificateCollection in DB by PublicKey
@author - Azeem Ashraf
*/
func (cd *Connection) GetAllCertificatebyPublicKey(PublicKey string) *promise.Promise {
	result := []model.CertificateCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

// GetSpecialForPkAndSeq ...
func (cd *Connection) GetSpecialForPkAndSeq(Publickey string, SequenceNo int64) *promise.Promise {
	// fmt.Println("Address to get special ", Publickey)
	// fmt.Println("Sequence no of the address ", SequenceNo)
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TempOrphan")
		err1 := c.FindOne(context.TODO(), bson.M{"publickey": Publickey, "sequenceno": SequenceNo}).Decode(&result)

		if err1 != nil {
			reject(err1)
		} else {
			resolve(result)
		}
	})

	return p
}

// GetTransactionByTxnhash ..
func (cd *Connection) GetTransactionByTxnhash(Txnhash string) *promise.Promise {
	result := model.TransactionCollectionBody{}
	// p := promise.NewPromise()
	session, err := cd.connect()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("Organizations")
		cursor, err1 := c.Find(context.TODO(), bson.M{"status": model.Approved.String()})

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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
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

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Organizations")

		count, er := c.CountDocuments(context.TODO(), bson.M{"subaccount": subAccount})
		if er != nil {
			reject(er)
		}

		options := options.FindOne()
		options.SetSkip(count - 1)

		err1 := c.FindOne(context.TODO(), bson.M{"subaccount": subAccount}, options).Decode(&result)
		if err1 != nil {
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

		if result.Status == model.Expired.String() {
			result2.Expiration = true
		} else {
			result2.Expiration = false
		}

		resolve(result2)
	})

	return p
}

func (cd *Connection) GetLastTestimonialbySubAccount(subAccount string) *promise.Promise {
	result := model.Testimonial{}
	result2 := apiModel.GetSubAccountStatusResponse{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Testimonials")
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

		if result.Status == model.Expired.String() {
			result2.Expiration = true
		} else {
			result2.Expiration = false
		}

		resolve(result2)
	})

	return p
}

func (cd *Connection) GetPendingAndRejectedOrganizations() *promise.Promise {
	result := []model.TestimonialOrganization{}

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Organizations")

		cursor, err := c.Find(context.TODO(), bson.M{"status": bson.M{"$in": []string{model.Pending.String(), model.Rejected.String()}}})
		err1 := cursor.All(context.TODO(), &result)

		if err1 != nil || len(result) == 0 {
			log.Error("Error while getting organizations from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetAllTransactionForPK_Paginated(Publickey string, page int, perPage int) *promise.Promise {
	result := model.TransactionCollectionBodyWithCount{}
	// p := promise.NewPromise()
	fmt.Println("----------r1")
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}
		defer session.EndSession(context.TODO())
		fmt.Println("----------r2")
		c := session.Client().Database(dbName).Collection("Transactions")
		// count, er := c.CountDocuments(context.TODO(), bson.M{"publickey": Publickey})
		count, er := c.CountDocuments(context.TODO(), bson.M{"$and": []interface{}{bson.M{"publickey": Publickey}, bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "10"}}}}})

		if er != nil {
			log.Error("Error while get f.count " + err.Error())
			reject(er)
		}

		offset := int64(page) * int64(perPage)
		if count < offset {
			perPage = int(count + int64(perPage) - offset)
			offset = count

		}

		opt := options.Find().SetSkip(count - offset).SetLimit(int64(perPage))
		cursor, err1 := c.Find(context.TODO(), bson.M{"$and": []interface{}{bson.M{"publickey": Publickey}, bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "10"}}}}}, opt)

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &(result.Transactions))
			result.Count = int64(count)
			if err2 != nil || len(result.Transactions) == 0 {
				reject(err2)
			} else {
				fmt.Println("----------r3")
				resolve(result)
			}
		}
	})

	return p
}

func (cd *Connection) GetAllTransactionForTdpId_Paginated(TdpId string, page int, perPage int) *promise.Promise {
	result := model.TransactionCollectionBodyWithCount{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to the db " + err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		count, er := c.CountDocuments(context.TODO(), bson.M{"$and": []interface{}{bson.M{"tdpid": TdpId}, bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "10"}}}}})

		if er != nil {
			log.Error("Error while get f.count " + err.Error())
			reject(er)
		}

		offset := int64(page) * int64(perPage)
		if count < offset {
			perPage = int(count + int64(perPage) - offset)
			offset = count

		}

		opt := options.Find().SetSkip(count - int64(page)*int64(perPage)).SetLimit(int64(perPage))
		cursor, err1 := c.Find(context.TODO(), bson.M{"$and": []interface{}{bson.M{"tdpid": TdpId}, bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "10"}}}}}, opt)
		err2 := cursor.All(context.TODO(), &(result.Transactions))
		result.Count = int64(count)
		if err1 != nil || err2 != nil || len(result.Transactions) == 0 {
			log.Error("Error while getting transactions " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetTransactionsbyIdentifier_Paginated(identifier string, page int, perPage int) *promise.Promise {
	result := model.TransactionCollectionBodyWithCount{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		count, er := c.CountDocuments(context.TODO(), bson.M{"$and": []interface{}{bson.M{"identifier": identifier}, bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "10"}}}}})

		if er != nil {
			log.Error("Error while get f.count " + err.Error())
			reject(er)
		}

		offset := int64(page) * int64(perPage)
		if count < offset {
			perPage = int(count + int64(perPage) - offset)
			offset = count

		}

		opt := options.Find().SetSkip(count - int64(page)*int64(perPage)).SetLimit(int64(perPage))
		cursor, err1 := c.Find(context.TODO(), bson.M{"$and": []interface{}{bson.M{"identifier": identifier}, bson.M{"txntype": bson.M{"$in": []string{"0", "2", "5", "6", "10"}}}}}, opt)
		err2 := cursor.All(context.TODO(), &(result.Transactions))
		result.Count = int64(count)
		fmt.Println(count)
		if err1 != nil || err2 != nil || len(result.Transactions) == 0 {
			// fmt.Println(err1)
			reject(err1)
		} else {
			resolve(result)
		}
	})

	return p
}

func (cd *Connection) GetTestimonialbyStatus(status string) *promise.Promise {
	result := []model.Testimonial{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Testimonials")
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

func (cd *Connection) GetTestimonialOrganizationbyStatus(status string) *promise.Promise {
	result := []model.TestimonialOrganization{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Organizations")
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

// get proof protocol by proof name
func (cd *Connection) GetProofProtocolByProofName(proofName string) *promise.Promise {
	resultProtocolObj := model.ProofProtocol{}

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			log.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("ProofProtocols")
		err = c.FindOne(context.TODO(), bson.M{"proofname": proofName}).Decode(&resultProtocolObj)
		if err != nil {
			log.Error("Error when fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(resultProtocolObj)
		}
	})
	return p
}

// get all approved organizations with pagination
func (cd *Connection) GetAllApprovedOrganizations_Paginated(perPage int, page int) *promise.Promise {
	result := []model.TestimonialOrganization{}

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("Organizations")
		count, er := c.CountDocuments(context.TODO(), bson.M{"status": model.Approved.String()})

		if er != nil {
			log.Error("Error while get f.count " + err.Error())
			reject(er)
		}

		if count < int64(perPage) {
			cursor, err1 := c.Find(context.TODO(), bson.M{})
			err2 := cursor.All(context.TODO(), &result)

			if err1 != nil || err2 != nil || len(result) == 0 {
				log.Error("Error while f.All " + err1.Error())
				reject(err1)
			} else {
				resolve(result)
			}

		}

		offset := int64(page) * int64(perPage)
		if count < offset {
			pagelimit := perPage
			perPage = int(count + int64(perPage) - offset)
			if (offset - count) < int64(pagelimit) {
				offset = count
			}

		}

		collation := options.Collation{
			Locale:    "en",
			CaseLevel: true,
		}

		opts := options.Find().SetSort(bson.M{"name": -1}).SetSkip(count - offset).SetLimit(int64(perPage)).SetCollation(&collation)
		cursor, err1 := c.Find(context.TODO(), bson.M{"status": model.Approved.String()}, opts)

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

func (cd *Connection) GetRealIdentifier(mapValue string) *promise.Promise {
	result := apiModel.IdentifierModel{}
	fmt.Println("----------identifier 1")
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			log.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		fmt.Println("----------identifier 2")
		c := session.Client().Database(dbName).Collection("IdentifierMap")
		err = c.FindOne(context.TODO(), bson.M{"mapvalue": mapValue, "identifier": bson.M{"$ne": ""}}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(&result)
		if err != nil {
			log.Error("Error when fetching data from DB " + err.Error())
			reject(err)
		} else {
			if result.Identifier != "" {
				fmt.Println("----------identifier is there")
				resolve(result)
			} else {
				fmt.Println("----------identifier not there")
				result.MapValue = mapValue
				result.Identifier = mapValue
				resolve(result)
			}
		}
	})
	return p
}

func (cd *Connection) GetMapValue(Identifier string) *promise.Promise {
	result := apiModel.IdentifierModel{}

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			log.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("IdentifierMap")
		err = c.FindOne(context.TODO(), bson.M{"identifier": Identifier}).Decode(&result)
		if err != nil {
			log.Error("Error when fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetRealIdentifierByMapValue(identifier string) *promise.Promise {
	result := apiModel.IdentifierModel{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("IdentifierMap")
		err1 := c.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&result)
		if err1 != nil {
			reject(err1)
		} else {
			result1 := []model.TransactionCollectionBody{}
			session, err := cd.connect()
			if err != nil {
				reject(err)
			}
			defer session.EndSession(context.TODO())
			c := session.Client().Database(dbName).Collection("Transactions")
			cursor, err1 := c.Find(context.TODO(), bson.M{"identifier": result.MapValue})

			if err1 != nil {
				reject(err1)
			} else {
				err2 := cursor.All(context.TODO(), &result1)
				if err2 != nil || len(result1) == 0 {
					reject(err2)
				} else {
					resolve(result1)
				}
			}
		}
	})
	return p
}

func (cd *Connection) GetRealIdentifiersByArtifactId(identifier []string) *promise.Promise {
	// map identifiers objcet array
	var mapIdentifiers []apiModel.IdentifierModel
	// only map identifiers
	var mapIdentifiersArray []string
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session1, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session1.EndSession(context.TODO())
		c1 := session1.Client().Database(dbName).Collection("IdentifierMap")
		// find map identifiers using real identifers
		rst, err1 := c1.Find(context.TODO(), bson.D{{"identifier", bson.D{{"$in", identifier}}}})
		// check for errors in the finding
		if err1 != nil {
			reject(err1)
		}
		// read the douments and assign it to mapIdentifiers
		if err2 := rst.All(context.TODO(), &mapIdentifiers); err != nil {
			reject(err2)
		}

		// creating array using only map identifiers
		for _, result := range mapIdentifiers {
			mapIdentifiersArray = append(mapIdentifiersArray, result.MapValue)
		}

		result1 := []model.TransactionCollectionBody{}
		session2, err3 := cd.connect()
		if err3 != nil {
			reject(err)
		}
		defer session2.EndSession(context.TODO())
		c2 := session2.Client().Database(dbName).Collection("Transactions")
		// find transacions using mapIdentifers in gateway
		cursor, err4 := c2.Find(context.TODO(), bson.D{{"identifier", bson.D{{"$in", mapIdentifiersArray}}}})
		if err4 != nil {
			reject(err1)
		}
		err5 := cursor.All(context.TODO(), &result1)
		if err5 != nil || len(result1) == 0 {
			reject(err5)
		} else {
			resolve(result1)
		}
	})
	p.Await()
	return p
}

func (cd *Connection) GetNFTMinterPKSolana(ImageBase64 string, blockchain string) *promise.Promise {
	result := model.NFTWithTransactionSolana{}
	promise := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		dbclient := session.Client().Database(dbName).Collection("NFTSolana")
		err1 := dbclient.FindOne(context.TODO(), bson.D{{"imagebase64", ImageBase64}, {"nftissuingblockchain", blockchain}}).Decode(&result)
		if err1 != nil {
			reject(err)
		} else {
			resolve(result)
		}
	})
	return promise
}

func (cd *Connection) GetNFTTxnForStellar(ImageBase64 string, blockchain string) *promise.Promise {
	result := model.NFTWithTransaction{}
	promise := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		dbclient := session.Client().Database(dbName).Collection("NFTStellar")
		err1 := dbclient.FindOne(context.TODO(), bson.D{{"imagebase64", ImageBase64}, {"nftissuingblockchain", blockchain}}).Decode(&result)
		if err1 != nil {
			reject(err)
		} else {
			resolve(result)
		}
	})
	return promise
}

func (cd *Connection) GetAllSellingNFTStellar_Paginated(sellingStatus string, distributorPK string) *promise.Promise {
	result := model.MarketPlaceNFTTrasactionWithCount{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("MarketPlaceNFT")
		count, er := c.CountDocuments(context.TODO(), bson.M{"$and": []interface{}{bson.M{"sellingstatus": sellingStatus}}})
		if er != nil {
			log.Error("Error while get f.count " + err.Error())
			reject(er)
		}
		if distributorPK != "withoutKey" {
			cursor, er := c.Find(context.TODO(), bson.M{
				"currentownernftpk": distributorPK,
				"$or": []interface{}{
					bson.M{"sellingstatus": "FORSALE"},
					bson.M{"sellingstatus": "NOTFORSALE"},
				},
			})
			if er != nil {
				reject(er)
			} else {
				err2 := cursor.All(context.TODO(), &(result.MarketPlaceNFTItems))
				result.Count = int64(count)
				if err2 != nil || len(result.MarketPlaceNFTItems) == 0 {
					reject(err2)
				} else {
					resolve(result)
				}
			}
		} else {
			cursor, er := c.Find(context.TODO(), bson.M{"$and": []interface{}{bson.M{"sellingstatus": sellingStatus}}})
			if er != nil {
				reject(er)
			} else {
				err2 := cursor.All(context.TODO(), &(result.MarketPlaceNFTItems))
				result.Count = int64(count)
				if err2 != nil || len(result.MarketPlaceNFTItems) == 0 {
					reject(err2)
				} else {
					resolve(result)
				}
			}
		}
	})
	return p
}

func (cd *Connection) GetNFTByNFTTxn(NFTTXNhash string) *promise.Promise {
	result := model.MarketPlaceNFT{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			log.Error("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("MarketPlaceNFT")
		err1 := c.FindOne(context.TODO(), bson.M{"nfttxnhash": NFTTXNhash}).Decode(&result)
		if err1 != nil {
			log.Error("Error while getting NFT TXN Hash from db " + err.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetLastNFTbyInitialDistributorPK(InitialDistributorPK string) *promise.Promise {
	result := []model.MarketPlaceNFT{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("MarketPlaceNFT")
		cursor, err1 := c.Find(context.TODO(), bson.M{"initialdistributorpk": InitialDistributorPK})

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

func (cd *Connection) GetNFTIssuerSK(isserPK string) *promise.Promise {
	result := []model.NFTKeys{}
	promise := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		dbclient := session.Client().Database(dbName).Collection("NFTKeys")
		cursor, err := dbclient.Find(context.TODO(), bson.M{"publickey": isserPK})
		if err != nil {
			reject(err)
		} else {
			err := cursor.All(context.TODO(), &(result))
			if err != nil {
				reject(err)
			} else {
				resolve(result)
			}
		}
	})
	return promise
}

func (cd *Connection) GetFormulaMapID(formulaID string) *promise.Promise {
	result := model.FormulaIDMap{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("FormulaIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"formulaid": formulaID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting FormulaIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetExpertMapID(expertID string) *promise.Promise {
	result := model.ExpertIDMap{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("ExpertIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"expertpk": expertID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting ExpertIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetValueMapID(valueID, formulaId string) *promise.Promise {
	result := model.ValueIDMap{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("ValueIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"valueid": valueID, "formulaid": formulaId}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting ValueIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetUnitMapID(unit string) *promise.Promise {
	result := model.UnitIDMap{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("UnitIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"unit": unit}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting UnitIdMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetRequestAmount(reqEntityType string, reqEntity string, fromTime time.Time, toTime time.Time) *promise.Promise {
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("APIThrottleCounter")
		count, err1 := c.CountDocuments(context.TODO(), bson.M{"requestentitytype": reqEntityType, "requestentity": reqEntity, "timestamp": bson.M{"$gt": fromTime, "$lt": toTime}})
		if err1 != nil {
			log.Error("Error while retrieving request count from DB " + err1.Error())
			reject(err1)
		} else {
			resolve(count)
		}
	})
	return p
}

// for checking the public keys against trust nwtwork
func (cd *Connection) GetTrustNetworkKeyMap(publicKey string) *promise.Promise {
	result := model.Expert{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TrustNetwork")
		err1 := c.FindOne(context.TODO(), bson.M{"publickey": publicKey}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting trust network key map from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetVariableDetails(formulaID string, key string, variableName string) *promise.Promise {
	var valueDef model.ValueIDMap

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("ValueIDMap")
		err = c.FindOne(context.TODO(), bson.M{"formulaid": formulaID, "key": key, "valuename": variableName}).Decode(&valueDef)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(valueDef)
		}
	})
	return p
}

func (cd *Connection) GetResourceMapID(resourceID string) *promise.Promise {
	result := model.ResourceIdMap{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("ResourceIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"resourceid": resourceID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting ResourceIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetMetricMapID(metricID string) *promise.Promise {
	result := model.MetricMapDetails{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("MetricIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"metricid": metricID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting ResourceIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetTenentMapID(tenentUUID string) *promise.Promise {
	result := model.TenentMapDetails{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TenentIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"tenentid": tenentUUID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting tenent id from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetActivityMapID(activityId string) *promise.Promise {
	result := model.ActivityMapDetails{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("ActivityIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"activityid": activityId}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting activity id from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetExpertFormulaCount(formulaID string) *promise.Promise {
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("ExpertFormula")
		count, err1 := c.CountDocuments(context.TODO(), bson.M{"formulaid": formulaID, "status": "Success"})
		if err1 != nil {
			logrus.Info("Error while getting FormulaIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(count)
		}
	})
	return p
}

func (cd Connection) GetRSAPublicKeyBySHA256PK(sha256pk string) *promise.Promise {
	result := []model.RSAPublickey{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("PGPAccounts")
		cursor, err1 := c.Find(context.TODO(), bson.M{"pgppksha256": sha256pk})

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

func (cd *Connection) GetWorkflowMapID(workflowId string) *promise.Promise {
	result := model.WorkflowMap{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("WorkflowIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"workflowid": workflowId}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting workflow id from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetArtifactTemplateMapID(artifactId string) *promise.Promise {
	result := model.ArtifactTemplateId{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("ArtifactTemplateIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"artifactid": artifactId}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting artifact id from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetValueMapDetails(formulaID string, key string) *promise.Promise {
	var valueDef model.ValueIDMap

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("ValueIDMap")
		err = c.FindOne(context.TODO(), bson.M{"formulaid": formulaID, "key": key}).Decode(&valueDef)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(valueDef)
		}
	})
	return p
}

func (cd *Connection) GetMetricStatus(metricId string) *promise.Promise {
	var metricMap model.MetricBindingStore

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("MetricBinds")
		err = c.FindOne(context.TODO(), bson.M{"metricid": metricId}).Decode(&metricMap)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(metricMap)
		}
	})
	return p
}

func (cd *Connection) GetFormulaStatus(formulaID string) *promise.Promise {
	var metricMap model.FormulaStore

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("ExpertFormula")
		err = c.FindOne(context.TODO(), bson.M{"formulaid": formulaID}).Decode(&metricMap)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(metricMap)
		}
	})
	return p
}

func (cd *Connection) GetBindKey(formulaID, key, metricId string) *promise.Promise {
	var bindKey model.BindKeyMap
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("BindKeyMap")
		err = c.FindOne(context.TODO(), bson.M{"formulaid": formulaID, "keyinblockchain": key, "metricId": metricId}).Decode(&bindKey)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(bindKey)
		}
	})
	return p
}

func (cd *Connection) GetEthFormulaStatus(formulaID string) *promise.Promise {
	var ethFormulaMap model.EthereumExpertFormula
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum formula status by formula ID", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumExpertFormula")
		err = c.FindOne(context.TODO(), bson.M{"formulaid": formulaID}).Decode(&ethFormulaMap)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(ethFormulaMap)
		}
	})
	return p
}

func (cd *Connection) GetPrimaryKeyMapID(artifactId string) *promise.Promise {
	result := model.ArtifactTemplateId{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("PrimaryKeyIdMap")
		err1 := c.FindOne(context.TODO(), bson.M{"artifactid": artifactId}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting artifact id from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd Connection) GetEthFormulaMapID(formulaID string) *promise.Promise {
	result := model.EthFormulaIDMap{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum formula mapped ID by formula ID", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumFormulaIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"formulaid": formulaID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting FormulaIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetTrustNetWorkUserbyID(id primitive.ObjectID) *promise.Promise {
	result := []model.LoggedInTrustNetworkUser{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
		cursor, err1 := c.Find(context.TODO(), bson.M{"_id": id})
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

func (cd *Connection) GetTrustNetworkUserEndorsment(pkHash string) *promise.Promise {
	var useredorsments []model.TrustNetWorkUser
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
		cursor, err1 := c.Find(context.TODO(), bson.M{"pgppkhash": pkHash})
		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &useredorsments)
			if err2 != nil {
				reject(err2)
			} else {
				resolve(useredorsments[len(useredorsments)-1])
			}
		}
	})
	return p
}

func (cd *Connection) ValidateTrustNetworkUser(email string, password string) *promise.Promise {
	result := []model.LoggedInTrustNetworkUser{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
		cursor, err1 := c.Find(context.TODO(), bson.M{"email": email, "password": password})
		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || (result) == nil {
				reject(err2)
			} else {
				resolve(result[len(result)-1])
			}
		}
	})
	return p
}
func (cd *Connection) GetTrustNetworkUserResetPassword(email string, password []byte) *promise.Promise {
	result := []model.TrustNetWorkUser{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
		cursor, err1 := c.Find(context.TODO(), bson.M{"email": email})
		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &result)
			if err2 != nil || (result) == nil {
				reject(err2)
			} else {
				resolve(result[len(result)-1])
			}
		}
	})
	return p
}
func (cd *Connection) GetAllTrustNetworkUsers() *promise.Promise {
	result := []model.LoggedInTrustNetworkUser{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
		opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
		cursor, err1 := c.Find(context.TODO(), bson.M{}, opts)
		if err1 != nil {
			reject(err1)
		} else {
			for cursor.Next(context.TODO()) {
				var user model.LoggedInTrustNetworkUser
				decodeErr := cursor.Decode(&user)
				if decodeErr != nil {
					reject(decodeErr)
				}
				result = append(result, user)
			}
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetTrustNetWorkUserbyEncryptedPW(id string) *promise.Promise {
	var useredorsments []model.TrustNetWorkUser
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
		cursor, err1 := c.Find(context.TODO(), bson.M{"password": id})
		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &useredorsments)
			if err2 != nil {
				reject(err2)
			} else {
				resolve(useredorsments[len(useredorsments)-1])
			}
		}
	})
	return p
}

func (cd *Connection) GetTrustNetWorkUserbyEmail(email string) *promise.Promise {
	var pwdResetuser []model.TrustNetWorkUser
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
		cursor, err1 := c.Find(context.TODO(), bson.M{"email": email})
		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &pwdResetuser)
			if err2 != nil || (pwdResetuser) == nil {
				reject(err2)
			} else {
				resolve(pwdResetuser[len(pwdResetuser)-1])
			}
		}
	})
	return p
}

func (cd *Connection) GetPGPAccountByStellarPK(StellarPK string) *promise.Promise {
	result := model.PGPAccount{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("PGPAccounts")
		err1 := c.FindOne(context.TODO(), bson.M{"stellarpublickey": StellarPK}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting PGP Account from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetEthMetricStatus(metricID string) *promise.Promise {
	var ethMetric model.EthereumMetricBind

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum metric by metric ID", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("EthereumMetricBind")
		err = c.FindOne(context.TODO(), bson.M{"metricid": metricID}).Decode(&ethMetric)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(ethMetric)
		}
	})
	return p
}

func (cd *Connection) GetEthFormulaByName(formulaName string) *promise.Promise {
	result := model.EthereumExpertFormula{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum formula by name", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumExpertFormula")
		err1 := c.FindOne(context.TODO(), bson.M{"formulaname": formulaName, "status": "SUCCESS"}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting Ethereum formula by name from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) EthereumGetValueMapID(valueID, formulaId string) *promise.Promise {
	result := model.ValueIDMap{}
	// p := promise.NewPromise()
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum mapped value ID", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthValueIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"valueid": valueID, "formulaid": formulaId}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting ValueIDMap from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

// GET pass numbered identifiers and get real identifers
// mapvalues - numbered identifiers
// identifers - real identifier
func (cd *Connection) GetIdentifierMap(mapIdentifiers []string) *promise.Promise {
	var identifiers []apiModel.IdentifierModel
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("IdentifierMap")
		rst, err1 := c.Find(context.TODO(), bson.D{{"mapvalue", bson.D{{"$in", mapIdentifiers}}}})
		if err2 := rst.All(context.TODO(), &identifiers); err != nil {
			reject(err2)
		}
		if err1 != nil {
			reject(err1)
		} else {
			resolve(identifiers)
		}
	})
	return p
}

func (cd *Connection) GetEthereumMetricLatestContract(metricID string) *promise.Promise {
	result := model.MetricLatestContract{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum latest contract for the given metric ID", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthMetricLatest")
		err1 := c.FindOne(context.TODO(), bson.M{"metricid": metricID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting latest metric contract address from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetEthMetricByMetricIdAndType(metricID string, contractType string) *promise.Promise {
	result := model.EthereumMetricBind{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum metric by metric id and type", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumMetricBind")
		err1 := c.FindOne(context.TODO(), bson.M{"metricid": metricID, "type": contractType}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting the requested metric from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd Connection) GetEthMetricMapID(metricID string) *promise.Promise {
	result := model.EthMetricIDMap{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum mapped metric ID", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumMetricIDMap")
		err1 := c.FindOne(context.TODO(), bson.M{"metricid": metricID}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting Metric ID Map from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetEthMetricStatusForFormula(metricID string, contractType string, formulaID string) *promise.Promise {
	var ethMetric model.EthereumMetricBind

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Error("Error when connecting to DB " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum metric status for formula(activity)", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())

		c := session.Client().Database(dbName).Collection("EthereumMetricBind")
		err = c.FindOne(context.TODO(), bson.M{"metricid": metricID, "type": contractType, "formulaid": formulaID}).Decode(&ethMetric)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(ethMetric)
		}
	})
	return p
}

//Get the Ethereum contract status
func (cd *Connection) GetPendingContractsByStatus(status string) *promise.Promise {
	result := []model.PendingContracts{}

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			notificationhandler.InformDBConnectionIssue("get Ethereum pending contracts by status", err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumPendingTransactions")
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

func (cd *Connection) GetEthFormulaBinAndAbiByIdentifier(identifier string) *promise.Promise {
	result := model.EthereumExpertFormula{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum formula bin and abi by identifier", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumExpertFormula")
		err1 := c.FindOne(context.TODO(), bson.M{"transactionuuid": identifier}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting Ethereum formula contract abi and bin by uuid from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetEthMetricBinAndAbiByIdentifier(identifier string) *promise.Promise {
	result := model.EthereumMetricBind{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum metric bin and abi by identifier", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumMetricBind")
		err1 := c.FindOne(context.TODO(), bson.M{"transactionuuid": identifier}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting Ethereum metric contract abi and bin by uuid from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetEthMetricByUUID(identifier string) *promise.Promise {
	result := model.EthereumMetricBind{}
	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			logrus.Info("Error while connecting to db " + err.Error())
			notificationhandler.InformDBConnectionIssue("get Ethereum metric by UUID", err.Error())
			reject(err)
		}
		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumMetricBind")
		err1 := c.FindOne(context.TODO(), bson.M{"transactionuuid": identifier}).Decode(&result)
		if err1 != nil {
			logrus.Info("Error while getting Ethereum metric by uuid from db " + err1.Error())
			reject(err1)
		} else {
			resolve(result)
		}
	})
	return p
}

func (cd *Connection) GetEthMetricsByMetricID(metricID string) *promise.Promise {
	ethMetric := []model.EthereumMetricBind{}

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			notificationhandler.InformDBConnectionIssue("get Ethereum metrics by metric ID", err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumMetricBind")
		cursor, err1 := c.Find(context.TODO(), bson.M{"metricid": metricID})

		if err1 != nil {
			reject(err1)
		} else {
			err2 := cursor.All(context.TODO(), &ethMetric)
			if err2 != nil || len(ethMetric) == 0 {
				reject(err2)
			} else {
				resolve(ethMetric)
			}
		}
	})

	return p
}

func (cd *Connection) GetPendingContractByIdentifier(identifier string) *promise.Promise {
	result := model.PendingContracts{}

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			notificationhandler.InformDBConnectionIssue("get Ethereum pending contract by identifier", err.Error())
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("EthereumPendingTransactions")
		err = c.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&result)
		if err != nil {
			log.Info("Fetching data from DB " + err.Error())
			reject(err)
		} else {
			resolve(result)
		}
	})

	return p

}

/*
GetLastTransactionbyIdentifier Retrieve Last merger Transaction Object from TransactionCollection in DB by Identifier and merge block number
*/
func (cd *Connection) GetLastMergeTransactionbyIdentifierAndOrder(identifier, tenantId string, mergeBlock int) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"identifier": identifier, "txntype": "7", "mergeblock": mergeBlock, "tenantid": tenantId})

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

/*
GetLastTransactionbyIdentifier Retrieve Last Transaction Object from TransactionCollection in DB by Identifier and tenantId
*/
func (cd *Connection) GetLastTransactionbyIdentifierAndTenantId(identifier, tenantId string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("Transactions")
		cursor, err1 := c.Find(context.TODO(), bson.M{"identifier": identifier, "tenantid": tenantId})

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

func (cd *Connection) GetNFTById1Id2Id3(blockchain string, nftidentifier string, imagebase64 string) *promise.Promise {
	result := model.PendingNFTS{}
	// p := promise.NewPromise()

	p := promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()
		if err != nil {
			// fmt.Println(err)
			reject(err)
		}

		defer session.EndSession(context.TODO())
		c := session.Client().Database(dbName).Collection("NFTStatus")
		err1 := c.FindOne(context.TODO(), bson.M{"blockchain": blockchain, "nftidentifier": nftidentifier, "imagebase64": imagebase64}).Decode(&result)

		if err1 != nil {
			reject(err1)
		} else {
			resolve(result)
		}
	})

	return p
}
func (cd *Connection) GetQueueData(ImageBase64 string, blockchain string) *promise.Promise {
	result := model.PendingNFTS{}
	promise := promise.New(func(resolve func(interface{}), reject func(error)) {
		session, err := cd.connect()
		if err != nil {
			reject(err)
		}
		defer session.EndSession(context.TODO())
		dbclient := session.Client().Database(dbName).Collection("NFTStatus")
		err1 := dbclient.FindOne(context.TODO(), bson.D{{"imagebase64", ImageBase64}, {"blockchain", blockchain}}).Decode(&result)
		if err1 != nil {
			reject(err)
		} else {
			resolve(result)
		}
	})
	return promise
}
