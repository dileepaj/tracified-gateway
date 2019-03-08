package dao

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"

	"gopkg.in/mgo.v2/bson"

	// "fmt"

	"github.com/chebyrash/promise"
)

//GetCOCbySender ...
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

//GetLastCOCbySubAccount ...
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
		result2.SequenceNo = result.SequenceNo
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

//GetCOCbyReceiver ...
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

//GetCOCbyAcceptTxn ...
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
//GetCOCbyRejectTxn ...
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

//GetCOCbyStatus ...
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

//GetLastCOCbyIdentifier ...
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
//GetLastTransactionbyIdentifier ...
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

//GetFirstTransactionbyIdentifier ...
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

//GetTransactionsbyIdentifier ...
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

//GetTransactionForTdpId ...
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

//GetTdpIdForTransaction ...simply returns a TXN collection when a txnid is given
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

//GetOrphanbyIdentifier ...
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

//GetProfilebyIdentifier ...
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

