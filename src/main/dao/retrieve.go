package dao

import (
	"main/model"

	"gopkg.in/mgo.v2/bson"

	"fmt"

	"github.com/chebyrash/promise"
)

func (cd *Connection) GetCOCbySender(sender string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"sender": sender}).All(&result)
		if err1 != nil || len(result) == 0 {
			fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

func (cd *Connection) GetCOCbyReceiver(receiver string) *promise.Promise {
	result := []model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"receiver": receiver}).All(&result)
		if err1 != nil || len(result) == 0 {
			fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

func (cd *Connection) GetCOCbyAcceptTxn(accepttxn string) *promise.Promise {
	result := model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"accepttxn": accepttxn}).One(&result)
		if err1 != nil {
			fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

func (cd *Connection) GetCOCbyRejectTxn(rejecttxn string) *promise.Promise {
	result := model.COCCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Find(bson.M{"rejecttxn": rejecttxn}).One(&result)
		if err1 != nil {
			fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}


func (cd *Connection) GetLastTransactionbyIdentifier(identifier string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"identifier": identifier}).All(&result)
		if err1 != nil || len(result) == 0 {
			fmt.Println(err1)
			reject(err1)

		}
		resolve(result[len(result)-1])

	})

	return p

}

func (cd *Connection) GetTransactionsbyIdentifier(identifier string) *promise.Promise {
	result := []model.TransactionCollectionBody{}
	// p := promise.NewPromise()

	var p = promise.New(func(resolve func(interface{}), reject func(error)) {
		// Do something asynchronously.
		session, err := cd.connect()

		if err != nil {
			fmt.Println(err)
			reject(err)

		}
		defer session.Close()

		c := session.DB("tracified-gateway").C("Transactions")
		err1 := c.Find(bson.M{"identifier": identifier}).All(&result)
		if err1 != nil || len(result) == 0 {
			fmt.Println(err1)
			reject(err1)

		}
		resolve(result)

	})

	return p

}

func (cd *Connection) GetTransactionForTdpId(TdpId string) *promise.Promise {
		result := model.TransactionCollectionBody{}
		// p := promise.NewPromise()
	
		var p = promise.New(func(resolve func(interface{}), reject func(error)) {
			// Do something asynchronously.
			session, err := cd.connect()
	
			if err != nil {
				fmt.Println(err)
				reject(err)
	
			}
			defer session.Close()
	
			c := session.DB("tracified-gateway").C("Transactions")
			err1 := c.Find(bson.M{"tdpid": TdpId}).One(&result)
			if err1 != nil {
				fmt.Println(err1)
				reject(err1)
	
			}
			resolve(result)
	
		})
	
		return p
	
	}

