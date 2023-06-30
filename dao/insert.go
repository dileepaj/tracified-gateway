package dao

import (
	"context"
	"log"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/model"
	notificationhandler "github.com/dileepaj/tracified-gateway/services/notificationHandler.go"
	"github.com/dileepaj/tracified-gateway/utilities"
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
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertCoc ------------------------", constants.INFO)
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("COC")
	_, err = c.InsertOne(context.TODO(), Coc)
	if err != nil {
		logger.LogWriter("Error while inserting to COC "+err.Error(), constants.ERROR)
	}
	return err
}

/*
InsertTransaction Insert a single Transaction Object to TransactionCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertTransaction(Coc model.TransactionCollectionBody) error {
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertTransaction ------------------------", constants.INFO)

	// result := model.TransactionCollectionBody{}
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("Transactions")
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"txnhash", Coc.TxnHash}}
	pByte, err := bson.Marshal(Coc)
	if err != nil {
		logger.LogWriter("Error when marshalling "+err.Error(), constants.ERROR)
		return err
	}
	var updateNew bson.M
	err = bson.Unmarshal(pByte, &updateNew)
	if err != nil {
		logger.LogWriter("Error when unmarshalling "+err.Error(), constants.ERROR)
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
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertSpecialToTempOrphan ------------------------", constants.INFO)
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("TempOrphan")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		logger.LogWriter("Error while inserting to TempOrphan "+err.Error(), constants.ERROR)
	}
	return err
}

/*
InsertToOrphan Insert a single Transaction Object to OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertToOrphan(Coc model.TransactionCollectionBody) error {
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertToOrphan ------------------------", constants.INFO)

	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("TempOrphan")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		logger.LogWriter("Error while inserting to Orphan "+err.Error(), constants.ERROR)
	}
	return err
}

/*
InsertProfile Insert a single Profile Object to ProfileCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertProfile(Coc model.ProfileCollectionBody) error {
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertProfile ------------------------", constants.INFO)

	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}

	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Profiles")
	_, err = c.InsertOne(context.TODO(), Coc)

	if err != nil {
		logger.LogWriter("Error while inserting to Profiles "+err.Error(), constants.ERROR)
	}
	return err
}

/*
InsertCertificate Insert a single Certificate Object to CertificateCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertCertificate(Cert model.CertificateCollectionBody) error {
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertCertificate ------------------------", constants.INFO)
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Certificates")
	_, err = c.InsertOne(context.TODO(), Cert)

	if err != nil {
		logger.LogWriter("Error while inserting to Certificates "+err.Error(), constants.ERROR)
	}
	return err
}

func (cd *Connection) InsertArtifact(artifacts model.ArtifactTransaction) error {
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertArtifact ------------------------", constants.INFO)

	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)

	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Artifacts")
	_, err = c.InsertOne(context.TODO(), artifacts)

	if err != nil {
		logger.LogWriter("Error while inserting to Artifacts "+err.Error(), constants.ERROR)
	}
	return err
}

func (cd *Connection) InsertOrganization(Org model.TestimonialOrganization) error {
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertOrganization ------------------------", constants.INFO)

	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Organizations")
	_, err = c.InsertOne(context.TODO(), Org)
	if err != nil {
		logger.LogWriter("Error while inserting to organizations "+err.Error(), constants.ERROR)
	}
	return err
}

func (cd *Connection) InsertTestimonial(Tes model.Testimonial) error {
	logger := utilities.NewCustomLogger()
	logger.LogWriter("--------------------------- InsertOrganization ------------------------", constants.INFO)

	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)

	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("Testimonials")
	_, err = c.InsertOne(context.TODO(), Tes)

	if err != nil {
		logger.LogWriter("Error while inserting to organizations "+err.Error(), constants.ERROR)
	}
	return err
}

// insert new proof presentation protocol
func (cd *Connection) InsertProofProtocol(protocol model.ProofProtocol) error {
	logger := utilities.NewCustomLogger()

	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error when connecting to DB "+err.Error(), constants.ERROR)

	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("ProofProtocols")
	_, err = c.InsertOne(context.TODO(), protocol)
	if err != nil {
		logger.LogWriter("Error when inserting data to DB "+err.Error(), constants.ERROR)

	}
	return err
}

func (cd *Connection) InsertIdentifier(id apiModel.IdentifierModel) error {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("IdentifierMap")
	_, err = c.InsertOne(context.TODO(), id)

	if err != nil {
		logger.LogWriter("Error while inserting to TempOrphan "+err.Error(), constants.ERROR)
	}
	return err
}

func (cd *Connection) InsertSolanaNFT(solanaNFT model.NFTWithTransactionSolana, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error when connecting to DB "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTSolana")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), solanaNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		logger.LogWriter("Error while inserting to NFTSolana DB "+err.Error(), constants.ERROR)
	}
	if err2 != nil {
		logger.LogWriter("Error while inserting to Marketplace DB "+err2.Error(), constants.ERROR)
	}
	return err, err2
}

func (cd *Connection) InsertPolygonNFT(polyNFT model.NFTWithTransactionContracts, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error when connecting to DB "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTPolygon")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), polyNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		logger.LogWriter("Error while inserting to NFTPolygon DB "+err.Error(), constants.ERROR)
	}
	if err2 != nil {
		logger.LogWriter("Error while inserting to Marketplace DB "+err2.Error(), constants.ERROR)
	}
	return err, err2
}

func (cd *Connection) InsertEthereumNFT(etherNFT model.NFTWithTransactionContracts, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error when connecting to DB "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTEthereum")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), etherNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		logger.LogWriter("Error while inserting to NFTEthereum DB "+err.Error(), constants.ERROR)
	}
	if err2 != nil {
		logger.LogWriter("Error while inserting to Marketplace DB "+err2.Error(), constants.ERROR)
	}
	return err, err2
}

func (cd *Connection) InsertStellarNFT(stellarNFT model.NFTWithTransaction, marketPlaceNFT model.MarketPlaceNFT) (error, error) {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error when connecting to DB "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTStellar")
	c2 := session.Client().Database(dbName).Collection("MarketPlaceNFT")
	_, err2 := c.InsertOne(context.TODO(), stellarNFT)
	_, err = c2.InsertOne(context.TODO(), marketPlaceNFT)
	if err != nil {
		logger.LogWriter("Error while inserting to NFTStellar DB "+err.Error(), constants.ERROR)
	}
	if err2 != nil {
		logger.LogWriter("Error while inserting to Marketplace DB "+err2.Error(), constants.ERROR)
	}
	return err, err2
}

func (cd *Connection) InsertStellarNFTKeys(nftKeys model.NFTKeys) error {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error when connecting to DB "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("NFTKeys")
	_, err = c.InsertOne(context.TODO(), nftKeys)
	if err != nil {
		logger.LogWriter("Error while inserting to NFTStellar DB "+err.Error(), constants.ERROR)
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

func (cd *Connection) InsertMetricBindingFormula(metricBind model.MetricBindingStore) (string, error) {
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

func (cd *Connection) InsertRSAKeyPair(rsaKey model.RSAKeyPair) error {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("RSAKeys")
	_, err = c.InsertOne(context.TODO(), rsaKey)
	if err != nil {
		logger.LogWriter("Error when inserting data to NFTStellar DB "+err.Error(), constants.ERROR)
	}
	return err
}

func (cd *Connection) InsertToWorkflowIDMap(tenentIDMap model.WorkflowMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("WorkflowIDMap")
	_, err = c.InsertOne(context.TODO(), tenentIDMap)
	if err != nil {
		logrus.Info("Error when inserting workflow id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToArtifactTemplateIDMap(artifactMap model.ArtifactTemplateId) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("ArtifactTemplateIDMap")
	_, err = c.InsertOne(context.TODO(), artifactMap)
	if err != nil {
		logrus.Info("Error when inserting artifact id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertBindKey(bindKey model.BindKeyMap) (string, error) {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("BindKeyMap")
	result, err := c.InsertOne(context.TODO(), bindKey)
	if err != nil {
		logrus.Info("Error when inserting MetricBinding to DB " + err.Error())
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (cd *Connection) InsertToEthFormulaDetails(ethFormulaMap model.EthereumExpertFormula) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("insert Ethereum expert formula", err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthereumExpertFormula")
	_, err = c.InsertOne(context.TODO(), ethFormulaMap)
	if err != nil {
		logrus.Info("Error when inserting formula details to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToPrimaryKeyIdMap(artifactMap model.PrimaryKeyMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("PrimaryKeyIdMap")
	_, err = c.InsertOne(context.TODO(), artifactMap)
	if err != nil {
		logrus.Info("Error when inserting artifact id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertEthFormulaIDMap(formulaIDMap model.EthFormulaIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("insert Ethereum formula mapped ID", err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthereumFormulaIDMap")
	_, err = c.InsertOne(context.TODO(), formulaIDMap)
	if err != nil {
		logrus.Info("Error when inserting formula id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) SaveTrustNetworkUser(user model.TrustNetWorkUser) (string, error) {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error when connecting to DB "+err.Error(), constants.ERROR)
		return "", err
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("TracifiedTrustNetwork")
	result, err1 := c.InsertOne(context.TODO(), user)
	if err1 != nil {
		logger.LogWriter("Error when inserting data to NFTStellar DB "+err1.Error(), constants.ERROR)
		return "", err1
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (cd *Connection) InsertPGPAccount(pgpAccount model.PGPAccount) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("PGPAccounts")
	_, err = c.InsertOne(context.TODO(), pgpAccount)
	if err != nil {
		logrus.Info("Error when inserting pgpAccount id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToEthMetricDetails(ethMetricDetails model.EthereumMetricBind) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("insert Ethereum metric bind", err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthereumMetricBind")
	_, err = c.InsertOne(context.TODO(), ethMetricDetails)
	if err != nil {
		logrus.Info("Error when inserting metric details to DB " + err.Error())
	}
	return err
}

func (cd *Connection) EthereumInsertToValueIDMap(valueIDMap model.ValueIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("insert Ethereum mapped value ID", err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthValueIDMap")
	_, err = c.InsertOne(context.TODO(), valueIDMap)
	if err != nil {
		logrus.Info("Error when inserting Counters to DB " + err.Error())
	}
	return err
}

func (cd *Connection) EthereumInsertToMetricLatestContract(contractObj model.MetricLatestContract) error {
	session, err := cd.connect()
	if err != nil {
		notificationhandler.InformDBConnectionIssue("insert latest metric contract", err.Error())
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthMetricLatest")
	_, err = c.InsertOne(context.TODO(), contractObj)
	if err != nil {
		logrus.Info("Error when inserting latest contract to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertEthMetricIDMap(metricIDMap model.EthMetricIDMap) error {
	session, err := cd.connect()
	if err != nil {
		logrus.Info("Error when connecting to DB " + err.Error())
		notificationhandler.InformDBConnectionIssue("insert Ethereum metric mapped ID", err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthereumMetricIDMap")
	_, err = c.InsertOne(context.TODO(), metricIDMap)
	if err != nil {
		logrus.Info("Error when inserting metric id to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertEthErrorMessage(errorMessage model.EthErrorMessage) error {
	session, err := cd.connect()
	if err != nil {
		notificationhandler.InformDBConnectionIssue("insert ethereum error message", err.Error())
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthErrorMessages")
	_, err = c.InsertOne(context.TODO(), errorMessage)
	if err != nil {
		logrus.Info("Error when inserting ethereum error messages to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertEthPendingContract(pendingTransaction model.PendingContracts) error {
	session, err := cd.connect()
	if err != nil {
		notificationhandler.InformDBConnectionIssue("insert ethereum pending contract", err.Error())
		logrus.Info("Error when connecting to DB " + err.Error())
	}
	defer session.EndSession(context.TODO())
	c := session.Client().Database(dbName).Collection("EthereumPendingTransactions")
	_, err = c.InsertOne(context.TODO(), pendingTransaction)
	if err != nil {
		logrus.Info("Error when inserting pending transaction to DB " + err.Error())
	}
	return err
}

func (cd *Connection) InsertToNFTStatus(NFT model.PendingNFTS) error {
	logger := utilities.NewCustomLogger()
	session, err := cd.connect()
	if err != nil {
		logger.LogWriter("Error while getting session "+err.Error(), constants.ERROR)
	}
	defer session.EndSession(context.TODO())

	c := session.Client().Database(dbName).Collection("NFTStatus")
	_, err = c.InsertOne(context.TODO(), NFT)

	if err != nil {
		logger.LogWriter("Error while inserting to NFTStatus "+err.Error(), constants.ERROR)
	}
	return err
}
