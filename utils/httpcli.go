package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTP Basic Auth
func BasicAuth(w http.ResponseWriter, r *http.Request, u string, p string) error {
	var err error

	if len(u) == 0 || len(p) == 0 {
		err = fmt.Errorf("Initialize http basic auth failed")
		return err
	}
	user, pass, ok := r.BasicAuth()
	if ok {
		if user != u || pass != p {
			err = fmt.Errorf("Status Forbidden")
			http.Error(w, fmt.Sprintln(err), http.StatusForbidden)
			//AccessLog(r, http.StatusForbidden)
			return err
		}
	} else {
		err = fmt.Errorf("Status Unauthorized")
		http.Error(w, fmt.Sprintln(err), http.StatusUnauthorized)
		//AccessLog(r, http.StatusUnauthorized)
		return err
	}

	return err
}

func HttpGet(url string, data []byte, u string, p string) ([]byte, error) {
	return errorHandler(request("GET", url, data, u, p))
}

func HttpDelete(url string, data []byte, u string, p string) ([]byte, error) {
	return errorHandler(request("DELETE", url, data, u, p))
}

func HttpPost(url string, data []byte, u string, p string) ([]byte, error) {
	return errorHandler(request("POST", url, data, u, p))
}

func errorHandler(data []byte, code int, err error) ([]byte, error) {
	// marathon internal error
	if err == nil && code > 300 {
		return nil, fmt.Errorf("%v", string(data))
	}
	return data, err
}

func request(method, url string, data []byte, u string, p string) ([]byte, int, error) {
	if url == "" {
		return nil, 0, fmt.Errorf("Url is null")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	/* Authenticate */
	if len(u) > 0 && len(p) > 0 {
		req.SetBasicAuth(u, p)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}
