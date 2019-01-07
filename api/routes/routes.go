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
		"/poc/{Txn}",
		businessFacades.CheckPOC,
	},
	Route{
		"FULLPOC",
		"POST",
		"/fullpoc/{Txn}/{PID}",
		businessFacades.CheckFullPOC,
	},
	Route{
		"POE",
		"GET",
		"/poe/{Txn}",
		businessFacades.CheckPOE,
	},
	Route{
		"POG",
		"GET",
		"/pog/{LastTxn}/{POGTxn}/{Identifier}",
		businessFacades.CheckPOG,
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
		businessFacades.SubmitXDR,
	},
	Route{
		"SubmitSplit",
		"POST",
		"/transaction/split",
		businessFacades.SubmitXDR,
	},Route{
		"SubmitGenesis",
		"POST",
		"/transaction/genesis",
		businessFacades.SubmitXDR,
	},Route{
		"SubmitTransformation",
		"POST",
		"/transaction/transformation",
		businessFacades.SubmitXDR,
	},Route{
		"SubmitMerge",
		"POST",
		"/transaction/merge",
		businessFacades.SubmitXDR,
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
	},
}
