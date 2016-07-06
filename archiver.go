package main

import (
	"fmt"
	"net/http"

	api "archiver/api"
	web "archiver/web"
)

func main() {
	// api service
	apiConf := api.NewConfig()
	api.ErrorLog(fmt.Sprintf("OK, api service listen at: %s\n", apiConf.Api.Port))
	go api.MgtCron()
	go func() {
		err := http.ListenAndServe(":"+apiConf.Api.Port, api.NewRouter(apiConf.Api.Uri))
		if err != nil {
			panic(err)
		}
	}()

	// web service
	webConf := web.NewConfig()
	web.ErrorLog(fmt.Sprintf("OK, web service listen at: %s\n", webConf.Web.Port))
	err := http.ListenAndServe(":"+webConf.Web.Port, web.NewRouter(webConf.Web.Uri))
	fmt.Println(err)
}
