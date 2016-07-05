package api

import (
	. "archiver/utils"
	"fmt"
)

// jobs info
type Job struct {
	Id          int64  `json:"job_id"`
	SchdId      int    `json:"schd_id"`
	StartTime   string `json:"start_time"`
	RunningTime int    `json:"running_time"`
	EndTime     string `json:"end_time"`
	Status      int    `json:"status"`
	Pid         int    `json:"pid"`
	Killed      int    `json:"killed"`
	TargetName  string `json:"target_name"`
	Host        string `json:"dbhost"`
	TaskId      string `json:"Task_id"`
	XbmHost     string `json:"Xbm_host"`
	XbmDir      string `json:"Xbm_dir"`
	XbmStatus   string `json:"Xbm_status"`
	StdoutLog   string `json:"stdout"`
	StderrLog   string `json:"stderr"`
}

type XbmInfo struct {
	TaskId  string `json:"Task_id"`
	AppName string `json:"Appname"`
	Name    string `json:"Name"`
	Master  string `json:"Master"`
	XbmHost string `json:"Xbm_host"`
	XbmDir  string `json:"Dir"`
	NcPort  string `json:"Ncport"`
	Status  string `json:"Status"`
}

//
func listJobs(id int, page int, offset int) ([]Job, error) {
	var j Job
	var jobs []Job

	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return jobs, err
	}
	if id > 0 {
		j.Id = int64(id)
		j.Get()
		jobs = append(jobs, j)
		return jobs, nil
	}
	q := fmt.Sprintf("SELECT * FROM jobs ORDER BY job_id DESC Limit %d, %d;", (page-1)*offset, offset)
	//ErrorLog(fmt.Sprintf("%s\n", q))
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return jobs, err
	}
	for _, row := range rows {
		j.Id = row.Int64(res.Map("job_id"))
		j.SchdId = row.Int(res.Map("scheduled_id"))
		j.StartTime = row.Str(res.Map("start_time"))
		j.EndTime = row.Str(res.Map("end_time"))
		j.RunningTime = row.Int(res.Map("running_time"))
		j.Status = row.Int(res.Map("status"))
		j.Pid = row.Int(res.Map("pid"))
		j.Killed = row.Int(res.Map("killed"))
		j.TargetName = row.Str(res.Map("target_name"))
		j.Host = row.Str(res.Map("dbhost"))
		//j.TaskId = row.Int(res.Map("backup_id"))
		j.TaskId = fmt.Sprintf("%d", row.Int(res.Map("backup_id")))
		j.XbmHost = row.Str(res.Map("backup_host"))
		j.XbmDir = row.Str(res.Map("backup_dir"))
		j.XbmStatus = row.Str(res.Map("backup_status"))
		j.StdoutLog = row.Str(res.Map("stdout"))
		j.StderrLog = row.Str(res.Map("stderr"))
		jobs = append(jobs, j)
	}

	return jobs, nil
}

//
func (j *Job) Get() error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	q := fmt.Sprintf("SELECT * FROM jobs WHERE job_id = %d;", j.Id)
	//	ErrorLog(fmt.Sprintf("%s\n", q))
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return err
	}
	for _, row := range rows {
		j.Id = row.Int64(res.Map("job_id"))
		j.SchdId = row.Int(res.Map("scheduled_id"))
		j.StartTime = row.Str(res.Map("start_time"))
		j.EndTime = row.Str(res.Map("end_time"))
		j.RunningTime = row.Int(res.Map("running_time"))
		j.Status = row.Int(res.Map("status"))
		j.Pid = row.Int(res.Map("pid"))
		j.Killed = row.Int(res.Map("killed"))
		j.TargetName = row.Str(res.Map("target_name"))
		j.Host = row.Str(res.Map("dbhost"))
		//j.TaskId = row.Int(res.Map("backup_id"))
		j.TaskId = fmt.Sprintf("%d", row.Int(res.Map("backup_id")))
		j.XbmHost = row.Str(res.Map("backup_host"))
		j.XbmDir = row.Str(res.Map("backup_dir"))
		j.XbmStatus = row.Str(res.Map("backup_status"))
		j.StdoutLog = row.Str(res.Map("stdout"))
		j.StderrLog = row.Str(res.Map("stderr"))
	}

	return nil
}

//
func (j *Job) DelRunningJobs() error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	q := fmt.Sprintf("DELETE FROM running_jobs WHERE scheduled_id = %d;", j.SchdId)
	ErrorLog(fmt.Sprintf("%s\n", q))
	_, _, err = dbadmin.Query(q)
	if err != nil {
		return err
	}

	return nil
}

//
func (j *Job) UpdateStatus2(pid int) error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	q := fmt.Sprintf("UPDATE jobs SET status=2, pid=%d WHERE job_id = %d;", pid, j.Id)
	ErrorLog(fmt.Sprintf("%s\n", q))
	_, _, err = dbadmin.Query(q)
	if err != nil {
		return err
	}

	return nil
}

//
func (j *Job) UpdateStatus3() error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	q := fmt.Sprintf("UPDATE jobs SET status=3, end_time=now(), running_time=TIMESTAMPDIFF(SECOND, start_time, now()) WHERE job_id = %d;", j.Id)
	ErrorLog(fmt.Sprintf("%s\n", q))
	_, _, err = dbadmin.Query(q)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) UpdateXbm1(x XbmInfo) error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	q := fmt.Sprintf("UPDATE jobs SET backup_id=%s, backup_name='%s', backup_host='%s' WHERE job_id=%d", x.TaskId, x.AppName, x.XbmHost, j.Id)
	ErrorLog(fmt.Sprintf("%s", q))
	_, _, err = dbadmin.Query(q)

	return nil
}

func (j *Job) UpdateXbm2(x XbmInfo) error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	q := fmt.Sprintf("UPDATE jobs SET backup_dir='%s', backup_status='%s' WHERE job_id=%d", x.XbmDir, x.Status, j.Id)
	ErrorLog(fmt.Sprintf("%s", q))
	_, _, err = dbadmin.Query(q)
	if err != nil {
		return fmt.Errorf("Update xbm status Error: %s\n", err)
	}

	return nil
}

func (j *Job) Logs() (map[string]string, error) {
	var err error
	var stdoutlog, stderrlog string
	var loginfo = make(map[string]string, 0)

	if CheckFileExist(j.StdoutLog) == false {
		stdoutlog = ""
	} else {
		stdoutlog, err = ReadFile(j.StdoutLog)
		if err != nil {
			CheckErr(fmt.Errorf(fmt.Sprintf("ReadJobLog stdout job_id %d Failed. %s", j.Id, err)))
		}
	}
	if CheckFileExist(j.StderrLog) == false {
		stderrlog = ""
	} else {
		stderrlog, err = ReadFile(j.StderrLog)
		if err != nil {
			CheckErr(fmt.Errorf(fmt.Sprintf("ReadJobLog stderr job_id %d Failed. %s", j.Id, err)))
		}
	}
	loginfo["stdout"], loginfo["stderr"] = stdoutlog, stderrlog

	return loginfo, nil
}
