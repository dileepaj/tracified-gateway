package dao

import (
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
)

//InsertCoc ...
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

//InsertTransaction ...
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

//InsertToOrphan ...
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

//InsertProfile ...
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

func (cd *Connection) InsertCertificate(Cert model.CertificateCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("Certificates")
	err1 := c.Insert(Cert)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}