package routes

import (
	"net/http"

	"main/api/businessFacades"
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
		"TDPInsert",
		"GET",
		"/TDPInsert/{hash}/{type}/{PreviousTXNID}/{profileId}",
		businessFacades.SaveData,
	},
	Route{
		"POC",
		"POST",
		"/poc/{Txn}/{PID}",
		businessFacades.CheckPOC,
	},
	Route{
		"POC",
		"POST",
		"/fullpoc/{Txn}/{PID}",
		businessFacades.CheckFullPOC,
	},
	Route{
		"POE",
		"GET",
		"/poe/{Txn}/{PID}/{Hash}",
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
	},Route{
		"POCDeveloperRetriever",
		"get",
		"/pocbctree/{Txn}",
		businessFacades.DeveloperRetriever,
	},
}
