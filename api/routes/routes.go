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

type Routes []Route

var routes = Routes{

	Route{
		"POC",
		"GET",
		"/proof/poc/{Txn}",
		businessFacades.CheckPOC,
	},
	Route{
		"FULLPOC",
		"POST",
		"/proof/fullpoc/{Txn}",
		businessFacades.CheckFullPOC,
	},
	Route{
		"POE",
		"GET",
		"/poe/{Txn}",
		businessFacades.CheckPOEV3,
	},
	Route{
		"POG",
		"GET",
		"/pog/{Identifier}",
		businessFacades.CheckPOGV3,
	},
	Route{
		"Transactions",
		"POST",
		"/transaction/type/{TType}",
		businessFacades.Transaction,
	}, Route{
		"TrustLine",
		"POST",
		"/create/Trustline",
		businessFacades.CreateTrust,
	}, Route{
		"SendAssestV2",
		"POST",
		"/send/asset",
		businessFacades.SendAssests,
	}, Route{
		"lockAcc",
		"POST",
		"/lock/registrarAcc",
		businessFacades.MultisigAccount,
	}, Route{
		"UnlockAcc",
		"POST",
		"/Appoint/Registrar",
		businessFacades.AppointRegistrar,
	}, Route{
		"transformV2",
		"POST",
		"/transform/V2",
		businessFacades.TransformV2,
	}, Route{
		"COC",
		"POST",
		"/COC/Transaction",
		businessFacades.COC,
	}, Route{
		"COCLink",
		"POST",
		"/COCLink/Transaction",
		businessFacades.COCLink,
	}, Route{
		"POCDeveloperRetriever",
		"get",
		"/pocbctree/{Txn}",
		businessFacades.DeveloperRetriever,
	},
	Route{
		"GetCOCCollectionBySender",
		"get",
		"/getcocbysender/{Sender}",
		businessFacades.GetCocBySender,
	}, Route{
		"GetCOCCollectionByReceiver",
		"get",
		"/getcocbyreceiver/{Receiver}",
		businessFacades.GetCocByReceiver,
	},
	Route{
		"InsertCOCCollection",
		"POST",
		"/insertcoccollection",
		businessFacades.InsertCocCollection,
	},
	Route{
		"InsertCOCCollection",
		"PUT",
		"/insertcoccollection",
		businessFacades.UpdateCocCollection,
	},
	Route{
		"SubmitXDR",
		"POST",
		"/transaction/dataPacket",
		businessFacades.SubmitData,
	},
	Route{
		"SubmitSplit",
		"POST",
		"/transaction/split",
		businessFacades.SubmitSplit,
	},Route{
		"SubmitGenesis",
		"POST",
		"/transaction/genesis",
		businessFacades.SubmitGenesis,
	},Route{
		"SubmitTransformation",
		"POST",
		"/transaction/transformation",
		businessFacades.SubmitTransformation,
	},Route{
		"SubmitMerge",
		"POST",
		"/transaction/merge",
		businessFacades.SubmitMerge,
	},Route{
		"SubmitTransfer",
		"POST",
		"/transaction/transfer",
		businessFacades.SubmitTransfer,
	},Route{
		"LastTxn",
		"GET",
		"/transaction/lastTxn/{Identifier}",
		businessFacades.LastTxn,
	},Route{
		"SubAccountStatus",
		"POST",
		"/transaction/coc/subAccountStatus",
		businessFacades.CheckAccountsStatus,
	}, Route{
		"POCGatewayRetrieverForTDP",
		"GET",
		"/gatewayTree/{Txn}",
		businessFacades.GatewayRetriever,
	},
	Route{
		"POCGatewayRetrieverForIdentifier",
		"GET",
		"/gatewayTreeWithIdentifier/{Identifier}",
		businessFacades.GatewayRetrieverWithIdentifier,
	},Route{
		"ConvertXDRToTXN",
		"POST",
		"/xdrToTxn",
		businessFacades.ConvertXDRToTXN,
	},Route{
		"LastCOC",
		"GET",
		"/lastCoc/{Identifier}",
		businessFacades.LastCOC,
	},Route{
		"Retrieve TDP for Transaction",
		"GET",
		"/tdpForTxn/{Txn}",
		businessFacades.TDPForTXN,
	},
}
