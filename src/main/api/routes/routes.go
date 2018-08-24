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
	},
}
