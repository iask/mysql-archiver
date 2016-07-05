package api

import (
	"fmt"
	"strings"
	"time"

	. "archiver/utils"

	"github.com/jakecoffman/cron"
)

// define cron job
type StartJob struct {
	schd Scheduled
	name string
}

//
type Cron struct {
	SchdId     int    `json:"scheduled_id"`
	SchdName   string `json:"scheduled_name"`
	Cron       string `json:"cron"`
	UpdateTime string `json:"update_time"`
	JobId      int    `json:"job_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Status     int    `json:"status"`
	Host       string `json:"dbhost"`
}

//
func listCron() ([]Cron, error) {
	var c Cron
	var cs []Cron

	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return cs, err
	}
	q := "SELECT c.scheduled_id, c.scheduled_name, c.cron, c.update_time, j.job_id, j.start_time, j.end_time, j.status, j.dbhost FROM crontab c JOIN scheduled s ON c.scheduled_id=s.scheduled_id LEFT OUTER JOIN (SELECT scheduled_id AS sid, MAX(job_id) AS jid FROM jobs GROUP BY scheduled_id) mj ON c.scheduled_id=mj.sid LEFT OUTER JOIN jobs j ON mj.jid=j.job_id"
	//ErrorLog(fmt.Sprintf("%s\n", q))
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return cs, err
	}
	for _, row := range rows {
		c.SchdId = row.Int(res.Map("scheduled_id"))
		c.SchdName = row.Str(res.Map("scheduled_name"))
		c.Cron = row.Str(res.Map("cron"))
		c.UpdateTime = row.Str(res.Map("update_time"))
		c.JobId = row.Int(res.Map("job_id"))
		c.StartTime = row.Str(res.Map("start_time"))
		c.EndTime = row.Str(res.Map("end_time"))
		c.Status = row.Int(res.Map("status"))
		c.Host = row.Str(res.Map("dbhost"))
		cs = append(cs, c)
	}

	return cs, nil
}

//
func getDeadSchd() (map[string]string, error) {
	var sid string

	deadjob := make(map[string]string, 0)
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return deadjob, err
	}

	// update active status for dead job
	sid = ""
	q := "SELECT scheduled_id, name FROM scheduled WHERE active=1 AND deadline <= now();"
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return deadjob, err
	} else {
		for _, row := range rows {
			sid += fmt.Sprintf("%d, ", row.Int(res.Map("scheduled_id")))
		}
		sid = strings.Trim(strings.Trim(sid, " "), ",")
		if len(sid) > 0 {
			q := fmt.Sprintf("UPDATE scheduled SET active=0 WHERE scheduled_id in (%s);", sid)
			ErrorLog(fmt.Sprintf("%s\n", q))
			_, _, err := dbadmin.Query(q)
			if err != nil {
				return deadjob, err
			}
		}
	}
	q = "SELECT scheduled_id, name FROM scheduled WHERE active = 0 OR deadline <= now();"
	rows, res, err = dbadmin.Query(q)
	if err != nil {
		ErrorLog(fmt.Sprintf("deadjob err: %v\n", deadjob))
		return deadjob, err
	}
	for _, row := range rows {
		id := fmt.Sprintf("%d", row.Int(res.Map("scheduled_id")))
		name := row.Str(res.Map("name"))
		deadjob[id] = name
	}

	return deadjob, nil
}

//
func getAliveSchd() ([]Scheduled, error) {
	var schds []Scheduled
	var s Scheduled

	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return schds, err
	}
	//ErrorLog(fmt.Sprintf("getAliveSchd inside dbadmin: %v\n", dbadmin))

	q := "SELECT * FROM scheduled WHERE active=1 AND deadline > now();"
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return schds, err
	} else {
		for _, row := range rows {
			s.Id = row.Int(res.Map("scheduled_id"))
			s.Name = row.Str(res.Map("name"))
			s.XboxTag = row.Str(res.Map("xboxtag"))
			s.Port = row.Int(res.Map("port"))
			s.Db = row.Str(res.Map("db"))
			s.Table = row.Str(res.Map("tbl"))
			s.Tables = s.getTables()
			s.TargetType = row.Int(res.Map("target_type"))
			s.TargetName = row.Str(res.Map("target_name"))
			s.QueryStr = row.Str(res.Map("querystr"))
			s.Cron = row.Str(res.Map("cron"))
			s.Deadline = row.Str(res.Map("deadline"))
			s.Weight = row.Int(res.Map("weight_id"))
			s.Charset = row.Str(res.Map("charset"))
			s.Active = row.Int(res.Map("active"))
			s.Onload, err = s.OnloadStatus()
			if err != nil {
				ErrorLog(fmt.Sprintf("Err: getAliveSchd get onload status fail, id:%d, %s", s.Id, err))
			}
			s.Hosts, err = s.getMasterHost()
			if err != nil {
				ErrorLog(fmt.Sprintf("Err: getAliveSched getMasterinfo fail, id:%d, %s", s.Id, err))
				//continue
			}
			//fmt.Printf("getAliveSchd: %s\n", s)
			s.Host = strings.Join(s.Hosts, ",")
			schds = append(schds, s)
		}
	}
	//ErrorLog(schds)

	return schds, nil
}

// cron job
func (j StartJob) Run() {
	var err error

	// do --replace
	s := j.schd
	err = s.Run(1)
	if err != nil {
		ErrorLog(fmt.Sprintf("Job Dry-run failed: %s\n", err))
	} else {
		err = s.Run(2)
		if err != nil {
			ErrorLog(fmt.Sprintf("Job Run failed: %s\n", err))
		}
		/*
			if s.TargetType == 1 {
				err = s.PushRemote()
				if err != nil {
					ErrorLog(fmt.Sprintf("Job Run failed: %s\n", err))
				}
			}
		*/
	}
}

// mangement cron task
func MgtCron() {
	//var err error

	co := cron.New()
	co.Start()
	defer func() {
		if err := recover(); err != nil {
			// error handling
			co.Stop()
			ErrorLog(fmt.Sprintf("Defer MgtCron Error: %s\n", err))
			return
		}
	}()

	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	CheckErr(err)
	for {
		// check added cron list
		actuals := make(map[string]string, 0)
		for _, entry := range co.Entries() {
			schdid := fmt.Sprintf("%d", entry.Job.(StartJob).schd.Id)
			actuals[schdid] = entry.Job.(StartJob).schd.Name
		}
		//ErrorLog(fmt.Sprintf("before: %v", actuals))

		// remove dead job
		jobsDead, _ := getDeadSchd()
		//ErrorLog(fmt.Sprintf("dead: %v", jobsDead))
		if jobsDead != nil {
			for schdid, schdname := range jobsDead {
				_, _, err := dbadmin.Query(fmt.Sprintf("DELETE FROM crontab WHERE scheduled_id= %s", schdid))
				if err != nil {
					ErrorLog(fmt.Sprintf("Remove crontab error: %s\n", err))
				}
				if _, ok := actuals[schdid]; ok == true {
					co.RemoveJob(schdid)
					ErrorLog(fmt.Sprintf("Remove Job: id %s, name %s \n", schdid, schdname))
				}
			}
		}

		// start add cron
		jobsAlived, err := getAliveSchd()
		if err != nil {
			ErrorLog(fmt.Sprintf("getAliveSchd Error: %s\n", err))
		}
		for _, s := range jobsAlived {
			schdid := fmt.Sprintf("%d", s.Id)
			if _, ok := actuals[schdid]; ok != true {
				//ErrorLog(fmt.Sprintf("Add job:%v \n", s))
				co.AddJob(s.Cron, StartJob{s, schdid}, schdid)
				_, _, err := dbadmin.Query(fmt.Sprintf("REPLACE INTO crontab (scheduled_id, scheduled_name, cron) VALUES(%d, '%s', '%s')", s.Id, s.Name, s.Cron))
				if err != nil {
					ErrorLog(fmt.Sprintf("Add crontab error: %s\n", err))
				}
			}
		}

		//		actuals = nil
		actuals = make(map[string]string, 0)
		for _, entry := range co.Entries() {
			jobid := fmt.Sprintf("%d", entry.Job.(StartJob).schd.Id)
			actuals[jobid] = entry.Job.(StartJob).schd.Name
		}
		//ErrorLog(fmt.Sprintf("after: %v\n", actuals))

		// time.Sleep(60 * time.Second)
		time.Sleep(time.Duration(API.CronCheckSec) * time.Second)
	}
}
