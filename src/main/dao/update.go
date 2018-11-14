package dao

import (
	"main/model"
	"fmt"

)

func (cd *Connection)UpdateTransaction(selector model.TransactionCollectionBody,update model.TransactionCollectionBody) error {
	
	session, err :=cd.connect()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer session.Close()

	up:=model.TransactionCollectionBody{
		Identifier:selector.Identifier,
		TdpID:selector.TdpID,
		PublicKey:selector.PublicKey,
		XDR:selector.XDR,
		TxnHash:update.TxnHash,
		// ProfileHash:update.ProfileHash,
		TxnType:selector.TxnType,
		Status:update.Status,
		// ProfileID:update.ProfileID,
	}

	c := session.DB("tracified-gateway").C("Transactions")
	err1 := c.Update(selector,up)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}
