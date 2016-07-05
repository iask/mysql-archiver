package api

import (
	. "archiver/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	BROCKER_NAME = "ap"
	SECRET_KEY   = "5930fc63523b7776eb90d0cf310979b0"
)

//
type LogOutPut struct {
	Status string            `json:"status"`
	Result map[string]string `json:"result"`
}

// jobs list json output
type CronOutPut struct {
	Status string `json:"status"`
	Result []Cron `json:"result"`
}

// scheduled list json output
type SchdOutPut struct {
	Status string      `json:"status"`
	Result []Scheduled `json:"result"`
}

// jobs list json output
type JobsOutPut struct {
	Status string `json:"status"`
	Result []Job  `json:"result"`
}

//
type JsonOutPut struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

//
type DbMapOutput struct {
	Err string    `json:"Err"`
	Clu []Cluster `json:"Clu"`
}

type schdFields struct {
	Fields []schdField `json:"fields"`
}

/*
 * list all onload cron
 */
func ListCron(w http.ResponseWriter, r *http.Request) {
	var err error
	var ret CronOutPut

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer list cron Error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		err = r.ParseForm()
		CheckErr(err)

		ret.Result, err = listCron()
		if err != nil {
			CheckErr(fmt.Errorf("Get cron list failed"))
		}
		ret.Status = "OK"
		m, _ := json.Marshal(ret)
		w.Write(m)
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 * POST data json format
 * {
 *      "id":n,
 * }
 */
func JobLog(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		j   Job
		ret LogOutPut
	)

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer ReadJobLog() error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			//loginfo["stdout"], loginfo["stderr"] = stdoutlog, stderrlog
			//ret.Result = loginfo
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		jobId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			jobId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		j.Id = int64(jobId)
		err = j.Get()
		if err != nil || j.Id == 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Job info Failed. %s", err)))
		}
		//ErrorLog(j)

		loginfo, err := j.Logs()
		CheckErr(err)

		ret.Status = "OK"
		ret.Result = loginfo
		Msg, _ := json.Marshal(ret)
		w.Write(Msg)
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}

	AccessLog(r, http.StatusOK)
	return
}

/*
 * list all jobs by paging
 */
func ListJobs(w http.ResponseWriter, r *http.Request) {
	var err error
	var ret JobsOutPut

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer List Jobs Error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		// get one jobs
		jobId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			jobId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}

		// get all jobs
		err = r.ParseForm()
		CheckErr(err)
		page, err := strconv.Atoi(r.Form.Get("page"))
		if err != nil {
			page = 1
		}
		if page <= 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("page=< 0 or page>%d", MAX_PAGE)))
		} else if page > MAX_PAGE {
			page = 1
		}

		offset, err := strconv.Atoi(r.Form.Get("offset"))
		if err != nil {
			offset = MAX_OFFSET
		}
		if offset <= 0 && offset > MAX_OFFSET {
			CheckErr(fmt.Errorf(fmt.Sprintf("offset =< 0 or offset >%d", MAX_OFFSET)))
		} else if offset > MAX_OFFSET {
			offset = MAX_OFFSET
		}

		ret.Result, err = listJobs(jobId, page, offset)
		if err != nil {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Jobs List failed, page=%d, %s", page, err)))
		}
		ret.Status = "OK"
		m, _ := json.Marshal(ret)
		w.Write(m)
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)

	return
}

/*
 * add a scheduled
 * {
 *       "name":"miui_test_6",
 *       "host": "lg-dba-cc01.bj",
 *       "port": 3306,
 *       "db": "dba_admin",
 *       "table": "t4",
 *       "target_type":"1",
 *       "target_name": "",
 *       "query": "1=1",
 *       "cron": "0 50 11 * * *",
 *       "deadline":"2016-05-01 00:00:00",
 *       "weight": 5,
 *       "charset": "utf8"
 * }
 */
