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

//Routes An Array of type Route
type Routes []Route

/*routes contains all the routes
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
		businessFacades.CheckPOCV3, //Calls the Deprecated POC for Gateway Version 1, Should call the new CheckPOCV3
	},
	Route{
		"FULLPOC",
		"POST",
		"/proof/fullpoc/{Txn}",
		businessFacades.CheckFullPOC, //Calls the Deprecated FULLPOC for Gateway Version 1
	},
	Route{
		"POE",
		"GET",
		"/proof/poe/{Txn}",
		businessFacades.CheckPOEV3, //Calls the Functional POE for Gateway Version 3
	},
	Route{
		"POG",
		"GET",
		"/proof/pog/{Txn}",
		businessFacades.CheckPOGV3Rewrite, //Calls the Functional POG for Gateway Version 3
	},
	Route{
		"POCOC",
		"GET",
		"/proof/pococ/{TxnId}",
		businessFacades.CheckPOCOCV3, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"GetCOCCollectionBySender",
		"get",
		"/getcocbysender/{Sender}",
		businessFacades.GetCocBySender, //Calls the Functional GetCocBySender for Gateway Version 3
	},
	Route{
		"GetCOCCollectionByReceiver",
		"get",
		"/getcocbyreceiver/{Receiver}",
		businessFacades.GetCocByReceiver, //Calls the Functional GetCocByReceiver for Gateway Version 3
	},
	Route{
		"InsertCOCCollection",
		"POST",
		"/insertcoccollection",
		businessFacades.InsertCocCollection, //Calls the Functional InsertCocCollection for Gateway Version 3
	},
	Route{
		"InsertCOCCollection",
		"PUT",
		"/insertcoccollection",
		businessFacades.UpdateCocCollection, //Calls the Functional UpdateCocCollection for Gateway Version 3
	},
	Route{
		"SubmitXDR",
		"POST",
		"/transaction/dataPacket",
		businessFacades.SubmitData, //Calls the Functional SubmitData for Gateway Version 3
	},
	Route{
		"SubmitSplit",
		"POST",
		"/transaction/split",
		businessFacades.SubmitSplit, //Calls the Functional SubmitSplit for Gateway Version 3
	},
	Route{
		"SubmitGenesis",
		"POST",
		"/transaction/genesis",
		businessFacades.SubmitGenesis, //Calls the Functional SubmitGenesis for Gateway Version 3
	},
	Route{
		"SubmitTransformation",
		"POST",
		"/transaction/transformation",
		businessFacades.SubmitTransformation, //Calls the Functional SubmitTransformation for Gateway Version 3
	},
	Route{
		"SubmitMerge",
		"POST",
		"/transaction/merge",
		businessFacades.SubmitMerge, //Calls the Functional SubmitMerge for Gateway Version 3
	},
	Route{
		"SubmitTransfer",
		"POST",
		"/transaction/transfer",
		businessFacades.SubmitTransfer, //Calls the Functional SubmitTransfer for Gateway Version 3
	},
	Route{
		"InsertCertificate",
		"POST",
		"/transaction/certificateInsert",
		businessFacades.SubmitCertificateInsert, //Calls the Functional SubmitCertificateInsert for Gateway Version 3
	},
	Route{
		"RenewCertificate",
		"POST",
		"/transaction/certificateRenew",
		businessFacades.SubmitCertificateRenewal, //Calls the Functional SubmitCertificateRenewal for Gateway Version 3
	},
	Route{
		"RevokeCertificate",
		"POST",
		"/transaction/certificateRevoke",
		businessFacades.SubmitCertificateRevoke, //Calls the Functional SubmitCertificateRevoke for Gateway Version 3
	},
	Route{
		"LastTxn",
		"GET",
		"/transaction/lastTxn/{Identifier}",
		businessFacades.LastTxn, //Calls the Functional LastTxn for Gateway Version 3
	},
	Route{
		"SubAccountStatus",
		"POST",
		"/transaction/coc/subAccountStatus",
		businessFacades.CheckAccountsStatus, //Calls the Functional CheckAccountsStatus for Gateway Version 3
	},
	Route{
		"POCDeveloperRetriever",
		"get",
		"/pocbctree/{Txn}",
		businessFacades.DeveloperRetriever, //Test
	},
	Route{
		"GET LOGS",
		"get",
		"/getLogsForToday/",
		businessFacades.RetrieveLogsForToday, //Test
	},
	Route{
		"POCGatewayRetrieverForTDP",
		"GET",
		"/gatewayTree/{Txn}",
		businessFacades.GatewayRetriever, //Test
	},
	Route{
		"POCGatewayRetrieverForIdentifier",
		"GET",
		"/gatewayTreeWithIdentifier/{Identifier}",
		businessFacades.GatewayRetrieverWithIdentifier, //Test
	},
	Route{
		"ConvertXDRToTXN",
		"POST",
		"/xdrToTxn",
		businessFacades.ConvertXDRToTXN, //Test
	},
	Route{
		"LastCOC",
		"GET",
		"/lastCoc/{Identifier}",
		businessFacades.LastCOC, //Test
	},
	Route{
		"Retrieve TDP for Transaction",
		"GET",
		"/tdpForTxn/{Txn}",
		businessFacades.TDPForTXN, //Test
	}, Route{
		"Retrieve Transaction for TDP",
		"GET",
		"/txnForTdp/{Txn}",
		businessFacades.TXNForTDP, //Test
	},
	Route{
		"Transactions",
		"POST",
		"/transaction/type/{TType}",
		businessFacades.Transaction, //Deprecated
	},
	Route{
		"TrustLine",
		"POST",
		"/create/Trustline",
		businessFacades.CreateTrust, //Deprecated
	},
	Route{
		"SendAssestV2",
		"POST",
		"/send/asset",
		businessFacades.SendAssests, //Deprecated
	},
	Route{
		"lockAcc",
		"POST",
		"/lock/registrarAcc",
		businessFacades.MultisigAccount, //Deprecated
	},
	Route{
		"UnlockAcc",
		"POST",
		"/Appoint/Registrar",
		businessFacades.AppointRegistrar, //Deprecated
	},
	Route{
		"transformV2",
		"POST",
		"/transform/V2",
		businessFacades.TransformV2, //Deprecated
	},
	Route{
		"COC",
		"POST",
		"/COC/Transaction",
		businessFacades.COC, //Deprecated
	},
	Route{
		"COCLink",
		"POST",
		"/COCLink/Transaction",
		businessFacades.COCLink, //Deprecated
	},
	Route{
		"TransactionId",
		"GET",
		"/TransactionId/{id}",
		businessFacades.GetTransactionId, //Test
	},
	Route{
		"TransactionIds",
		"GET",
		"/GetTransactionsForTDP/{id}",
		businessFacades.GetTransactionsForTDP, //Test
	},
	Route{
		"TransactionIdsForTDPs",
		"POST",
		"/GetTransactionsForTDPs",
		businessFacades.GetTransactionsForTdps, //Test
	},
	Route{
		"TransactionIdsForPK",
		"GET",
		"/GetTransactionsForPK/{id}",
		businessFacades.GetTransactionsForPK, //Test
	},
	Route{
		"RetriveTransactionId",
		"GET",
		"/GetTransactionId/{id}",
		businessFacades.RetriveTransactionId, //Test
	},
	Route{
		"QueryTransactionsByKey",
		"GET",
		"/GetTransactions",
		businessFacades.QueryTransactionsByKey, //multisearch
	},
	Route{
		"GetCOCByTxn",
		"GET",
		"/GetCOCByTxn/{txn}",
		businessFacades.GetCOCByTxn, //multisearch
	},
	Route{
		"RetriveTransactionId",
		"GET",
		"/GetTransactionId/{id}",
		businessFacades.RetriveTransactionId, //Test
	},
	Route{
		"RetrievePreviousTranasctions",
		"GET",
		"/RetrievePreviousTranasctions",
		businessFacades.RetrievePreviousTranasctions, //Test
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
		businessFacades.ArtifactTransactions, //Test
	},
	Route{
		"InsertOrganization",
		"POST",
		"/organization",
		businessFacades.InsertOrganization, //Test
	},
	Route{
		"GetAllOrganizations",
		"GET",
		"/approved/organization",
		businessFacades.GetAllOrganizations, //Test
	},
	Route{
		"GetOrganizationByPublicKey",
		"GET",
		"/organization/{PK}",
		businessFacades.GetOrganizationByPublicKey, //Test
	},
	Route{
		"UpdateOrganization",
		"PUT",
		"/organization",
		businessFacades.UpdateOrganization, //Test
	},
	Route{
		"InsertTestimonial",
		"POST",
		"/testimonial",
		businessFacades.InsertTestimonial, //Test
	},
	Route{
		"GetTestimonialBySender",
		"GET",
		"/testimonial/sender/{PK}",
		businessFacades.GetTestimonialBySender, //Test
	},
	Route{
		"GetTestimonialByReciever",
		"GET",
		"/testimonial/reciever/{PK}",
		businessFacades.GetTestimonialByReciever, //Test
	},
	Route{
		"UpdateTestimonial",
		"PUT",
		"/testimonial",
		businessFacades.UpdateTestimonial, //Test
	},
	Route{
		"SubAccountStatusExtended",
		"POST",
		"/transaction/subAccountStatus",
		businessFacades.CheckAccountsStatusExtended, //Test
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
		businessFacades.GetProofPresentationProtocolByProofName, //Test - calls the GetProofPresentationProtocolByProofName in the ProofPresesntationHandlers
	},
	Route{
		"InsertProofPresentationProtocol",
		"POST",
		"/proofprotocol",
		businessFacades.InsertProofPresentationProtocol, //Test - calls the InsertProofPresentationProtocol in the ProofPresesntationHandlers
	},
	Route{
		"UpdateProofPresentationProtocol",
		"PUT",
		"/proofprotocol",
		businessFacades.UpdateProofPresentationProtocol, //Test - calls the UpdateProofPresesntationProtocol in the ProofPresesntationHandlers
	},
	Route{
		"DeleteProofPresentationProtocolByProofName",
		"DELETE",
		"/deleteproofprotocol/{proofname}",
		businessFacades.DeleteProofPresentationProtocolByProofName, //Test - calls the DeleteProofPresentationProtocolByProtocolName in the ProofPresesntationHandlers
	},
	Route{
		"GetAllOrganizationsPaginated",
		"GET",
		"/approved/organizationPaginated",
		businessFacades.GetAllOrganizations_Paginated, //Test
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
}
