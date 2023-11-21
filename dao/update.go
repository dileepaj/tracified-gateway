package dao

import (
	"context"
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/model"
	notificationhandler "github.com/dileepaj/tracified-gateway/services/notificationHandler.go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
UpdateTransaction  Update a Transaction Object from TransactionCollection in DB
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

/*
UpdateCOC Update a COC Object from COCCollection in DB on the basis of the status
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

/*
UpdateCertificate Update a Certificate Object from CertificateCollection in DB
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

func (cd *Connection) Updateorganization(selector model.TestimonialOrganization, update model.TestimonialOrganization) error {
	fmt.Println("----------------------------------- UpdateOrganization---------------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())
	switch update.Status {
	case model.Approved.String():
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
			RejectXDR:      selector.RejectXDR,
			TxnHash:        update.TxnHash,
			Author:         selector.Author,
			SubAccount:     selector.SubAccount,
			SequenceNo:     selector.SequenceNo,
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
		break

	case model.Rejected.String():
		up := model.TestimonialOrganization{
			Name:           selector.Name,
			Description:    selector.Description,
			Logo:           selector.Logo,
			Email:          selector.Email,
			Phone:          selector.Phone,
			PhoneSecondary: selector.PhoneSecondary,
			AcceptTxn:      selector.AcceptTxn,
			AcceptXDR:      selector.AcceptXDR,
			RejectTxn:      selector.RejectTxn,
			RejectXDR:      update.RejectXDR,
			TxnHash:        update.TxnHash,
			Author:         selector.Author,
			SubAccount:     selector.SubAccount,
			SequenceNo:     selector.SequenceNo,
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
		break

	case model.Expired.String():
		up := model.TestimonialOrganization{
			Name:           selector.Name,
			Description:    selector.Description,
			Logo:           selector.Logo,
			Email:          selector.Email,
			Phone:          selector.Phone,
			PhoneSecondary: selector.PhoneSecondary,
			AcceptTxn:      selector.AcceptTxn,
			AcceptXDR:      selector.AcceptXDR,
			RejectTxn:      selector.RejectTxn,
			RejectXDR:      update.RejectXDR,
			TxnHash:        update.TxnHash,
			Author:         selector.Author,
			SubAccount:     selector.SubAccount,
			SequenceNo:     selector.SequenceNo,
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
		break
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
	switch update.Status {
	case model.Approved.String():
		up := model.Testimonial{
			Sender:      selector.Sender,
			Reciever:    selector.Reciever,
			AcceptTxn:   selector.AcceptTxn,
			RejectTxn:   selector.RejectTxn,
			AcceptXDR:   update.AcceptXDR,
			RejectXDR:   selector.RejectXDR,
			TxnHash:     update.TxnHash,
			Subaccount:  selector.Subaccount,
			SequenceNo:  selector.SequenceNo,
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
		break

	case model.Rejected.String():
		up := model.Testimonial{
			Sender:      selector.Sender,
			Reciever:    selector.Reciever,
			AcceptTxn:   selector.AcceptTxn,
			RejectTxn:   selector.RejectTxn,
			AcceptXDR:   selector.AcceptXDR,
			RejectXDR:   update.RejectXDR,
			TxnHash:     update.TxnHash,
			SequenceNo:  selector.SequenceNo,
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
		break

	case model.Expired.String():
		up := model.Testimonial{
			Sender:      selector.Sender,
			Reciever:    selector.Reciever,
			AcceptTxn:   selector.AcceptTxn,
			RejectTxn:   selector.RejectTxn,
			AcceptXDR:   selector.AcceptXDR,
			RejectXDR:   selector.RejectXDR,
			TxnHash:     selector.TxnHash,
			SequenceNo:  selector.SequenceNo,
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
		break
	}

	return err
}

func (cd *Connection) UpdateProofPresesntationProtocol(selector model.ProofProtocol, update model.ProofProtocol) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.ProofProtocol{
		ProofName:            update.ProofName,
		ProofDescriptiveName: update.ProofDescriptiveName,
		NumberofSteps:        update.NumberofSteps,
		Steps:                update.Steps,
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

	c := session.Client().Database(dbName).Collection("ProofProtocols")
	_, err = c.UpdateOne(context.TODO(), bson.M{"proofname": selector.ProofName}, bson.D{{Key: "$set", Value: updateNew}})

	if err != nil {
		fmt.Println("Error while updating proof protocols " + err.Error())
	}
	return err
}

func (cd *Connection) UpdateBuyingStatus(selector model.MarketPlaceNFT, updateStatus string, updateCurrentPK string, updatePreviousPK string) error {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error while getting session " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())
	up := model.MarketPlaceNFT{
		Identifier:                       selector.Identifier,
		Collection:                       selector.Collection,
		Categories:                       selector.Categories,
		ImageBase64:                      selector.ImageBase64,
		NftTransactionExistingBlockchain: selector.NftTransactionExistingBlockchain,
		NftIssuingBlockchain:             selector.NftIssuingBlockchain,
		NFTTXNhash:                       selector.NFTTXNhash,
		Timestamp:                        selector.Timestamp,
		NftURL:                           selector.NftURL,
		NftContentName:                   selector.NftContentName,
		NftContent:                       selector.NftContent,
		NFTArtistName:                    selector.NFTArtistName,
		NFTArtistURL:                     selector.NFTArtistURL,
		InitialIssuerPK:                  selector.InitialIssuerPK,
		InitialDistributorPK:             selector.InitialDistributorPK,
		TrustLineCreatedAt:               selector.TrustLineCreatedAt,
		MainAccountPK:                    selector.MainAccountPK,
		Description:                      selector.Description,
		Copies:                           selector.Copies,
		PreviousOwnerNFTPK:               updatePreviousPK,
		CurrentOwnerNFTPK:                updateCurrentPK,
		OriginPK:                         updateCurrentPK,
		SellingStatus:                    updateStatus,
		Amount:                           selector.Amount,
		Price:                            selector.Price,
	}
	c := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}
	var filter bson.D
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
	_, err = c.UpdateOne(context.TODO(), bson.M{"nfttxnhash": selector.NFTTXNhash}, bson.D{{Key: "$set", Value: updateNew}})
	if err != nil {
		log.Println("Error while updating NFT Stellar " + err.Error())
	}
	return err
}

func (cd *Connection) UpdateSellingStatus(selector model.MarketPlaceNFT, updateStatus string, updateAmount string, updatePrice string) error {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error while getting session " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())
	up := model.MarketPlaceNFT{
		Identifier:                       selector.Identifier,
		Collection:                       selector.Collection,
		Categories:                       selector.Categories,
		ImageBase64:                      selector.ImageBase64,
		NftTransactionExistingBlockchain: selector.NftTransactionExistingBlockchain,
		NftIssuingBlockchain:             selector.NftIssuingBlockchain,
		NFTTXNhash:                       selector.NFTTXNhash,
		Timestamp:                        selector.Timestamp,
		NftURL:                           selector.NftURL,
		NftContentName:                   selector.NftContentName,
		NftContent:                       selector.NftContent,
		NFTArtistName:                    selector.NFTArtistName,
		NFTArtistURL:                     selector.NFTArtistURL,
		InitialIssuerPK:                  selector.InitialIssuerPK,
		InitialDistributorPK:             selector.InitialDistributorPK,
		TrustLineCreatedAt:               selector.TrustLineCreatedAt,
		MainAccountPK:                    selector.MainAccountPK,
		Description:                      selector.Description,
		Copies:                           selector.Copies,
		PreviousOwnerNFTPK:               selector.PreviousOwnerNFTPK,
		CurrentOwnerNFTPK:                selector.CurrentOwnerNFTPK,
		OriginPK:                         selector.OriginPK,
		SellingStatus:                    updateStatus,
		Amount:                           updateAmount,
		Price:                            updatePrice,
	}
	c := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	pByte, err := bson.Marshal(selector)
	if err != nil {
		return err
	}
	var filter bson.D
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
	_, err = c.UpdateOne(context.TODO(), bson.M{"nfttxnhash": selector.NFTTXNhash}, bson.D{{Key: "$set", Value: updateNew}})
	if err != nil {
		log.Println("Error while updating NFT Stellar " + err.Error())
	}
	return err
}

// auto count sequence incrementer
func (cd *Connection) GetNextSequenceValue(Id string) (model.Counters, error) {
	var result model.Counters
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return model.Counters{}, err
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("Counters")
	err = c.FindOneAndUpdate(
		context.TODO(),
		bson.M{"id": Id}, // <- Find block
		bson.D{{"$inc", bson.D{{"sequencevalue", 1}}}},
		options.FindOneAndUpdate().SetReturnDocument(options.After), options.FindOneAndUpdate().SetUpsert(true), // <- Set option to return document after update (important)
	).Decode(&result)
	if err != nil {
		fmt.Println("Error while updating proof protocols " + err.Error())
	}
	return result, err
}

func (cd *Connection) UpdateCounterOnThrottler(ID primitive.ObjectID, newIndex int) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return err
	}
	c := session.Client().Database(dbName).Collection("APIThrottleCounter")
	// id, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.D{{"_id", ID}}
	update := bson.D{{"$set", bson.D{{"currentcount", newIndex}}}}
	_, errWhenUpdate := c.UpdateOne(context.TODO(), filter, update)
	if errWhenUpdate != nil {
		logrus.Error("Error when updating the throttler counter in DB : " + errWhenUpdate.Error())
		return errWhenUpdate
	}
	return nil
}

func (cd *Connection) UpdateOrganizationInfo(data model.TestimonialOrganization) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return err
	}
	c := session.Client().Database(dbName).Collection("Organizations")
	filter := bson.D{{"sequenceno", data.SequenceNo}}
	update := bson.D{{"$set", bson.D{{"pgpdata", data.PGPData}, {"status", data.Status}, {"sequenceno", data.SequenceNo}}}}
	_, errWhenUpdate := c.UpdateOne(context.TODO(), filter, update)
	if errWhenUpdate != nil {
		log.Println(err.Error())
		return errWhenUpdate
	}
	return nil
	//end of function
}

func (cd *Connection) UpdateMetricBindStatus(metricID string, txnUUID string, update model.MetricBindingStore) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.MetricBindingStore{
		MetricId:              update.MetricId,
		MetricMapID:           update.MetricMapID,
		Metric:                update.Metric,
		User:                  update.User,
		TotalNoOfManageData:   update.TotalNoOfManageData,
		NoOfManageDataInTxn:   update.NoOfManageDataInTxn,
		TransactionTime:       update.TransactionTime,
		TransactionCost:       update.TransactionCost,
		Memo:                  update.Memo,
		TxnHash:               update.TxnHash,
		TxnSenderPK:           update.TxnSenderPK,
		XDR:                   update.XDR,
		SequenceNo:            update.SequenceNo,
		Status:                update.Status,
		Timestamp:             update.Timestamp,
		ErrorMessage:          update.ErrorMessage,
		TxnUUID:               update.TxnUUID,
		ActivityManageDataMap: update.ActivityManageDataMap,
		TransactionOrderCount: update.TransactionOrderCount,
	}

	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("MetricBinds")
	_, err = c.UpdateOne(context.TODO(), bson.M{"metricid": metricID, "txnuuid": txnUUID}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}

func (cd *Connection) UpdateFormulaStatus(formulaID string, txnUUID string, update model.FormulaStore) error {
	logrus.Info("--------", formulaID, "---", txnUUID)
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.FormulaStore{
		MetricExpertFormula: update.MetricExpertFormula,
		Verify:              update.Verify,
		FormulaID:           update.FormulaID,
		FormulaMapID:        update.FormulaMapID,
		VariableCount:       update.VariableCount,
		ExecutionTemplate:   update.ExecutionTemplate,
		TotalNoOfManageData: update.TotalNoOfManageData,
		NoOfManageDataInTxn: update.NoOfManageDataInTxn,
		Memo:                update.Memo,
		TxnHash:             update.TxnHash,
		TxnSenderPK:         update.TxnSenderPK,
		XDR:                 update.XDR,
		SequenceNo:          update.SequenceNo,
		Status:              update.Status,
		Timestamp:           update.Timestamp,
		TransactionTime:     update.TransactionTime,
		TransactionCost:     update.TransactionCost,
		ErrorMessage:        update.ErrorMessage,
		TxnUUID:             update.TxnUUID,
	}

	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("ExpertFormula")
	_, err = c.UpdateOne(context.TODO(), bson.M{"formulaid": formulaID, "txnuuid": txnUUID}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}

func (cd *Connection) UpdateEthereumFormulaStatus(formulaID string, txnUUID string, update model.EthereumExpertFormula) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum expert formula", err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.EthereumExpertFormula{
		FormulaID:           update.FormulaID,
		FormulaName:         update.FormulaName,
		ExecutionTemplate:   update.ExecutionTemplate,
		MetricExpertFormula: update.MetricExpertFormula,
		VariableCount:       update.VariableCount,
		TemplateString:      update.TemplateString,
		BINstring:           update.BINstring,
		ABIstring:           update.ABIstring,
		SetterNames:         update.SetterNames,
		ContractAddress:     update.ContractAddress,
		ContractName:        update.ContractName,
		Status:              update.Status,
		Timestamp:           update.Timestamp,
		TransactionHash:     update.TransactionHash,
		TransactionCost:     update.TransactionCost,
		TransactionUUID:     update.TransactionUUID,
		TransactionSender:   update.TransactionSender,
		Verify:              update.Verify,
		ErrorMessage:        update.ErrorMessage,
		ActualStatus:        update.ActualStatus,
	}

	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("EthereumExpertFormula")
	_, err = c.UpdateOne(context.TODO(), bson.M{"formulaid": formulaID, "transactionuuid": txnUUID}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}

func (cd *Connection) UpdateTrustNetworkUserEndorsment(pkhash string, update model.TrustNetWorkUser) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Println("Error while connecting to DB " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())
	up := model.TrustNetWorkUser{
		Name:               update.Name,
		Company:            update.Company,
		Email:              update.Email,
		Password:           update.Password,
		Contact:            update.Contact,
		Industry:           update.Industry,
		StellerPK:          update.StellerPK,
		PGPPK:              update.PGPPK,
		PGPPKHash:          update.PGPPKHash,
		PgpSecret:          update.PgpSecret,
		StellarSecret:      update.StellarSecret,
		DigitalSignature:   update.DigitalSignature,
		Signaturehash:      update.Signaturehash,
		Date:               update.Date,
		Endorsments:        update.Endorsments,
		TXNOrgRegistration: update.TXNOrgRegistration,
	}
	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}
	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}
	c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
	_, err = c.UpdateOne(context.TODO(), bson.M{"pgppkhash": pkhash}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}
func (cd *Connection) UpdateEthereumMetricStatus(metricID string, txnUUID string, update model.EthereumMetricBind) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum metric", err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.EthereumMetricBind{
		MetricID:          update.MetricID,
		MetricName:        update.MetricName,
		Metric:            update.Metric,
		TemplateString:    update.TemplateString,
		BINstring:         update.BINstring,
		ABIstring:         update.ABIstring,
		ContractAddress:   update.ContractAddress,
		ContractName:      update.ContractName,
		Status:            update.Status,
		Timestamp:         update.Timestamp,
		TransactionHash:   update.TransactionHash,
		TransactionCost:   update.TransactionCost,
		TransactionUUID:   update.TransactionUUID,
		TransactionSender: update.TransactionSender,
		User:              update.User,
		ErrorMessage:      update.ErrorMessage,
		FormulaIDs:        update.FormulaIDs,
		ValueIDs:          update.ValueIDs,
		FormulaID:         update.FormulaID,
		Type:              update.Type,
		ActualStatus:      update.ActualStatus,
	}

	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("EthereumMetricBind")
	_, err = c.UpdateOne(context.TODO(), bson.M{"metricid": metricID, "transactionuuid": txnUUID}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}
func (cd *Connection) UpdateTrustNetworkUserPassword(pkhash string, update model.TrustNetWorkUser, newPassword string) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Println("Error while connecting to DB " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())
	up := model.TrustNetWorkUser{
		Name:               update.Name,
		Company:            update.Company,
		Email:              update.Email,
		Password:           newPassword,
		PasswordResetCode:  nil,
		Contact:            update.Contact,
		Industry:           update.Industry,
		StellerPK:          update.StellerPK,
		PGPPK:              update.PGPPK,
		PGPPKHash:          update.PGPPKHash,
		PgpSecret:          update.PgpSecret,
		StellarSecret:      update.StellarSecret,
		DigitalSignature:   update.DigitalSignature,
		Signaturehash:      update.Signaturehash,
		Date:               update.Date,
		Endorsments:        update.Endorsments,
		TXNOrgRegistration: update.TXNOrgRegistration,
	}
	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}
	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}
	c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
	_, err = c.UpdateOne(context.TODO(), bson.M{"pgppkhash": pkhash}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}
func (cd *Connection) UpdateTrustNetworkResetUserPassword(pkhash string, update model.TrustNetWorkUser, ResetPassword []byte) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Println("Error while connecting to DB " + err.Error())
		return err
	}
	defer session.EndSession(context.TODO())
	up := model.TrustNetWorkUser{
		Name:               update.Name,
		Company:            update.Company,
		Email:              update.Email,
		Password:           update.Password,
		PasswordResetCode:  ResetPassword,
		Contact:            update.Contact,
		Industry:           update.Industry,
		StellerPK:          update.StellerPK,
		PGPPK:              update.PGPPK,
		PGPPKHash:          update.PGPPKHash,
		PgpSecret:          update.PgpSecret,
		StellarSecret:      update.StellarSecret,
		DigitalSignature:   update.DigitalSignature,
		Signaturehash:      update.Signaturehash,
		Date:               update.Date,
		Endorsments:        update.Endorsments,
		TXNOrgRegistration: update.TXNOrgRegistration,
	}
	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}
	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}
	c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
	_, err = c.UpdateOne(context.TODO(), bson.M{"pgppkhash": pkhash}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}

func (cd *Connection) UpdateEthereumMetricLatestContract(metricID string, update model.MetricLatestContract) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum latest metric contract", err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.MetricLatestContract{
		MetricID:        update.MetricID,
		ContractAddress: update.ContractAddress,
		Type:            update.Type,
	}

	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("EthMetricLatest")
	_, err = c.UpdateOne(context.TODO(), bson.M{"metricid": metricID}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}

func (cd *Connection) UpdateEthereumPendingContract(transactionHash string, contractAddress, identifier string, update model.PendingContracts) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum pending contract", err.Error())
		return err
	}
	defer session.EndSession(context.TODO())

	up := model.PendingContracts{
		TransactionHash: update.TransactionHash,
		ContractAddress: update.ContractAddress,
		Status:          update.Status,
		CurrentIndex:    update.CurrentIndex,
		ErrorMessage:    update.ErrorMessage,
		ContractType:    update.ContractType,
		Identifier:      update.Identifier,
		Nonce:           update.Nonce,
		GasPrice:        update.GasPrice,
		GasLimit:        update.GasLimit,
	}

	pByte, err := bson.Marshal(up)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	c := session.Client().Database(dbName).Collection("EthereumPendingTransactions")
	_, err = c.UpdateOne(context.TODO(), bson.M{"transactionhash": transactionHash, "contractaddress": contractAddress, "identifier": identifier}, bson.D{{Key: "$set", Value: updateNew}})

	return err
}

func (cd *Connection) UpdateEthFormulaStatusByUUID(txnUUID string, status int, errorMessage string) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum formula status for the given UUID", err.Error())
		return err
	}
	c := session.Client().Database(dbName).Collection("EthereumExpertFormula")
	filter := bson.D{{Key: "transactionuuid", Value: txnUUID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}, {Key: "errormessage", Value: errorMessage}}}}
	_, errWhenUpdate := c.UpdateOne(context.TODO(), filter, update)
	if errWhenUpdate != nil {
		logrus.Error("Error when updating ethereum formula status in DB : " + errWhenUpdate.Error())
		return errWhenUpdate
	}
	return nil
}

func (cd *Connection) UpdateEthMetricStatusByUUID(txnUUID string, status int, errorMessage string) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum metric status for the given UUID", err.Error())
		return err
	}
	c := session.Client().Database(dbName).Collection("EthereumMetricBind")
	filter := bson.D{{Key: "transactionuuid", Value: txnUUID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}, {Key: "errormessage", Value: errorMessage}}}}
	_, errWhenUpdate := c.UpdateOne(context.TODO(), filter, update)
	if errWhenUpdate != nil {
		logrus.Error("Error when updating ethereum metric status in DB : " + errWhenUpdate.Error())
		return errWhenUpdate
	}
	return nil
}

// UpdateSelectedEthMetricFields updates the selected fields of the Ethereum metric bind except the metric request
func (cd *Connection) UpdateSelectedEthMetricFields(metricID string, txnUUID string, metricObj model.EthereumMetricBind) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum metric fields for the given metric ID", err.Error())
		return err
	}
	c := session.Client().Database(dbName).Collection("EthereumMetricBind")
	updateFields := bson.D{}

	if metricObj.MetricID != "" {
		updateFields = append(updateFields, bson.E{Key: "metricid", Value: metricObj.MetricID})
	}
	if metricObj.MetricName != "" {
		updateFields = append(updateFields, bson.E{Key: "metricname", Value: metricObj.MetricName})
	}
	if metricObj.ContractName != "" {
		updateFields = append(updateFields, bson.E{Key: "contractname", Value: metricObj.ContractName})
	}
	if metricObj.ContractAddress != "" {
		updateFields = append(updateFields, bson.E{Key: "contractaddress", Value: metricObj.ContractAddress})
	}
	if metricObj.TemplateString != "" {
		updateFields = append(updateFields, bson.E{Key: "templatestring", Value: metricObj.TemplateString})
	}
	if metricObj.BINstring != "" {
		updateFields = append(updateFields, bson.E{Key: "binstring", Value: metricObj.BINstring})
	}
	if metricObj.ABIstring != "" {
		updateFields = append(updateFields, bson.E{Key: "abistring", Value: metricObj.ABIstring})
	}
	if metricObj.Timestamp != "" {
		updateFields = append(updateFields, bson.E{Key: "timestamp", Value: metricObj.Timestamp})
	}
	if metricObj.TransactionHash != "" {
		updateFields = append(updateFields, bson.E{Key: "transactionhash", Value: metricObj.TransactionHash})
	}
	if metricObj.TransactionCost != "" {
		updateFields = append(updateFields, bson.E{Key: "transactioncost", Value: metricObj.TransactionCost})
	}
	if metricObj.TransactionUUID != "" {
		updateFields = append(updateFields, bson.E{Key: "transactionuuid", Value: metricObj.TransactionUUID})
	}
	if metricObj.TransactionSender != "" {
		updateFields = append(updateFields, bson.E{Key: "transactionsender", Value: metricObj.TransactionSender})
	}
	if metricObj.Status != 0 {
		updateFields = append(updateFields, bson.E{Key: "status", Value: metricObj.Status})
	}
	if metricObj.ErrorMessage != "" {
		updateFields = append(updateFields, bson.E{Key: "errormessage", Value: metricObj.ErrorMessage})
	}
	if metricObj.Type != "" {
		updateFields = append(updateFields, bson.E{Key: "type", Value: metricObj.Type})
	}
	if metricObj.FormulaID != "" {
		updateFields = append(updateFields, bson.E{Key: "formulaid", Value: metricObj.FormulaID})
	}
	if metricObj.FormulaIDs != nil {
		updateFields = append(updateFields, bson.E{Key: "formulaid", Value: metricObj.FormulaIDs})
	}
	if metricObj.ValueIDs != nil {
		updateFields = append(updateFields, bson.E{Key: "valueids", Value: metricObj.ValueIDs})
	}
	if metricObj.User.ID != "" {
		updateFields = append(updateFields, bson.E{Key: "user.id", Value: metricObj.User.ID})
	}
	if metricObj.User.Publickey != "" {
		updateFields = append(updateFields, bson.E{Key: "user.publickey", Value: metricObj.User.Publickey})
	}
	if metricObj.User.TenantID != "" {
		updateFields = append(updateFields, bson.E{Key: "user.tenantID", Value: metricObj.User.TenantID})
	}
	if metricObj.ActualStatus != 0 {
		updateFields = append(updateFields, bson.E{Key: "actualstatus", Value: metricObj.ActualStatus})
	}
	filter := bson.D{{Key: "metricid", Value: metricID}, {Key: "transactionuuid", Value: txnUUID}}
	update := bson.D{{Key: "$set", Value: updateFields}}
	_, errWhenUpdate := c.UpdateOne(context.TODO(), filter, update)
	if errWhenUpdate != nil {
		logrus.Error("Error when updating ethereum metric fields in DB : " + errWhenUpdate.Error())
		return errWhenUpdate
	}
	return nil
}

// UpdateSelectedEthFormulaFields updates the selected fields of the Ethereum expert formula except the formula request and execution template
func (cd *Connection) UpdateSelectedEthFormulaFields(formulaID string, txnUUID string, formulaObj model.EthereumExpertFormula) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("update Ethereum formula fields for the given formula ID", err.Error())
		return err
	}
	c := session.Client().Database(dbName).Collection("EthereumExpertFormula")
	updateFields := bson.D{}
	if formulaObj.FormulaID != "" {
		updateFields = append(updateFields, bson.E{Key: "formulaid", Value: formulaObj.FormulaID})
	}
	if formulaObj.FormulaName != "" {
		updateFields = append(updateFields, bson.E{Key: "formulaname", Value: formulaObj.FormulaName})
	}
	if formulaObj.VariableCount != 0 {
		updateFields = append(updateFields, bson.E{Key: "variablecount", Value: formulaObj.VariableCount})
	}
	if formulaObj.ContractName != "" {
		updateFields = append(updateFields, bson.E{Key: "contractname", Value: formulaObj.ContractName})
	}
	if formulaObj.TemplateString != "" {
		updateFields = append(updateFields, bson.E{Key: "templatestring", Value: formulaObj.TemplateString})
	}
	if formulaObj.BINstring != "" {
		updateFields = append(updateFields, bson.E{Key: "binstring", Value: formulaObj.BINstring})
	}
	if formulaObj.ABIstring != "" {
		updateFields = append(updateFields, bson.E{Key: "abistring", Value: formulaObj.ABIstring})
	}
	if formulaObj.GOstring != "" {
		updateFields = append(updateFields, bson.E{Key: "gostring", Value: formulaObj.GOstring})
	}
	if formulaObj.SetterNames != nil {
		updateFields = append(updateFields, bson.E{Key: "setternames", Value: formulaObj.SetterNames})
	}
	if formulaObj.ContractAddress != "" {
		updateFields = append(updateFields, bson.E{Key: "contractaddress", Value: formulaObj.ContractAddress})
	}
	if formulaObj.Status != 0 {
		updateFields = append(updateFields, bson.E{Key: "status", Value: formulaObj.Status})
	}
	if formulaObj.Timestamp != "" {
		updateFields = append(updateFields, bson.E{Key: "timestamp", Value: formulaObj.Timestamp})
	}
	if formulaObj.TransactionHash != "" {
		updateFields = append(updateFields, bson.E{Key: "transactionhash", Value: formulaObj.TransactionHash})
	}
	if formulaObj.TransactionCost != "" {
		updateFields = append(updateFields, bson.E{Key: "transactioncost", Value: formulaObj.TransactionCost})
	}
	if formulaObj.TransactionUUID != "" {
		updateFields = append(updateFields, bson.E{Key: "transactionuuid", Value: formulaObj.TransactionUUID})
	}
	if formulaObj.TransactionSender != "" {
		updateFields = append(updateFields, bson.E{Key: "transactionsender", Value: formulaObj.TransactionSender})
	}
	if formulaObj.Verify.Payload != "" {
		updateFields = append(updateFields, bson.E{Key: "verify.payload", Value: formulaObj.Verify.Payload})
	}
	if formulaObj.Verify.Signature != "" {
		updateFields = append(updateFields, bson.E{Key: "verify.signature", Value: formulaObj.Verify.Signature})
	}
	if formulaObj.Verify.PublicKey != "" {
		updateFields = append(updateFields, bson.E{Key: "verify.publickey", Value: formulaObj.Verify.PublicKey})
	}
	if formulaObj.ErrorMessage != "" {
		updateFields = append(updateFields, bson.E{Key: "errormessage", Value: formulaObj.ErrorMessage})
	}
	if formulaObj.ActualStatus != 0 {
		updateFields = append(updateFields, bson.E{Key: "actualstatus", Value: formulaObj.ActualStatus})
	}

	filter := bson.D{{Key: "formulaid", Value: formulaID}, {Key: "transactionuuid", Value: txnUUID}}
	update := bson.D{{Key: "$set", Value: updateFields}}
	_, errWhenUpdate := c.UpdateOne(context.TODO(), filter, update)
	if errWhenUpdate != nil {
		logrus.Error("Error when updating ethereum formula fields in DB : " + errWhenUpdate.Error())
		return errWhenUpdate
	}
	return nil
}

func (cd *Connection) UpdateCOCState(data model.UpdateCOCState) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while connecting to DB " + err.Error())
		return err
	}
	c := session.Client().Database(dbName).Collection("cocTransfers")
	filter := bson.D{{Key: "_id", Value: data.Id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "cocstatus", Value: data.COCStatus}}}}
	_, errWhenUpdate := c.UpdateOne(context.TODO(), filter, update)
	if errWhenUpdate != nil {
		logrus.Error("Error Updating COC status: " + errWhenUpdate.Error())
		return errWhenUpdate
	}
	return nil
}
