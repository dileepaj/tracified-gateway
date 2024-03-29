package routes

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/businessFacades"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes An Array of type Route
type Routes []Route

/*
routes contains all the routes
@author Azeem Ashraf, Jajeththanan Sabapathipillai
*/
var routes = Routes{
	Route{
		"Get server health",
		"GET",
		"/health",
		businessFacades.HealthCheck,
	},
	Route{
		"POC",
		"GET",
		"/proof/poc/{Txn}",
		businessFacades.CheckPOCV3, // Calls the Deprecated POC for Gateway Version 1, Should call the new CheckPOCV3
	},
	Route{
		"FULLPOC",
		"POST",
		"/proof/fullpoc/{Txn}",
		businessFacades.CheckFullPOC, // Calls the Deprecated FULLPOC for Gateway Version 1
	},
	Route{
		"POE",
		"GET",
		"/proof/poe",
		businessFacades.CheckPOEV3, // Calls the Functional POE for Gateway Version 3
	},
	Route{
		"POG",
		"GET",
		"/proof/pog/{Txn}",
		businessFacades.CheckPOGV3Rewrite, // Calls the Functional POG for Gateway Version 3
	},
	Route{
		"POCOC",
		"GET",
		"/proof/pococ/{TxnId}",
		businessFacades.CheckPOCOCV3, // Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"GetCOCCollectionBySender",
		"get",
		"/getcocbysender/{Sender}",
		businessFacades.GetCocBySender, // Calls the Functional GetCocBySender for Gateway Version 3
	},
	Route{
		"GetCOCCollectionByReceiver",
		"get",
		"/getcocbyreceiver/{Receiver}",
		businessFacades.GetCocByReceiver, // Calls the Functional GetCocByReceiver for Gateway Version 3
	},
	Route{
		"InsertCOCCollection",
		"POST",
		"/insertcoccollection",
		businessFacades.InsertCocCollection, // Calls the Functional InsertCocCollection for Gateway Version 3
	},
	Route{
		"InsertCOCCollection",
		"PUT",
		"/insertcoccollection",
		businessFacades.UpdateCocCollection, // Calls the Functional UpdateCocCollection for Gateway Version 3
	},
	Route{
		"SubmitXDR",
		"POST",
		"/transaction/dataPacket",
		businessFacades.SubmitData, // Calls the Functional SubmitData for Gateway Version 3
	},
	Route{
		"SubmitSplit",
		"POST",
		"/transaction/split",
		businessFacades.SubmitData, // Calls the Functional SubmitSplit for Gateway Version 3
	},
	Route{
		"SubmitGenesis",
		"POST",
		"/transaction/genesis",
		businessFacades.SubmitData, // Calls the Functional SubmitGenesis for Gateway Version 3
	},
	Route{
		"SubmitTransformation",
		"POST",
		"/transaction/transformation",
		businessFacades.SubmitData, // Calls the Functional SubmitTransformation for Gateway Version 3
	},
	Route{
		"SubmitMerge",
		"POST",
		"/transaction/merge",
		businessFacades.SubmitData, // Calls the Functional SubmitMerge for Gateway Version 3
	},
	Route{
		"SubmitTransfer",
		"POST",
		"/transaction/transfer",
		businessFacades.SubmitData, // Calls the Functional SubmitTransfer for Gateway Version 3
	},
	Route{
		"InsertCertificate",
		"POST",
		"/transaction/certificateInsert",
		businessFacades.SubmitCertificateInsert, // Calls the Functional SubmitCertificateInsert for Gateway Version 3
	},
	Route{
		"RenewCertificate",
		"POST",
		"/transaction/certificateRenew",
		businessFacades.SubmitCertificateRenewal, // Calls the Functional SubmitCertificateRenewal for Gateway Version 3
	},
	Route{
		"RevokeCertificate",
		"POST",
		"/transaction/certificateRevoke",
		businessFacades.SubmitCertificateRevoke, // Calls the Functional SubmitCertificateRevoke for Gateway Version 3
	},
	Route{
		"LastTxn",
		"GET",
		"/transaction/lastTxn/{Identifier}",
		businessFacades.LastTxn, // Calls the Functional LastTxn for Gateway Version 3
	},
	Route{
		"SubAccountStatus",
		"POST",
		"/transaction/coc/subAccountStatus",
		businessFacades.CheckAccountsStatus, // Calls the Functional CheckAccountsStatus for Gateway Version 3
	},
	Route{
		"POCDeveloperRetriever",
		"get",
		"/pocbctree/{Txn}",
		businessFacades.DeveloperRetriever, // Test
	},
	Route{
		"GET LOGS",
		"get",
		"/getLogsForToday/",
		businessFacades.RetrieveLogsForToday, // Test
	},
	Route{
		"POCGatewayRetrieverForTDP",
		"GET",
		"/gatewayTree/{Txn}",
		businessFacades.GatewayRetriever, // Test
	},
	Route{
		"POCGatewayRetrieverForIdentifier",
		"GET",
		"/gatewayTreeWithIdentifier/{Identifier}",
		businessFacades.GatewayRetrieverWithIdentifier, // Test
	},
	Route{
		"ConvertXDRToTXN",
		"POST",
		"/xdrToTxn",
		businessFacades.ConvertXDRToTXN, // Test
	},
	Route{
		"LastCOC",
		"GET",
		"/lastCoc/{Identifier}",
		businessFacades.LastCOC, // Test
	},
	Route{
		"Retrieve TDP for Transaction",
		"GET",
		"/tdpForTxn/{Txn}",
		businessFacades.TDPForTXN, // Test
	},
	Route{
		"Retrieve Transaction for TDP",
		"GET",
		"/txnForTdp/{Txn}",
		businessFacades.TXNForTDP, // Test
	},
	Route{
		"Transactions",
		"POST",
		"/transaction/type/{TType}",
		businessFacades.Transaction, // Deprecated
	},
	Route{
		"TrustLine",
		"POST",
		"/create/Trustline",
		businessFacades.CreateTrust, // Deprecated
	},
	Route{
		"SendAssestV2",
		"POST",
		"/send/asset",
		businessFacades.SendAssests, // Deprecated
	},
	Route{
		"lockAcc",
		"POST",
		"/lock/registrarAcc",
		businessFacades.MultisigAccount, // Deprecated
	},
	Route{
		"UnlockAcc",
		"POST",
		"/Appoint/Registrar",
		businessFacades.AppointRegistrar, // Deprecated
	},
	Route{
		"transformV2",
		"POST",
		"/transform/V2",
		businessFacades.TransformV2, // Deprecated
	},
	Route{
		"COC",
		"POST",
		"/COC/Transaction",
		businessFacades.COC, // Deprecated
	},
	Route{
		"COCLink",
		"POST",
		"/COCLink/Transaction",
		businessFacades.COCLink, // Deprecated
	},
	Route{
		"TransactionId",
		"GET",
		"/TransactionId/{id}",
		businessFacades.GetTransactionId, // Test
	},
	Route{
		"TransactionIds",
		"GET",
		"/GetTransactionsForTDP/{id}",
		businessFacades.GetTransactionsForTDP, // Test
	},
	Route{
		"TransactionIdsForTDPs",
		"POST",
		"/GetTransactionsForTDPs",
		businessFacades.GetTransactionsForTdps, // Test
	},
	Route{
		"TransactionIdsForPK",
		"GET",
		"/GetTransactionsForPK/{id}",
		businessFacades.GetTransactionsForPK, // Test
	},
	Route{
		"RetriveTransactionId",
		"GET",
		"/GetTransactionId/{id}",
		businessFacades.RetriveTransactionId, // Test
	},
	Route{
		"QueryTransactionsByKey",
		"GET",
		"/GetTransactions",
		businessFacades.QueryTransactionsByKey, // multisearch
	},
	Route{
		"GetCOCByTxn",
		"GET",
		"/GetCOCByTxn/{txn}",
		businessFacades.GetCOCByTxn, // multisearch
	},
	Route{
		"RetriveTransactionId",
		"GET",
		"/GetTransactionId/{id}",
		businessFacades.RetriveTransactionId, // Test
	},
	Route{
		"RetrievePreviousTranasctions",
		"GET",
		"/RetrievePreviousTranasctions",
		businessFacades.RetrievePreviousTranasctions, // Test
	},
	Route{
		"GetTotalRecordsCountInTransactionCollection",
		"GET",
		"/RetrievePreviousTranasctionsCount",
		businessFacades.RetrievePreviousTranasctionsCount,
	},
	Route{
		"ArtifactTransactions",
		"POST",
		"/Insert/ArtifactTransactions",
		businessFacades.ArtifactTransactions, // Test
	},
	Route{
		"InsertOrganization",
		"POST",
		"/organization",
		businessFacades.InsertOrganization, // Test
	},
	Route{
		"GetAllOrganizations",
		"GET",
		"/approved/organization",
		businessFacades.GetAllOrganizations, // Test
	},
	Route{
		"GetOrganizationByPublicKey",
		"GET",
		"/organization/{PK}",
		businessFacades.GetOrganizationByPublicKey, // Test
	},
	Route{
		"UpdateOrganization",
		"PUT",
		"/organization",
		businessFacades.UpdateOrganization, // Test
	},
	Route{
		"InsertTestimonial",
		"POST",
		"/testimonial",
		businessFacades.InsertTestimonial, // Test
	},
	Route{
		"GetTestimonialBySender",
		"GET",
		"/testimonial/sender/{PK}",
		businessFacades.GetTestimonialBySender, // Test
	},
	Route{
		"GetTestimonialByReciever",
		"GET",
		"/testimonial/reciever/{PK}",
		businessFacades.GetTestimonialByReciever, // Test
	},
	Route{
		"UpdateTestimonial",
		"PUT",
		"/testimonial",
		businessFacades.UpdateTestimonial, // Test
	},
	Route{
		"SubAccountStatusExtended",
		"POST",
		"/transaction/subAccountStatus",
		businessFacades.CheckAccountsStatusExtended, // Test
	},
	Route{
		"GetAllPendingAndRejectedOrganizations",
		"GET",
		"/notapproved/organization",
		businessFacades.GetAllPendingAndRejectedOrganizations,
	},
	Route{
		"BlockchainRetrieverWithHash",
		"GET",
		"/blockchain/{txn}",
		businessFacades.BlockchainDataRetreiverWithHash,
	},
	Route{
		"BlockchainTreeRetrieverWithHash",
		"GET",
		"/pocv4/{txn}",
		businessFacades.BlockchainTreeRetreiverWithHash,
	},
	Route{
		"GetProofPresentationProtocolByProofName",
		"GET",
		"/getproofprotocol/{proofname}",
		businessFacades.GetProofPresentationProtocolByProofName, // Test - calls the GetProofPresentationProtocolByProofName in the ProofPresesntationHandlers
	},
	Route{
		"InsertProofPresentationProtocol",
		"POST",
		"/proofprotocol",
		businessFacades.InsertProofPresentationProtocol, // Test - calls the InsertProofPresentationProtocol in the ProofPresesntationHandlers
	},
	Route{
		"UpdateProofPresentationProtocol",
		"PUT",
		"/proofprotocol",
		businessFacades.UpdateProofPresentationProtocol, // Test - calls the UpdateProofPresesntationProtocol in the ProofPresesntationHandlers
	},
	Route{
		"DeleteProofPresentationProtocolByProofName",
		"DELETE",
		"/deleteproofprotocol/{proofname}",
		businessFacades.DeleteProofPresentationProtocolByProofName, // Test - calls the DeleteProofPresentationProtocolByProtocolName in the ProofPresesntationHandlers
	},
	Route{
		"GetAllOrganizationsPaginated",
		"GET",
		"/approved/organizationPaginated",
		businessFacades.GetAllOrganizations_Paginated, // Test
	},
	Route{
		"EnableCorsAndResponse",
		"GET",
		"/enable-cors",
		businessFacades.EnableCorsAndResponse,
	},
	Route{
		"Get all trasacion by identifer",
		"GET",
		"/transaction/identifier/{identifier}",
		businessFacades.TxnForIdentifier,
	},
	Route{
		"Get all trasacion by Artifact",
		"GET",
		"/transaction/identifier/artifact/{artifactid}",
		businessFacades.TxnForArtifact,
	},
	Route{
		"MintNFTStellarAndAddToDB",
		"POST",
		"/nft/mintStellar",
		businessFacades.MintNFTStellar, //Calls the MintNFTStellar in the nftHandler to mint in Stellar
	},
	Route{
		"MintNFTStellarAndAddToDB",
		"POST",
		"/nft/wallet/mint-stellar",
		businessFacades.MintWalletNFTStellar, //Calls the MintNFTStellar in the nftHandler to mint in Stellar
	},
	Route{
		"MintNFTSolanaAndAddToDB",
		"POST",
		"/nft/mintSolana",
		businessFacades.MintNFTSolana, //Calls the MintNFTSolana in the nftHandler to mint in Solana
	},
	Route{
		"MintNFTContractAndAddToDB",
		"POST",
		"/nft/mintcontract",
		businessFacades.MintNFTContract, //Calls the MintNFTContract in the nftHandler to save contract based minted data
	},
	Route{
		"GetMinterForSolana",
		"GET",
		"/nft/minter/{ImageBase64}/{blockchain}",
		businessFacades.RetrieveSolanaMinter, //Calls the minter from Solana after minting
	},
	Route{
		"RetrieveStellarTxn",
		"GET",
		"/nft/gettxn/{ImageBase64}/{blockchain}",
		businessFacades.RetrieveStellarTxn, //get txn based on the stellar minted nft
	},
	Route{
		"RetriveNFTByStatusAndPK",
		"GET",
		"/nft/getnft",
		businessFacades.RetriveNFTByStatusAndPK, //Calls the RetriveNFTByStatusAndPK in the StellarHandler
	},
	Route{
		"GetNFTIssuerAccount",
		"GET",
		"/nft/issueaccount",
		businessFacades.GetNFTIssuerAccount, //Calls the CreateNFTIssuerAccount in the StellarHandler
	},
	Route{
		"UpdateSellingStatus",
		"PUT",
		"/nft/updatesell",
		businessFacades.UpdateSellingStatus, //Calls the UpdateSellingStatus in the StellarHandler
	},
	Route{
		"UpdateBuyingStatus",
		"PUT",
		"/nft/updatebuy",
		businessFacades.UpdateBuyingStatus, //Calls the UpdateBuyingStatus in the StellarHandler
	},
	Route{
		"GetLastNFTByIdentifier",
		"GET",
		"/lastnft/{InitialDistributorPK}",
		businessFacades.GetLastNFTbyIdentifier, //Calls the GetLastNFTByIdentifier in the StellarHandler
	},
	Route{
		"FundAndGetAccount",
		"GET",
		"/nft/fundAccount",
		businessFacades.FundAndGetAccount, //Calls the FundAccount in the StellarHandler
	},

	Route{
		"GetSponsorAccountXDR",
		"GET",
		"/nft/sponsor",
		businessFacades.GetSponsorAccountXDR, //Calls the GetSponsorAccountXDR in the StellarHandler
	},
	Route{
		"GetSponsorAccount",
		"GET",
		"/trustnetwork/sponsor",
		businessFacades.SponsorAccount, //Calls the GetSponsorAccountXDR in the StellarHandler
	},
	Route{
		"GetSponsorTrustXDR",
		"GET",
		"/nft/sponsortrust",
		businessFacades.GetSponsorTrustXDR, //Calls the GetSponsorTrustXDR in the StellarHandler
	},
	Route{
		"Storing the social impact equation on blockchain",
		"POST",
		"/socialimapact/buildformula",
		businessFacades.BuildSocialImpactExpertFormula,
	},
	Route{
		"Binding metrics with traceability data and Master data",
		"POST",
		"/socialimapact/metricbinding",
		businessFacades.BindMetric,
	},
	Route{
		"save PGP key",
		"POST",
		"/pgp",
		businessFacades.SavePGPKey,
	},
	Route{
		"retrieve PGP publickey",
		"GET",
		"/pgp/{sha256pk}",
		businessFacades.GetRSAPublicKeyBySHA256PK,
	},
	Route{
		"Save Trust Network User",
		"POST",
		"/trustnetwork",
		businessFacades.SaveTrustNetworkUser,
	},
	Route{
		"Get trust networkuser by ID",
		"GET",
		"/trustnetwork/{id}",
		businessFacades.GetTrustNetWorkUserbyID,
	},
	Route{
		"Get trust networkuser by encryptedpassword",
		"GET",
		"/trustnetwork/validate/{password}",
		businessFacades.GetTrustNetWorkUserbyEncryptedPW,
	},
	Route{
		"Endorse User",
		"PUT",
		"/trustnetwork/endorse",
		businessFacades.EndorseTrustNetworkUser,
	},
	Route{
		"Login User",
		"POST",
		"/trustnetwork/login",
		businessFacades.ValidateTrustNetworkUser,
	},
	Route{
		"Get User endorsement count",
		"GET",
		"/trustnetwork/endorsements/{id}",
		businessFacades.GetTrustNetworkUserEndorsmentCount,
	},
	Route{
		"Get all trust network users",
		"GET",
		"/trustnetwork",
		businessFacades.GetAllTrustNetworkUsers,
	},
	Route{
		"Reset Trust Network pasword",
		"PUT",
		"/trustnetwork/reset",
		businessFacades.Resetpassword,
	},
	Route{
		"Update Password",
		"PUT",
		"/trustnetwork",
		businessFacades.UpdatePassword,
	},

	Route{
		"Saving PGP Account details",
		"POST",
		"/pgp/",
		businessFacades.SavePGPAccount,
	},
	Route{
		"Get PGP Public Key and Username",
		"GET",
		"/pgp/getaccounts/{stellarPublicKey}",
		businessFacades.GetPGPAccountByStellarPK,
	},
	Route{
		"TransferNFTonSolana",
		"POST",
		"/atatransfer",
		businessFacades.TransferNFTS, //Calls the MintNFTSolana in the nftHandler to mint in Solana
	},
	Route{
		"LockAquired",
		"POST",
		"/lock",
		businessFacades.BuyHandlerLock, //Calls the BuyHandlerLock in the nftHandler to queue the buy requests
	},
	Route{
		"GetQueuedByImageBase64AndBlockchain",
		"GET",
		"/queue/{ImageBase64}/{blockchain}/{version}",
		businessFacades.RetrieveQueueData, //Calls the RetrieveQueueData from Queue after being processed
	},
	Route{
		"BlockchainTreeRetrieverWithHashMerkleTree",
		"GET",
		"/pocv4/merkletree/{txn}",
		businessFacades.BlockchainTreeRetreiverWithHashWithMerkleTree,
	},
	Route{
		"SponsorFOAppUserTransactions",
		"GET",
		"/foapp/txn",
		businessFacades.SubmitFOData, //Calls the SponsorFOAppUserTransactions in the StellarHandler
	},
	Route{
		"GetSponsoringAccount",
		"GET",
		"/foapp/{foUser}",
		businessFacades.CreateSponsorer, //Calls the GetSponsoringAccount in the StellarHandler
	},
	Route{
		"ActivateFoUser",
		"GET",
		"/activate",
		businessFacades.ActivateFOUser, //Calls the ActivateFoUser in the StellarHandler
	},
	Route{
		"POC",
		"GET",
		"/tdpdata/txnhash",
		businessFacades.NewCheckPOEV3, // Calls the Deprecated POC for Gateway Version 1, Should call the new CheckPOCV3

	},
	Route{
		"BreakTrustline",
		"POST",
		"/nft/breaktrustline",
		businessFacades.BreakTrustline, //Calls the BrealTrustline in the StellarHandler
	},
	Route{
		"SubmitMarketXDR",
		"POST",
		"/marketxdr",
		businessFacades.SubmitMarketXDR, //Calls the ActivateFoUser in the StellarHandler
	},
}
