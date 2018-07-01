package routes

import (
	"main/api/businessfacades"
	"net/http"
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
		businessfacades.SaveDataHash,
	},
	Route{
		"Index",
		"GET",
		"/",
		businessfacades.Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		businessfacades.TodoIndex,
	},
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		businessfacades.TodoCreate,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		businessfacades.TodoShow,
	},
}
