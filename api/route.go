package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

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
func NewRouter(Uri string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(Uri + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"SchdsAdd", "POST", "/schds", myHandle(AddSchd),
	},
	Route{
		"SchdsDel", "DELETE", "/schds/{id:[0-9]+}", myHandle(DelSchd),
	},
	Route{
		"SchdsModify", "POST", "/schds/{id:[0-9]+}", myHandle(UpdateSchd),
	},
	Route{
		"SchdList", "GET", "/schds", myHandle(ListSchds),
	},
	Route{
		"Schd", "GET", "/schds/{id:[0-9]+}", myHandle(ListSchds),
	},
	Route{
		"SchdRun", "GET", "/schds/{id:[0-9]+}/{mode:[1-3]}", myHandle(RunSchd),
	},
	Route{
		"JobList", "GET", "/jobs", myHandle(ListJobs),
	},
	Route{
		"Jobs8Id", "GET", "/jobs/{id:[0-9]+}", myHandle(ListJobs),
	},
	Route{
		"JobLog", "GET", "/jobs/{id:[0-9]+}/log", myHandle(JobLog),
	},
	Route{
		"Crons", "GET", "/crons", myHandle(ListCron),
	},
}
