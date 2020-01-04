package dao

import (
	"fmt"
	"strconv"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"

	"gopkg.in/mgo.v2/bson"

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
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"sender": sender}).All(&result)
		if err1 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

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
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
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
		result2.Receiver = result.Receiver
		bumpSeq, err := strconv.Atoi(result.SequenceNo)
		if err == nil {
			fmt.Println(bumpSeq)
			bumpSeq = bumpSeq
			fmt.Println(bumpSeq)
		}
		result2.SequenceNo = strconv.Itoa(bumpSeq)
		result2.SubAccount = result.SubAccount
		if result.Status == "pending" {
			result2.Available = false
		} else {
			result2.Available = true
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"receiver": receiver}).All(&result)

		if err1 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

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
			// fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"accepttxn": accepttxn}).One(&result)
		if err1 != nil {
			// fmt.Println(err1)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"rejecttxn": rejecttxn}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"status": status}).All(&result)
		if err1 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

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
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		count, er := c.Find(bson.M{"identifier": identifier}).Count()
		if er != nil {
			// fmt.Println(er)
			reject(er)
		}

		err1 := c.Find(bson.M{"identifier": identifier}).Skip(count - 1).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		er := c.Find(bson.M{"txnhash": txnHash}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"identifier": identifier}).All(&result)
		if err1 != nil || len(result) == 0 {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result[len(result)-1])

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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"identifier": identifier}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"identifier": identifier}).All(&result)
		if err1 != nil || len(result) == 0 {
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
			// fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"tdpid": TdpId}).One(&result)
		if err1 != nil {
			// fmt.Println(err1)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"identifier": Identifer, "tdpid": ""}).One(&result)
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
			// fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"tdpid": TdpId}).All(&result)
		if err1 != nil {
			// fmt.Println(err1)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"txnhash": Txn}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Orphan")
		err1 := c.Find(bson.M{"identifier": identifier}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Profiles")
		err1 := c.Find(bson.M{"identifier": identifier}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Profiles")
		err1 := c.Find(bson.M{"profileid": ProfileID}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Certificates")
		err1 := c.Find(bson.M{"publickey": PublicKey}).All(&result)
		if err1 != nil || len(result) == 0 {
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Certificates")
		err1 := c.Find(bson.M{"certificateid": CertificateID}).All(&result)
		if err1 != nil || len(result) == 0 {
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Certificates")
		err1 := c.Find(bson.M{"publickey": PublicKey}).All(&result)
		if err1 != nil {
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"tdpid": tdpid}).All(&result)
		fmt.Println(result)

		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"tdpid": TdpId, "identifer": identifer}).One(&result)
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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"publickey": Publickey}).All(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

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
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"txnhash": Txnhash}).All(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

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
		defer session.Close()

		c := session.DB("tracified-gateway").C("TempOrphan")
		err1 := c.Find(bson.M{"publickey": Publickey, "sequenceno": SequenceNo}).One(&result)
		if err1 != nil {
			// fmt.Println(err1)
			reject(err1)

		} else {
			resolve(result)

		}

	})

	return p

}
