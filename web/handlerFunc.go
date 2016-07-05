package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	. "archiver/utils"

	"github.com/gorilla/mux"
)

// scheduled list json output
type SchdOutput struct {
	Status string      `json:"status"`
	Result []Scheduled `json:"result"`
}

// jobs list json output
type JobsOutPut struct {
	Status string `json:"status"`
	Result []Job  `json:"result"`
}

// jobs list json output
type CronOutPut struct {
	Status string `json:"status"`
	Result []Cron `json:"result"`
}

// jobs list json output
type TagsOutPut struct {
	Status string `json:"status"`
	Result []Tags `json:"result"`
}

func getUser(r *http.Request) string {
	c, err := r.Cookie("user")
	if err == nil {
		return c.Value
	}
	return ""
}

/*
 * add a scheduled
 * {
 *	 "name":"miui_test_6",
 *	 "xboxtag": "cop.xxx.xxx_pdl.xxx_service.xxx",
 *	 "port": 3306,
 *	 "db": "dba_admin",
 *	 "table": "t4",
 *	 "target_type":"1",
 *	 "query": "1=1",
 *	 "cron": "0 50 11 * * *",
 *	 "deadline":"2016-05-01 00:00:00",
 *	 "weight": 5,
 *	 "charset": "utf8"
 * }
 */
func AddSchd(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ErrorLog(fmt.Sprintf("Web Add Schd Error: %s\n", err))
			errMsg := fmt.Sprintf("Web Add Schd Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "POST" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		ErrorLog("/* Add schedule start */")
		PostData, err := ioutil.ReadAll(r.Body)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("%s\n", PostData))
		s := Scheduled{}
		b := strings.Split(string(PostData), "&")
		for _, v := range b {
			c := strings.Split(v, "=")
			switch c[0] {
			case "name":
				s.Name, _ = url.QueryUnescape(c[1])
			case "xboxtag":
				s.XboxTag, _ = url.QueryUnescape(c[1])
			case "port":
				s.Port, _ = strconv.Atoi(c[1])
			case "db":
				s.Db, _ = url.QueryUnescape(c[1])
			case "table":
				s.Table, _ = url.QueryUnescape(c[1])
			case "query":
				s.QueryStr, _ = url.QueryUnescape(c[1])
			case "target_type":
				s.TargetType, _ = strconv.Atoi(c[1])
			case "charset":
				s.Charset, _ = url.QueryUnescape(c[1])
			case "weight":
				s.Weight, _ = strconv.Atoi(c[1])
			case "deadline":
				s.Deadline, _ = url.QueryUnescape(c[1])
			case "cron":
				s.Cron, _ = url.QueryUnescape(c[1])
			default:
				continue
			}
		}
		//ErrorLog(fmt.Sprintf("%v\n", s))
		schdJson, _ := json.Marshal(s)
		//ErrorLog(fmt.Sprintf("%s\n", schdJson))

		ApiUrl := fmt.Sprintf("%s/schds", WEB.ApiUrl)
		ret, err := HttpPost(ApiUrl, schdJson, BASIC_AUTH.User, BASIC_AUTH.Pass)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("%s\n", ret))
		w.Write(ret)

		ErrorLog("/* Add schedule end */")
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept POST Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 * update a scheduled
 * {
 *	 "name":"miui_test_6",
 *	 "target_type":"1",
 *	 "query": "1=1",
 *	 "cron": "0 50 11 * * *",
 *	 "deadline":"2016-05-01 00:00:00",
 *	 "weight": 5,
 *	 "charset": "utf8",,,,
 *	 "active": 0
 * }
 */
