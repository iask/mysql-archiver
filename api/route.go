package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Ver is the version of the api
const Ver = "/v1"

//var Ver = API.Uri

//Route is the struct of the route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes is the slice of the route
type Routes []Route

//NewRouter news routers
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"SchdsAdd", "POST", Ver + "/schds", myHandle(AddSchd),
	},
	Route{
		"SchdsDel", "DELETE", Ver + "/schds/{id:[0-9]+}", myHandle(DelSchd),
	},
	Route{
		"SchdsModify", "POST", Ver + "/schds/{id:[0-9]+}", myHandle(UpdateSchd),
	},
	Route{
		"SchdList", "GET", Ver + "/schds", myHandle(ListSchds),
	},
	Route{
		"Schd", "GET", Ver + "/schds/{id:[0-9]+}", myHandle(ListSchds),
	},
	Route{
		"SchdRun", "GET", Ver + "/schds/{id:[0-9]+}/{mode:[1-3]}", myHandle(RunSchd),
	},
	Route{
		"JobList", "GET", Ver + "/jobs", myHandle(ListJobs),
	},
	Route{
		"Jobs8Id", "GET", Ver + "/jobs/{id:[0-9]+}", myHandle(ListJobs),
	},
	Route{
		"JobLog", "GET", Ver + "/jobs/{id:[0-9]+}/log", myHandle(JobLog),
	},
	Route{
		"Crons", "GET", Ver + "/crons", myHandle(ListCron),
	},
}
