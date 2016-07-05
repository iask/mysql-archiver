package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

var Ver = "/archiver"

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
	// for static url
	router.Methods("GET").
		Name("Static").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router
}

var routes = Routes{
	Route{
		"SchdsAdd", "POST", Ver + "/schds", myHandle(AddSchd),
	},
	Route{
		"SchdsDel", "GET", Ver + "/schds/del/{id:[0-9]+}", myHandle(DelSchd),
	},
	Route{
		"SchdsModify", "POST", Ver + "/schds/{id:[0-9]+}", myHandle(UpdateSchd),
	},
	Route{
		"SchdList", "GET", Ver + "/schds", myHandle(ListSchd),
	},
	Route{
		"SchdGet", "GET", Ver + "/schds/{id:[0-9]+}", myHandle(GetSchd),
	},
	Route{
		"SchdRun", "GET", Ver + "/schds/{id:[0-9]+}/{mode:[1-2]}", myHandle(DryRun),
	},
	Route{
		"JobList", "GET", Ver + "/jobs", myHandle(ListJobs),
	},
	Route{
		"JobLog", "GET", Ver + "/jobs/{id:[0-9]+}/log", myHandle(JobLog),
	},
	Route{
		"Crons", "GET", Ver + "/crons", myHandle(ListCron),
	},
	Route{
		"Logout", "GET", Ver + "/logout", myHandle(Logout),
	},
	Route{
		"XboxTags", "GET", Ver + "/tags", myHandle(XboxTags),
	},
}
