package dao

import (
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
)

func (cd *Connection) InsertCoc(Coc model.COCCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("COC")
	err1 := c.Insert(Coc)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}

func (cd *Connection) InsertTransaction(Coc model.TransactionCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("Transactions")
	err1 := c.Insert(Coc)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}

func (cd *Connection) InsertToOrphan(Coc model.TransactionCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("Orphan")
	err1 := c.Insert(Coc)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}

func (cd *Connection) InsertProfile(Coc model.ProfileCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("Profiles")
	err1 := c.Insert(Coc)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}