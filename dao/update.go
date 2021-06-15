package dao

import (
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
)

/*UpdateTransaction  Update a Transaction Object from TransactionCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) UpdateTransaction(selector model.TransactionCollectionBody, update model.TransactionCollectionBody) error {
	fmt.Println("----------------------------------- UpdateTransaction ---------------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session "+err.Error())
		return err
	}
	defer session.Close()
	Previous:=selector.PreviousTxnHash
	if update.PreviousTxnHash!=""{
		Previous=update.PreviousTxnHash
	}
	up := model.TransactionCollectionBody{
		Identifier: selector.Identifier,
		TdpId:      selector.TdpId,
		PublicKey:  selector.PublicKey,
		XDR:        selector.XDR,
		TxnHash:    update.TxnHash,
		TxnType: selector.TxnType,
		Status:  update.Status,
		ProfileID:update.ProfileID,
		PreviousTxnHash:Previous,
		FromIdentifier1:selector.FromIdentifier1,
		FromIdentifier2:selector.FromIdentifier2,
		ItemAmount:selector.ItemAmount,
		ItemCode:selector.ItemCode,
	}
	c := session.DB(dbName).C("Transactions")
	err = c.Update(selector, up)
	if err != nil {
		fmt.Println("Error while updating Transactions "+err.Error())
	}
	return err
}

/*UpdateCOC Update a COC Object from COCCollection in DB on the basis of the status
@author - Azeem Ashraf
*/
func (cd *Connection) UpdateCOC(selector model.COCCollectionBody, update model.COCCollectionBody) error {
	fmt.Println("----------------------------------- UpdateCOC ---------------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session "+err.Error())
		return err
	}
	defer session.Close()
	fmt.Println(update.Status)
	switch update.Status {
	case "accepted":
		up := model.COCCollectionBody{
			TxnHash:    update.TxnHash,
			Sender:     selector.Sender,
			Receiver:   selector.Receiver,
			AcceptXdr:  update.AcceptXdr,
			RejectXdr:  selector.RejectXdr,
			AcceptTxn:  update.AcceptTxn,
			RejectTxn:  selector.RejectTxn,
			Identifier: selector.Identifier,
			Status:     update.Status,
			SubAccount: selector.SubAccount,
			SequenceNo: selector.SequenceNo,
		}
		c := session.DB(dbName).C("COC")
		err = c.Update(selector, up)
		if err != nil {
			fmt.Println("Error while updating COC case accepted"+err.Error())
			return err
		}
		break
	case "rejected":
		up := model.COCCollectionBody{
			TxnHash:    update.TxnHash,
			Sender:     selector.Sender,
			Receiver:   selector.Receiver,
			AcceptXdr:  selector.AcceptXdr,
			RejectXdr:  update.RejectXdr,
			AcceptTxn:  selector.AcceptTxn,
			RejectTxn:  update.RejectTxn,
			Identifier: selector.Identifier,
			Status:     update.Status,
			SubAccount: selector.SubAccount,
			SequenceNo: selector.SequenceNo,
		}
		c := session.DB(dbName).C("COC")
		err = c.Update(selector, up)
		if err != nil {
			fmt.Println("Error while updating COC case rejected"+err.Error())
			return err
		}
		break
	case "expired":
		up := model.COCCollectionBody{
			TxnHash:    selector.TxnHash,
			Sender:     selector.Sender,
			Receiver:   selector.Receiver,
			AcceptXdr:  selector.AcceptXdr,
			RejectXdr:  selector.RejectXdr,
			AcceptTxn:  selector.AcceptTxn,
			RejectTxn:  selector.RejectTxn,
			Identifier: selector.Identifier,
			Status:     update.Status,
			SubAccount: selector.SubAccount,
			SequenceNo: selector.SequenceNo,
		}
		c := session.DB(dbName).C("COC")
		err = c.Update(selector, up)
		if err != nil {
			fmt.Println("Error while updating COC case expired "+err.Error())
			return err
		}
		break
	}
	return err
}

/*UpdateCertificate Update a Certificate Object from CertificateCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) UpdateCertificate(selector model.TransactionCollectionBody, update model.TransactionCollectionBody) error {
	fmt.Println("----------------------------------- UpdateCertificate ---------------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session "+err.Error())
		return err
	}
	defer session.Close()
	c := session.DB(dbName).C("Certificates")
	err = c.Update(selector, update)
	if err != nil {
		fmt.Println("Error while updating certificates "+err.Error())
	}
	return err
}