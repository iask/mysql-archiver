package web

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
	// for static url
	router.Methods("GET").
		Name("Static").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router
}

var routes = Routes{
	Route{
		"SchdsAdd", "POST", "/schds", myHandle(AddSchd),
	},
	Route{
		"SchdsDel", "GET", "/schds/del/{id:[0-9]+}", myHandle(DelSchd),
	},
	Route{
		"SchdsModify", "POST", "/schds/{id:[0-9]+}", myHandle(UpdateSchd),
	},
	Route{
		"SchdList", "GET", "/schds", myHandle(ListSchd),
	},
	Route{
		"SchdGet", "GET", "/schds/{id:[0-9]+}", myHandle(GetSchd),
	},
	Route{
		"SchdRun", "GET", "/schds/{id:[0-9]+}/{mode:[1-2]}", myHandle(DryRun),
	},
	Route{
		"JobList", "GET", "/jobs", myHandle(ListJobs),
	},
	Route{
		"JobLog", "GET", "/jobs/{id:[0-9]+}/log", myHandle(JobLog),
	},
	Route{
		"Crons", "GET", "/crons", myHandle(ListCron),
	},
	Route{
		"Logout", "GET", "/logout", myHandle(Logout),
	},
	Route{
		"XboxTags", "GET", "/tags", myHandle(XboxTags),
	},
}
