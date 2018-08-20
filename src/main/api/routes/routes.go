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
		"GET",
		"/poc/{Txn}/{PID}/{dbTree}",
		businessFacades.CheckPOC,
	},
	Route{
		"POE",
		"GET",
		"/poe/{Txn}/{PID}/{Hash}",
		businessFacades.CheckPOE,
	},
	Route{
		"Transactions",
		"POST",
		"/transaction/type/{TType}",
		businessFacades.Transaction,
	},
}
