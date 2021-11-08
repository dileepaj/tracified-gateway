package dao

import (
	"context"
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
	"go.mongodb.org/mongo-driver/bson"
)

/*UpdateTransaction  Update a Transaction Object from TransactionCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) UpdateTransaction(selector model.TransactionCollectionBody, update model.TransactionCollectionBody) error {
	fmt.Println("----------------------------------- UpdateTransaction ---------------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
		return err
	}

	defer session.EndSession(context.TODO())

	Previous := selector.PreviousTxnHash
	if update.PreviousTxnHash != "" {
		Previous = update.PreviousTxnHash
	}
	up := model.TransactionCollectionBody{
		Identifier:      selector.Identifier,
		TdpId:           selector.TdpId,
		PublicKey:       selector.PublicKey,
		XDR:             selector.XDR,
		TxnHash:         update.TxnHash,
		TxnType:         selector.TxnType,
		Status:          update.Status,
		ProfileID:       update.ProfileID,
		PreviousTxnHash: Previous,
		FromIdentifier1: selector.FromIdentifier1,
		FromIdentifier2: selector.FromIdentifier2,
		ItemAmount:      selector.ItemAmount,
		ItemCode:        selector.ItemCode,
	}
	c := session.Client().Database(dbName).Collection("Transactions")

	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}

	var filter bson.M
	err = bson.Unmarshal(pByte, &filter)
	if err != nil {
		return err
	}

	pByte, err = bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	_, err = c.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: updateNew}})
	if err != nil {
		fmt.Println("Error while updating Transactions " + err.Error())
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
		fmt.Println("Error while getting session " + err.Error())
		return err
	}

	defer session.EndSession(context.TODO())

	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}

	var filter bson.M
	err = bson.Unmarshal(pByte, &filter)
	if err != nil {
		return err
	}

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

		pByte, err := bson.Marshal(up)
		if err != nil {
			return err
		}

		var update bson.M
		err = bson.Unmarshal(pByte, &update)
		if err != nil {
			return err
		}

		c := session.Client().Database(dbName).Collection("COC")
		_, err = c.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})
		if err != nil {
			fmt.Println("Error while updating COC case accepted" + err.Error())
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

		pByte, err := bson.Marshal(up)
		if err != nil {
			return err
		}

		var update bson.M
		err = bson.Unmarshal(pByte, &update)
		if err != nil {
			return err
		}

		c := session.Client().Database(dbName).Collection("COC")
		_, err = c.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})

		if err != nil {
			fmt.Println("Error while updating COC case rejected" + err.Error())
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

		pByte, err := bson.Marshal(up)
		if err != nil {
			return err
		}

		var update bson.M
		err = bson.Unmarshal(pByte, &update)
		if err != nil {
			return err
		}

		c := session.Client().Database(dbName).Collection("COC")
		_, err = c.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})

		if err != nil {
			fmt.Println("Error while updating COC case expired " + err.Error())
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
		fmt.Println("Error while getting session " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}

	var filter bson.M
	err = bson.Unmarshal(pByte, &filter)
	if err != nil {
		return err
	}

	pByte, err = bson.Marshal(update)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("Certificates")
	_, err = c.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: updateNew}})

	if err != nil {
		fmt.Println("Error while updating certificates " + err.Error())
	}
	return err
}

func (cd *Connection) UpdateOrganization(selector model.TestimonialOrganization, update model.TestimonialOrganization) error {
	fmt.Println("----------------------------------- UpdateOrganization---------------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.TestimonialOrganization{
		Name:           selector.Name,
		Description:    selector.Description,
		Logo:           selector.Logo,
		Email:          selector.Email,
		Phone:          selector.Phone,
		PhoneSecondary: selector.PhoneSecondary,
		AcceptTxn:      selector.AcceptTxn,
		AcceptXDR:      update.AcceptXDR,
		RejectTxn:      selector.RejectTxn,
		RejectXDR:      update.RejectXDR,
		TxnHash:        update.TxnHash,
		Author:         selector.Author,
		SubAccount:     selector.SubAccount,
		Status:         update.Status,
		ApprovedBy:     update.ApprovedBy,
		ApprovedOn:     update.ApprovedOn,
	}

	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}

	var filter bson.M
	err = bson.Unmarshal(pByte, &filter)
	if err != nil {
		return err
	}

	pByte, err = bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("Organizations")
	_, err = c.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: updateNew}})

	if err != nil {
		fmt.Println("Error while updating Organization " + err.Error())
	}
	return err
}

func (cd *Connection) UpdateTestimonial(selector model.Testimonial, update model.Testimonial) error {
	fmt.Println("----------------------------------- UpdateTestimonial ---------------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
		return err
	}

	defer session.EndSession(context.TODO())

	up := model.Testimonial{
		Sender:      selector.Sender,
		Reciever:    selector.Reciever,
		AcceptTxn:   selector.AcceptTxn,
		RejectTxn:   selector.RejectTxn,
		AcceptXDR:   update.AcceptXDR,
		RejectXDR:   update.RejectXDR,
		TxnHash:     update.TxnHash,
		Subaccount:  selector.Subaccount,
		Status:      update.Status,
		Testimonial: selector.Testimonial,
	}

	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}

	var filter bson.M
	err = bson.Unmarshal(pByte, &filter)
	if err != nil {
		return err
	}

	pByte, err = bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("Testimonials")
	_, err = c.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: updateNew}})

	if err != nil {
		fmt.Println("Error while updating Testimonials " + err.Error())
	}
	return err
}
