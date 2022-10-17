package dao

import (
	"context"
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
InsertCoc Insert a single COC Object to COCCollection in DB
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

/*
InsertTransaction Insert a single Transaction Object to TransactionCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertTransaction(Coc model.TransactionCollectionBody) error {
	log.Println("--------------------------- InsertTransaction ------------------------")
	// result := model.TransactionCollectionBody{}
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
	log.Printf("found document %v", result)
	return err
}

/*
InsertSpecialToTempOrphan Insert a single Transaction Object to TempOrphan in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertSpecialToTempOrphan(Coc model.TransactionCollectionBody) error {
	fmt.Println("--------------------------- InsertSpecialToTempOrphan ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("TempOrphan")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		fmt.Println("Error while inserting to TempOrphan " + err.Error())
	}
	return err
}

/*
InsertToOrphan Insert a single Transaction Object to OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertToOrphan(Coc model.TransactionCollectionBody) error {
	fmt.Println("--------------------------- InsertToOrphan ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("TempOrphan")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		fmt.Println("Error while inserting to Orphan " + err.Error())
	}
	return err
}

/*
InsertProfile Insert a single Profile Object to ProfileCollection in DB
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

/*
InsertCertificate Insert a single Certificate Object to CertificateCollection in DB
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

// insert new proof presentation protocol
func (cd *Connection) InsertProofProtocol(protocol model.ProofProtocol) error {
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("ProofProtocols")
	_, err = c.InsertOne(context.TODO(), protocol)
	if err != nil {
		fmt.Println("Error when inserting data to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertIdentifier(id apiModel.IdentifierModel) error {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error while getting session " + err.Error())
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("IdentifierMap")
	_, err = c.InsertOne(context.TODO(), id)

	if err != nil {
		log.Println("Error while inserting to TempOrphan " + err.Error())
	}
	return err
}

func (cd *Connection) InsertSolanaNFT(solanaNFT model.NFTWithTransactionSolana, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTSolana")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), solanaNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		log.Println("Error when inserting data to NFTSolana DB " + err.Error())
	}
	if err2 != nil {
		log.Println("Error when inserting data to MarketPlaceNFT DB " + err.Error())
	}
	return err, err2
}

func (cd *Connection) InsertPolygonNFT(polyNFT model.NFTWithTransactionContracts, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTPolygon")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), polyNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		log.Println("Error when inserting data to NFTPolygon DB " + err.Error())
	}
	if err2 != nil {
		log.Println("Error when inserting data to MarketPlaceNFT DB " + err.Error())
	}
	return err, err2
}

func (cd *Connection) InsertEthereumNFT(etherNFT model.NFTWithTransactionContracts, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTEthereum")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), etherNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		log.Println("Error when inserting data to NFTEthereum DB " + err.Error())
	}
	if err2 != nil {
		log.Println("Error when inserting data to MarketPlaceNFT DB " + err.Error())
	}
	return err, err2
}

func (cd *Connection) InsertStellarNFT(stellarNFT model.NFTWithTransaction, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTStellar")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), stellarNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		log.Println("Error when inserting data to NFTStellar DB " + err.Error())
	}
	if err2 != nil {
		log.Println("Error when inserting data to MarketPlaceNFT DB " + err.Error())
	}
	return err, err2
}

func (cd *Connection) InsertStellarNFTKeys(nftKeys model.NFTKeys) error {
	session, err := cd.connect()
	if err != nil {
		log.Println("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTKeys")
	_, err = c.InsertOne(context.TODO(), nftKeys)
	if err != nil {
		log.Println("Error when inserting data to NFTStellar DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertFormulaIDMap(formulaIDMap model.FormulaIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("FormulaIDMap")
	_, err = c.InsertOne(context.TODO(), formulaIDMap)
	if err != nil {
		logrus.Info("Error when inserting Counters to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertExpertIDMap(expertIDMap model.ExpertIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("ExpertIDMap")
	_, err = c.InsertOne(context.TODO(), expertIDMap)
	if err != nil {
		logrus.Info("Error when inserting Counters to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToValueIDMap(valueIDMap model.ValueIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("ValueIDMap")
	_, err = c.InsertOne(context.TODO(), valueIDMap)
	if err != nil {
		logrus.Info("Error when inserting Counters to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToUnitIDMap(unitMap model.UnitIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("UnitIDMap")
	_, err = c.InsertOne(context.TODO(), unitMap)
	if err != nil {
		logrus.Info("Error when inserting Counters to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToAPIThrottler(throttellerReq model.ThrottlerRecord) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("APIThrottleCounter")
	_, err = c.InsertOne(context.TODO(), throttellerReq)
	if err != nil {
		logrus.Info("Error when inserting new API request to DB " + err.Error())
	}
	return err
}

// Insert ExpertFormula Details to DB
func (cd *Connection) InsertExpertFormula(expertFormula model.FormulaStore) (string, error) {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("ExpertFormula")
	result, err := c.InsertOne(context.TODO(), expertFormula)
	if err != nil {
		logrus.Info("Error when inserting Expert Formula to DB " + err.Error())
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (cd *Connection) InsertToResourceIDMap(resourceIDMap model.ResourceIdMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("ResourceIDMap")
	_, err = c.InsertOne(context.TODO(), resourceIDMap)
	if err != nil {
		logrus.Info("Error when inserting resource id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertMetricMapID(metricIDMap model.MetricMapDetails) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("MetricIDMap")
	_, err = c.InsertOne(context.TODO(), metricIDMap)
	if err != nil {
		logrus.Info("Error when inserting metric id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertTenentMapID(tenentIDMap model.TenentMapDetails) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("TenentIDMap")
	_, err = c.InsertOne(context.TODO(), tenentIDMap)
	if err != nil {
		logrus.Info("Error when inserting tenent id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertActivityID(activityDetails model.ActivityMapDetails) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("ActivityIDMap")
	_, err = c.InsertOne(context.TODO(), activityDetails)
	if err != nil {
		logrus.Info("Error when inserting activity id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertMetricBindingFormula(metricBind model.MetricDataBindingRequest) (string, error) {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("MetricBinds")
	result, err := c.InsertOne(context.TODO(), metricBind)
	if err != nil {
		logrus.Info("Error when inserting MetricBinding to DB " + err.Error())
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (cd *Connection) InsertToWorkflowIDMAP(tenentIDMap model.WorkflowMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("WorkflowIDMAP")
	_, err = c.InsertOne(context.TODO(), tenentIDMap)
	if err != nil {
		logrus.Info("Error when inserting workflow id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToArtifactIDMAP(artifactMap model.ArtifactIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("ArtifactIDMAP")
	_, err = c.InsertOne(context.TODO(), artifactMap)
	if err != nil {
		logrus.Info("Error when inserting artifact id to DB " + err.Error())
	}
	return err
}