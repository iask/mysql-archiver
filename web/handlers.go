package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"archiver/utils"
)

//BodyMaxSize Limit max http request body size
const BodyMaxSize = 1048576

//Response gives the client msg
type Response struct {
	Msg string
}

func response(w http.ResponseWriter, status int, r *Response) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	//	logger.Infof("%+v", r)
	//	logger.Output(2, "DEBUG", r.Message)
	return json.NewEncoder(w).Encode(r)
}

func readRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return ioutil.ReadAll(io.LimitReader(r.Body, BodyMaxSize))
}

func httpStatus(w http.ResponseWriter, r *http.Request) {
	//FIXME: Need more status
	http.Error(w, "Status OK", http.StatusOK)
}

func myHandle(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if BASIC_AUTH.State {
			err := utils.BasicAuth(w, r, BASIC_AUTH.User, BASIC_AUTH.Pass)
			if err != nil {
				ErrorLog(fmt.Sprintf("Basic Auth Error: %s\n", err))
				return
			}
		}
		f(w, r)
	}
}
