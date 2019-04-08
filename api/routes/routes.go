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
		"POC",
		"GET",
		"/proof/poc/{Txn}",
		businessFacades.CheckPOC, //Calls the Deprecated POC for Gateway Version 1, Should call the new CheckPOCV3
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
		"/proof/pog/{Identifier}",
		businessFacades.CheckPOGV3, //Calls the Functional POG for Gateway Version 3
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
		businessFacades.GetCocBySender, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"GetCOCCollectionByReceiver",
		"get",
		"/getcocbyreceiver/{Receiver}",
		businessFacades.GetCocByReceiver, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"InsertCOCCollection",
		"POST",
		"/insertcoccollection",
		businessFacades.InsertCocCollection, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"InsertCOCCollection",
		"PUT",
		"/insertcoccollection",
		businessFacades.UpdateCocCollection, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"SubmitXDR",
		"POST",
		"/transaction/dataPacket",
		businessFacades.SubmitData, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"SubmitSplit",
		"POST",
		"/transaction/split",
		businessFacades.SubmitSplit, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"SubmitGenesis",
		"POST",
		"/transaction/genesis",
		businessFacades.SubmitGenesis, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"SubmitTransformation",
		"POST",
		"/transaction/transformation",
		businessFacades.SubmitTransformation, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"SubmitMerge",
		"POST",
		"/transaction/merge",
		businessFacades.SubmitMerge, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"SubmitTransfer",
		"POST",
		"/transaction/transfer",
		businessFacades.SubmitTransfer, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"InsertCertificate",
		"POST",
		"/transaction/certificateInsert",
		businessFacades.SubmitCertificateInsert, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"RenewCertificate",
		"POST",
		"/transaction/certificateRenew",
		businessFacades.SubmitCertificateRenewal, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"RevokeCertificate",
		"POST",
		"/transaction/certificateRevoke",
		businessFacades.SubmitCertificateRevoke, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"LastTxn",
		"GET",
		"/transaction/lastTxn/{Identifier}",
		businessFacades.LastTxn, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"SubAccountStatus",
		"POST",
		"/transaction/coc/subAccountStatus",
		businessFacades.CheckAccountsStatus, //Calls the Functional POCOC for Gateway Version 3
	},
	Route{
		"POCDeveloperRetriever",
		"get",
		"/pocbctree/{Txn}",
		businessFacades.DeveloperRetriever, //Test
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
}