func AddSchd(w http.ResponseWriter, r *http.Request) {
	var ret SchdOutPut

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer AddSchd Error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "POST" {
		PostData, err := ioutil.ReadAll(r.Body)
		CheckErr(err)
		ErrorLog(fmt.Sprintf("%s", PostData))
		s := Scheduled{}
		err = json.Unmarshal(PostData, &s)
		CheckErr(err)
		ErrorLog(fmt.Sprintf("%v", s))
		ErrorLog(fmt.Sprintf("======> Add schedule start sid: %d <======\n", s.Id))
		if len(s.QueryStr) == 0 {
			ret.Status = "Err: Qurey string do not empty."
		} else {
			err = s.Add()
			if err != nil {
				CheckErr(fmt.Errorf(fmt.Sprintf("Add Scheduled Failed. %s", err)))
			}
			err = s.Get()
			if err != nil {
				CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. %s", err)))
			}
			if s.Id == 0 {
				CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. id=0")))
			} else {
				ret.Status = "OK"
			}
		}

		ret.Result = append(ret.Result, s)
		m, _ := json.Marshal(ret)
		w.Write(m)
		ErrorLog(fmt.Sprintf("======> Add schedule end sid: %d <======\n", s.Id))
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept POST Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 * list all scheduled
 */
func ListSchds(w http.ResponseWriter, r *http.Request) {
	var err error
	var ret SchdOutPut
	var schds []Scheduled

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer ListSchd Error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
		schdId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			schdId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		schds, err = listSchd(schdId)
		CheckErr(err)
		ret.Status, ret.Result = "OK", schds
		m, _ := json.Marshal(ret)
		w.Write(m)
	} else {
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)

	return
}

/*
 * POST data jsoin format
 * {
 *      "id":1,
 *      "runmode":1
 * }
 */
func RunSchd(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		s   Scheduled
		ret JsonOutPut
	)

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer Run Schd Error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "GET" {
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
		//fmt.Printf("schdid=%d,mode=%d\n", schdId, runMode)

		//ErrorLog(fmt.Sprintf("======> Run schedule start sid: %d <======\n", schdId))
		s.Id = schdId
		err = s.Get()
		if err != nil || s.Id == 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. %s", err)))
		}

		// do job
		err = s.Run(runMode)
		if err != nil {
			ret.Status = "Err"
			ret.Result = fmt.Sprintf("%s", err)
		} else {
			ret.Status = "OK"
			ret.Result = "Succ"
		}
		m, _ := json.Marshal(ret)
		w.Write(m)
		//ErrorLog(fmt.Sprintf("======> Run schedule start sid: %d <======\n", s.Id))
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept GET Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 * POST data json format
 * {
 *      "id":n,
 *      "fields" : [
 *          {"type":"string", "fields":"deadline", "data":"2016-05-10 121212"},
 *          {"type":"string", "fields":"deadline", "data":"2016-05-10 121212"},
 *      ]
 * }
 */
func UpdateSchd(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		s   Scheduled
		ret SchdOutPut
	)
	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer Update Schd Error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "POST" {
		schdId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			schdId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}
		s.Id = schdId

		schdJson, err := ioutil.ReadAll(r.Body)
		CheckErr(err)
		ErrorLog(fmt.Sprintf("postdata:%T %s\n", schdJson, schdJson))
		s1 := Scheduled{}
		err = json.Unmarshal(schdJson, &s1)
		CheckErr(err)
		ErrorLog(fmt.Sprintf("sfs: %T %v\n", s1, s1))

		err = s.Get()
		if err != nil || s.Id == 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. %s", err)))
		}

		// do update scheduled
		_, err = s.Update(s1)
		if err != nil {
			CheckErr(fmt.Errorf(fmt.Sprintf("Update Scheduled by id %d Failed. %s", s.Id, err)))
		}

		err = s.Get()
		if err != nil || s.Id == 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. %s", err)))
		}
		ret.Status = "OK"
		ret.Result = append(ret.Result, s)
		m, _ := json.Marshal(ret)
		w.Write(m)
		ErrorLog(fmt.Sprintf("======> Update schedule end sid: %d <======\n", s.Id))
	} else if r.Method == "PATCH" {
		vars := mux.Vars(r)
		fmt.Printf("patch: %v\n", vars)
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept POST Method"))
	}
	AccessLog(r, http.StatusOK)
	return
}

/*
 * POST data json format
 * {
 *      "schd_id":n,
 *      ]
 * }
 */
func GetSchd(w http.ResponseWriter, r *http.Request) {
	var ret SchdOutPut

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer Get Schd Error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "POST" {
		//ErrorLog("/* Get schedule start */")
		PostData, err := ioutil.ReadAll(r.Body)
		CheckErr(err)
		//ErrorLog(fmt.Sprintf("postdata:%T %s\n", PostData, PostData))
		s := Scheduled{}
		err = json.Unmarshal(PostData, &s)
		CheckErr(err)
		if s.Id <= 0 {
			http.Error(w, fmt.Sprintf("Sorry: Bad Request Args, [id] %d", s.Id), http.StatusBadRequest)
			AccessLog(r, http.StatusBadRequest)
			return
		}
		err = s.Get()
		if err != nil || s.Id == 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. %s", err)))
		}
		//ErrorLog(s)
		err = s.Get()
		if err != nil || s.Id == 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. %s", err)))
		}

		ret.Status = "OK"
		ret.Result = append(ret.Result, s)
		m, _ := json.Marshal(ret)
		w.Write(m)
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept POST Method"))
	}

	AccessLog(r, http.StatusOK)
	return
}

/*
 * POST data json format
 * {
 *      "schd_id":n,
 * }
 */
func DelSchd(w http.ResponseWriter, r *http.Request) {
	var err error
	var ret JsonOutPut

	defer func() {
		if err := recover(); err != nil {
			// error handling
			ErrorLog(fmt.Sprintf("Defer drop schd error: %s\n", err))
			ret.Status = fmt.Sprintf("Err: %s", err)
			Msg, _ := json.Marshal(ret)
			http.Error(w, string(Msg), http.StatusInternalServerError)
			AccessLog(r, http.StatusInternalServerError)
			return
		}
	}()
	if r.Method == "DELETE" {
		schdId := 0
		vars := mux.Vars(r)
		if len(vars["id"]) > 0 {
			schdId, err = strconv.Atoi(vars["id"])
			CheckErr(err)
		}

		if schdId <= 0 {
			CheckErr(fmt.Errorf(fmt.Sprintf("Get Scheduled info Failed. %s", schdId)))
		}

		s := Scheduled{Id: schdId}
		err = s.Get()
		if err != nil {
			CheckErr(fmt.Errorf(fmt.Sprintf("Check Scheduled info Failed. %s", err)))
		}
		ErrorLog(s)

		// do update scheduled
		//ErrorLog(fmt.Sprintf("======> Drop schedule start sid: %d <======\n", s.Id))
		err = s.Del()
		if err != nil {
			CheckErr(fmt.Errorf(fmt.Sprintf("Delete Scheduled by id %d Failed. %s", s.Id, err)))
		}

		ret.Status, ret.Result = "OK", fmt.Sprintf("Success drop scheduled by id: %d.'}", s.Id)
		m, _ := json.Marshal(ret)
		w.Write(m)
		//ErrorLog(fmt.Sprintf("======> Drop schedule End sid: %d <======\n", s.Id))
	} else {
		// error handling
		CheckErr(fmt.Errorf("Sorry: Only Accept DELETE Method"))
	}

	AccessLog(r, http.StatusOK)
	return
}