func UpdateSchd(w http.ResponseWriter, r *http.Request) {
	var err error
	var s Scheduled

	defer func() {
		if err := recover(); err != nil {
			ErrorLog(fmt.Sprintf("Web Update Schd Error: %s\n", err))
			errMsg := fmt.Sprintf("Web Update Schd Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "POST" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		schdId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			schdId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		s.Id = schdId

		PostData, err := ioutil.ReadAll(r.Body)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("update schd data: %s\n", PostData))
		pd := strings.Split(string(PostData), "&")
		for _, v := range pd {
			c := strings.Split(v, "=")
			switch c[0] {
			case "name":
				s.Name, _ = url.QueryUnescape(c[1])
			case "query":
				s.QueryStr, _ = url.QueryUnescape(c[1])
			case "target_type":
				s.TargetType, _ = strconv.Atoi(c[1])
			case "weight":
				s.Weight, _ = strconv.Atoi(c[1])
			case "deadline":
				s.Deadline, _ = url.QueryUnescape(c[1])
			case "cron":
				s.Cron, _ = url.QueryUnescape(c[1])
			case "charset":
				s.Charset, _ = url.QueryUnescape(c[1])
			case "active":
				s.Active, _ = strconv.Atoi(c[1])
			default:
				continue
			}
		}
		//ErrorLog(fmt.Sprintf("schd: %v\n", s))
		schdJson, _ := json.Marshal(s)
		//ErrorLog(fmt.Sprintf("%s\n", schdJson))

		ErrorLog("/* Update schedule start */")
		ApiUrl := fmt.Sprintf("%s/schds/%d", WEB.ApiUrl, s.Id)
		ret, err := HttpPost(ApiUrl, schdJson, BASIC_AUTH.User, BASIC_AUTH.Pass)
		CheckErr(err)
		ErrorLog("/* Update schedule end */")
		//ErrorLog(fmt.Sprintf("%s\n", ret))

		w.Write(ret)

	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept POST Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 *
 */
func DelSchd(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if err := recover(); err != nil {
			ErrorLog(fmt.Sprintf("Web Delete Schd Error: %s\n", err))
			errMsg := fmt.Sprintf("Web Delete Schd Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		schdId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			schdId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		ErrorLog("/* Del schedule start */")
		ApiUrl := fmt.Sprintf("%s/schds/%s", WEB.ApiUrl, schdId)
		ret, err := HttpDelete(ApiUrl, nil, BASIC_AUTH.User, BASIC_AUTH.Pass)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("ret: %s\n", ret))
		w.Write(ret)
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept DELETE Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 *
 */
func GetSchd(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if err := recover(); err != nil {
			ErrorLog(fmt.Sprintf("Web Get Schd Error: %s\n", err))
			errMsg := fmt.Sprintf("Web Get Schd Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		schdId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			schdId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		ApiUrl := fmt.Sprintf("%s/schds/%d", WEB.ApiUrl, schdId)
		ret, err := HttpGet(ApiUrl, nil, BASIC_AUTH.User, BASIC_AUTH.Pass)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("ret: %s\n", ret))
		w.Write(ret)
	} else {
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 * list all scheduled
 */
func ListSchd(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Web list schd Error: %s\n", err))
			errMsg := fmt.Sprintf("Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		ApiUrl := fmt.Sprintf("%s/schds", WEB.ApiUrl)
		retJson, err := HttpGet(ApiUrl, nil, BASIC_AUTH.User, BASIC_AUTH.Pass)
		ret := SchdOutput{}
		err = json.Unmarshal(retJson, &ret)
		CheckErr(err)
		if ret.Status != "OK" {
			CheckErr(fmt.Errorf("%s", ret.Status))
		}

		t, err := template.ParseFiles(WEB.TemplatePath+"/listtask.html", WEB.TemplatePath+"/header.html", WEB.TemplatePath+"/footer.html")
		CheckErr(err)
		err = t.Execute(w, ret)
		CheckErr(err)
	} else {
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}

	AccessLog(r, http.StatusOK)
	return
}

/*
 * list top 50 jobs
 */
func ListJobs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Web list jobs Error: %s\n", err))
			errMsg := fmt.Sprintf("Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()

	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		ret := JobsOutPut{}
		ApiUrl := fmt.Sprintf("%s/jobs", WEB.ApiUrl)
		retJson, err := HttpGet(ApiUrl, nil, BASIC_AUTH.User, BASIC_AUTH.Pass)
		err = json.Unmarshal(retJson, &ret)
		CheckErr(err)
		//ErrorLog(ret)
		if ret.Status != "OK" {
			CheckErr(fmt.Errorf("%s", ret.Status))
		}

		t, err := template.ParseFiles(WEB.TemplatePath+"/listjobs.html", WEB.TemplatePath+"/header.html", WEB.TemplatePath+"/footer.html")
		CheckErr(err)
		//err = templates.Execute(w, ret)
		err = t.Execute(w, ret)
		CheckErr(err)

	} else {
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 *
 */
func JobLog(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if err := recover(); err != nil {
			ErrorLog(fmt.Sprintf("Web read job log error: %s\n", err))
			errMsg := fmt.Sprintf("Web read job log error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		jobId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			jobId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		ApiUrl := fmt.Sprintf("%s/jobs/%d/log", WEB.ApiUrl, jobId)
		ret, err := HttpGet(ApiUrl, nil, BASIC_AUTH.User, BASIC_AUTH.Pass)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("ret: %s\n", ret))
		w.Write(ret)
	} else {
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 * list onload cron
 */
func ListCron(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ErrorLog(fmt.Sprintf("Web list cron Error: %s\n", err))
			errMsg := fmt.Sprintf("Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		ret := CronOutPut{}
		ApiUrl := fmt.Sprintf("%s/crons", WEB.ApiUrl)
		retJson, err := HttpGet(ApiUrl, nil, BASIC_AUTH.User, BASIC_AUTH.Pass)
		//fmt.Printf("retjson: %s\n", retJson)
		err = json.Unmarshal(retJson, &ret)
		CheckErr(err)
		if ret.Status != "OK" {
			CheckErr(fmt.Errorf("%s", ret.Status))
		}

		t, err := template.ParseFiles(fmt.Sprintf("%s/listcron.html", WEB.TemplatePath), fmt.Sprintf("%s/header.html", WEB.TemplatePath), fmt.Sprintf("%s/footer.html", WEB.TemplatePath))
		CheckErr(err)
		//err = templates.Execute(w, ret)
		err = t.Execute(w, ret)
		CheckErr(err)
	} else {
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 *
 */
func DryRun(w http.ResponseWriter, r *http.Request) {
	var err error

	defer func() {
		if err := recover(); err != nil {
			ErrorLog(fmt.Sprintf("Web Get Schd Error: %s\n", err))
			errMsg := fmt.Sprintf("Web Get Schd Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		schdId := 0
		runMode := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			schdId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		if len(vars["mode"]) > 0 {
			runMode, err = strconv.Atoi(vars["mode"])
			CheckErr(err)
		}
		ApiUrl := fmt.Sprintf("%s/schds/%d/%d", WEB.ApiUrl, schdId, runMode)
		ret, err := HttpGet(ApiUrl, nil, BASIC_AUTH.User, BASIC_AUTH.Pass)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("ret: %s\n", ret))
		w.Write(ret)
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}

	AccessLog(r, http.StatusOK)
	return
}

/*
 * list all scheduled
 */
func XboxTags(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Web list xbox tag Error: %s\n", err))
			errMsg := fmt.Sprintf("Error: %s", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		if !isLogin(r) {
			Login(w, r)
			return
		}
		x := NewXboxTree()
		username := getUser(r)
		if len(username) <= 0 {
			return
		}
		tags, err := x.Get(username)
		CheckErr(err)
		ret := TagsOutPut{"OK", tags}
		m, _ := json.Marshal(ret)
		w.Write(m)
	} else {
		CheckErr(fmt.Errorf("Sorry: Only Accept POST Method"))
	}

	AccessLog(r, http.StatusOK)
	return
}

/*
var templates = template.Must(template.ParseFiles(
	"archiver/template/header.html",
	"archiver/template/footer.html",
	"archiver/template/listtask.html",
	"archiver/template/header.html",
	"archiver/template/footer.html",
	"archiver/template/listjobs.html",
))

*/
