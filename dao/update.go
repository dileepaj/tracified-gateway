package dao

import (
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
)

func (cd *Connection) UpdateTransaction(selector model.TransactionCollectionBody, update model.TransactionCollectionBody) error {

	session, err := cd.connect()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer session.Close()

	up := model.TransactionCollectionBody{
		Identifier: selector.Identifier,
		TdpId:      selector.TdpId,
		PublicKey:  selector.PublicKey,
		XDR:        selector.XDR,
		TxnHash:    update.TxnHash,
		// ProfileHash:update.ProfileHash,
		TxnType: selector.TxnType,
		Status:  update.Status,
		// ProfileID:update.ProfileID,
	}

	c := session.DB("tracified-gateway").C("Transactions")
	err1 := c.Update(selector, up)
	if err1 != nil {
		fmt.Println(err1)
	}

	return err
}

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
