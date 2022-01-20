package dao

import (
	"context"
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/model"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"
)

/*InsertCoc Insert a single COC Object to COCCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertCoc(Coc model.COCCollectionBody) error {
	fmt.Println("--------------------------- InsertCoc ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("COC")
	_, err = c.InsertOne(context.TODO(), Coc)
	if err != nil {
		fmt.Println("Error while inserting to COC " + err.Error())
	}
	return err
}

/*InsertTransaction Insert a single Transaction Object to TransactionCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertTransaction(Coc model.TransactionCollectionBody) error {
	log.Println("--------------------------- InsertTransaction ------------------------")
	//result := model.TransactionCollectionBody{}
	session, err := cd.connect()
	if err != nil {
		log.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("Transactions")

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"txnhash", Coc.TxnHash}}

	pByte, err := bson.Marshal(Coc)
	if err != nil {
		return err
	}

	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		return err
	}

	update := bson.D{{"$set", updateNew}}
	result, err := c.UpdateOne(context.TODO(), filter, update, opts)

	fmt.Printf("found document %v", result)
	return err
}

/*InsertToOrphan Insert a single Transaction Object to OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertSpecialToTempOrphan(Coc model.TransactionCollectionBody) error {
	fmt.Println("--------------------------- InsertToOrphan ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Orphan")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		fmt.Println("Error while inserting to Orphan " + err.Error())
	}
	return err
}

/*InsertToOrphan Insert a single Transaction Object to OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertToOrphan(Coc model.TransactionCollectionBody) error {
	fmt.Println("--------------------------- InsertToOrphan ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Orphan")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		fmt.Println("Error while inserting to Orphan " + err.Error())
	}
	return err
}

/*InsertProfile Insert a single Profile Object to ProfileCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertProfile(Coc model.ProfileCollectionBody) error {
	fmt.Println("--------------------------- InsertProfile ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Profiles")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		fmt.Println("Error while inserting to Profiles " + err.Error())
	}
	return err
}

/*InsertCertificate Insert a single Certificate Object to CertificateCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertCertificate(Cert model.CertificateCollectionBody) error {
	fmt.Println("--------------------------- InsertCertificate ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Certificates")
	_, err = c.InsertOne(context.TODO(), Cert)

	if err != nil {
		fmt.Println("Error while inserting to Certificates " + err.Error())
	}
	return err
}

func (cd *Connection) InsertArtifact(artifacts model.ArtifactTransaction) error {
	fmt.Println("--------------------------- InsertArtifact ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Artifacts")
	_, err = c.InsertOne(context.TODO(), artifacts)

	if err != nil {
		fmt.Println("Error while inserting to Artifacts " + err.Error())
	}
	return err
}

func (cd *Connection) InsertOrganization(Org model.TestimonialOrganization) error {
	fmt.Println("--------------------------- InsertOrganization ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Organizations")
	_, err = c.InsertOne(context.TODO(), Org)
	if err != nil {
		fmt.Println("Error while inserting to organizations " + err.Error())
	}
	return err
}

func (cd *Connection) InsertTestimonial(Tes model.Testimonial) error {
	fmt.Println("--------------------------- InsertOrganization ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Testimonials")
	_, err = c.InsertOne(context.TODO(), Tes)

	if err != nil {
		fmt.Println("Error while inserting to organizations " + err.Error())
	}
	return err
}


//insert new proof presentation protocol 
func (cd *Connection) InsertProofProtocol(protocol model.ProofProtocol) error{
	session, err := cd.connect()
	if err != nil{
		fmt.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("ProofProtocols")
	_, err = c.InsertOne(context.TODO(), protocol)
	if err != nil{
		fmt.Println("Error when inserting data to DB " + err.Error())
	}
	return err
}

//insert NFT to NFTStellar collection
func (cd *Connection) InsertStellarNFT(stellarNFT model.NFTWithTransaction,marketPlaceNFT model.MarketPlaceNFT) (error,error){
	session, err := cd.connect()
	if err != nil{
		fmt.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("NFTStellar")
	c2:= session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_,err2:=c.InsertOne(context.TODO(), stellarNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil{
		fmt.Println("Error when inserting data to NFTStellar DB " + err.Error())
	}
	if err2 != nil{
		fmt.Println("Error when inserting data to MarketPlaceNFT DB " + err.Error())
	}
	return err, err2
}