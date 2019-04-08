package dao

import (
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
)

/*UpdateTransaction  Update a Transaction Object from TransactionCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) UpdateTransaction(selector model.TransactionCollectionBody, update model.TransactionCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
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

	c := session.DB("tracified-gateway").C("Transactions")
	err1 := c.Update(selector, up)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}

/*UpdateCOC Update a COC Object from COCCollection in DB on the basis of the status
@author - Azeem Ashraf
*/
func (cd *Connection) UpdateCOC(selector model.COCCollectionBody, update model.COCCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
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

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Update(selector, up)
		if err1 != nil {
			fmt.Println(err1)
			return err1

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

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Update(selector, up)
		if err1 != nil {
			fmt.Println(err1)
			return err1

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

		c := session.DB("tracified-gateway").C("COC")
		err1 := c.Update(selector, up)
		if err1 != nil {
			fmt.Println(err1)
			return err1
		}

		break
	}

	return err
}

/*UpdateCertificate Update a Certificate Object from CertificateCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) UpdateCertificate(selector model.TransactionCollectionBody, update model.TransactionCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer session.Close()

	c := session.DB("tracified-gateway").C("Certificates")
	err1 := c.Update(selector, update)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}