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
		"DataHash",
		"POST",
		"/data/{hash}/{secret}/{profileId}/{rootHash}",
		businessFacades.SaveDataHash,
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
}
