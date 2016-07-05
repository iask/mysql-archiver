package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	. "archiver/utils"
)

// task info
type Scheduled struct {
	Id         int      `json:"schd_id"`
	Name       string   `json:"name"`
	XboxTag    string   `json:"xboxtag"`
	Host       string   `json:"host"`
	Hosts      []string `json:"hosts"`
	Port       int      `json:"port"`
	Db         string   `json:"db"`
	Dbs        []string `json:"dbs"`
	Table      string   `json:"table"`
	Tables     []string `json:"tables"`
	TargetType int      `json:"target_type"`
	TargetName string   `json:"target_name"`
	LogPath    string   `json:"log_path"`
	DataPath   string   `json:"data_path"`
	QueryStr   string   `json:"query"`
	Cron       string   `json:"cron"`
	Deadline   string   `json:"deadline"`
	Charset    string   `json:"charset"`
	Weight     int      `json:"weight"`
	Active     int      `json:"active"`
	Onload     int      `json:"onload"`
}

//
type Instance struct {
	Ip           string            `json:"IP"`
	Port         string            `json:"Port"`
	Ismaster     bool              `json:"Ismaster"`
	Endpoint     string            `json:"Endpoint"`
	Service      string            `json:"Service"`
	Tags         map[string]string `json:"Tags"`
	Owt          string            `json:"Owt"`
	Pdl          string            `json:"Pdl"`
	Vip          string            `json:"Vip"`
	ServiceGroup string            `json:"ServiceGroup"`
	Remark       string            `json:"Remark"`
}

//
type Cluster struct {
	Master Instance   `json:"Master"`
	Slaves []Instance `json:"Slaves"`
}

// define update scheduled post data
type schdField struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	Data  string `json:"data"`
}

//
type XboxTag struct {
	Host string `json:"host"`
	Tag  string `json:"tag"`
}

//
func (s *Scheduled) AddJob(log string) (int64, error) {
	var jobId int64

	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return jobId, err
	}

	q := fmt.Sprintf("INSERT INTO jobs (scheduled_id, start_time, target_name, dbhost, stdout, stderr) VALUES (%d, now(), '%s', '%s', '%s','%s');", s.Id, s.TargetName, s.Host, fmt.Sprintf("%s.log", log), fmt.Sprintf("%s.err", log))
	ErrorLog(fmt.Sprintf("%s\n", q))
	_, res, err := dbadmin.Query(q)
	if err != nil {
		return jobId, err
	} else {
		jobId = int64(res.InsertId())
	}

	return jobId, nil
}

//
func (s *Scheduled) AddRunningJobs(pid int) error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}

	q := fmt.Sprintf("INSERT INTO running_jobs (scheduled_id, pid, target_name, dbhost) VALUES (%d, %d, '%s', '%s');", s.Id, pid, s.TargetName, s.Host)
	ErrorLog(fmt.Sprintf("%s\n", q))
	_, _, err = dbadmin.Query(q)
	if err != nil {
		return err
	}

	return nil
}

/*
 * get master info
 * {
 *	"host":"",
 *	"tag":"",
 * }
 */
func (s *Scheduled) getMasterHost() ([]string, error) {
	var master []string
	var t XboxTag
	var host string

	//ErrorLog(fmt.Sprintf("schd: %v\n", s))
	t.Tag = s.XboxTag

	dbmapReq, err := json.Marshal(t)
	if err != nil {
		return master, err
	}

	retJson, err := HttpPost(DBMAP.Url, dbmapReq, DBMAP.User, DBMAP.Pass)
	if err != nil {
		return master, err
	}

	//ErrorLog(fmt.Sprintf("DbMap: %s\n", retJson))
	ret := DbMapOutput{}
	err = json.Unmarshal(retJson, &ret)
	if err != nil {
		return master, err
	}
	if ret.Err != "OK" {
		return master, fmt.Errorf(ret.Err)
	}

	for _, clu := range ret.Clu {
		//fmt.Printf("---> master: %s\n", clu.Master)
		if len(clu.Master.Vip) > 0 {
			host = clu.Master.Vip
		} else {
			host = clu.Master.Ip
		}

		dbIns := &DbInstance{host, strconv.Itoa(s.Port), API.DbSuperUser, API.DbSuperPass, ""}
		jobdb, err := dbIns.Connect()
		if err != nil {
			return master, err
		}
		defer jobdb.Close()
		rows, _, err := jobdb.Query("SHOW SLAVE HOSTS;")
		if err != nil {
			return master, err
		}
		if len(rows) == 0 {
			return master, fmt.Errorf("Host: %s is not master.", host)
		}
		master = append(master, host)
	}

	return master, nil
}

/*
 * split db list
 */
func (s *Scheduled) getDbs() []string {
	var ds []string

	ds = strings.Split(s.Db, ",")
	for k, v := range ds {
		ds[k] = strings.Trim(v, " ")
	}

	return ds
}

/*
 * split table list
 */
func (s *Scheduled) getTables() []string {
	var ts []string

	ts = strings.Split(s.Table, ",")
	for k, v := range ts {
		ts[k] = strings.Trim(v, " ")
	}

	return ts
}

// add Scheduled
func (s *Scheduled) Add() error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	s.QueryStr = strings.Replace(strings.Replace(s.QueryStr, "\"", "''", -1), "'", "''", -1)
	q := fmt.Sprintf("INSERT INTO scheduled (name, xboxtag, port, db, tbl, querystr, cron, target_type, deadline, charset, weight_id) VALUES('%s', '%s', %d, '%s', '%s', '%s', '%s', %d, '%s', '%s', %d )",
		s.Name, s.XboxTag, s.Port, s.Db, s.Table, s.QueryStr, s.Cron, s.TargetType, s.Deadline, s.Charset, s.Weight)
	ErrorLog(q)
	_, res, err := dbadmin.Query(q)
	if err != nil {
		return err
	}
	s.Id = int(res.InsertId())

	return nil
}

// get Scheduled info by Id
func (s *Scheduled) Get() error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}
	q := fmt.Sprintf("SELECT * FROM scheduled WHERE scheduled_id = %d;", s.Id)
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return err
	}
	for _, row := range rows {
		s.Id = row.Int(res.Map("scheduled_id"))
		s.Name = row.Str(res.Map("name"))
		s.XboxTag = row.Str(res.Map("xboxtag"))
		s.Port = row.Int(res.Map("port"))
		s.Db = row.Str(res.Map("db"))
		s.Table = row.Str(res.Map("tbl"))
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
			ErrorLog(fmt.Sprintf("Error: s.Get() get onload status fail, schdid:%d, %s", s.Id, err))
		}
		s.Hosts, _ = s.getMasterHost()
		s.Host = strings.Join(s.Hosts, ",")
		s.Dbs = s.getDbs()
		s.Tables = s.getTables()
	}

	return nil
}

// get Scheduled info by Id
func (s *Scheduled) OnloadStatus() (int, error) {
	var status int

	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return status, err
	}
	q := fmt.Sprintf("SELECT * FROM crontab WHERE scheduled_id = %d;", s.Id)
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return status, err
	}
	for _, row := range rows {
		status = row.Int(res.Map("scheduled_id"))
	}

	return status, nil
}

// Del Scheduled info by Id
func (s *Scheduled) Del() error {
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE FROM scheduled WHERE scheduled_id = %d;", s.Id)
	_, _, err = dbadmin.Query(q)
	if err != nil {
		return err
	}

	return nil
}

//
func (s *Scheduled) checkRunning() (int, error) {
	var id int

	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return id, err
	}

	q := fmt.Sprintf("SELECT * FROM running_jobs WHERE scheduled_id = %d;", s.Id)
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return id, err
	}
	for _, row := range rows {
		id = row.Int(res.Map("scheduled_id"))
	}

	return id, nil
}

func (s *Scheduled) checkDeadline() error {
	tNow := time.Now()
	tz, _ := time.LoadLocation("Local")
	deadline, err := time.ParseInLocation("2006-01-02 15:04:05", s.Deadline, tz)
	if err != nil {
		return err
	}
	tdiff := tNow.Sub(deadline)
	if tdiff.Seconds() >= 0 {
		return fmt.Errorf(fmt.Sprintf("Please check deadline by task_id: %d", s.Id))
	}

	return nil
}

// do job
func (s *Scheduled) Run(runMode int) error {
	var err error
	var j Job
	var q, ptArgv, dstTable, _targetName string

	actions := map[int]string{
		1: "--dry-run",
		2: "--replace",
		3: "--purge",
	}

	for _, Host := range s.Hosts {
		if len(Host) <= 0 {
			return fmt.Errorf("Please get master info at first.")
		}

		if runMode > 1 {
			err := s.checkDeadline()
			if err != nil {
				return err
			}
			sid, err := s.checkRunning()
			if err != nil {
				return err
			}
			if s.Id == sid {
				return fmt.Errorf(fmt.Sprintf("Warning: Previous task is running. scheduled_id is %d", s.Id))
			}
		}

		// mkdir targetPath
		tNow := time.Now()
		timeNow := tNow.Format("20060102_150405")
		basePath := fmt.Sprintf("%s/%s", API.DataPath, Host)
		targetPath := fmt.Sprintf("%s/%s", basePath, timeNow)
		if CheckFileExist(basePath) == false {
			err = os.MkdirAll(basePath, 0774)
			if err != nil {
				return err
			}
		}

		s.LogPath = fmt.Sprintf("%s/log", targetPath)
		if CheckFileExist(s.LogPath) == false && runMode != 1 {
			err = os.MkdirAll(s.LogPath, 0774)
			if err != nil {
				return err
			}
		}
		s.DataPath = fmt.Sprintf("%s/data", targetPath)
		if CheckFileExist(s.DataPath) == false && runMode != 1 {
			err = os.MkdirAll(s.DataPath, 0774)
			if err != nil {
				return err
			}
		}
		for _, srcDB := range s.Dbs {
			dbIns := &DbInstance{Host, strconv.Itoa(s.Port), API.DbSuperUser, API.DbSuperPass, srcDB}
			jobdb, err := dbIns.Connect()
			if err != nil {
				return err
			}
			defer jobdb.Close()

			var targetLog string
			for _, srcTable := range s.Tables {
				tNow = time.Now()
				timeNow = tNow.Format("20060102_150405")
				targetName := fmt.Sprintf("%s_%s", srcTable, timeNow)
				targetFile := fmt.Sprintf("%s/%s", s.DataPath, targetName)

				_targetName = ""
				switch s.TargetType {
				case 1:
					s.TargetName = fmt.Sprintf("%s.sql", targetFile)
					if runMode != 3 {
						_targetName = fmt.Sprintf("--file=%s", s.TargetName)
					}
					ptArgv = fmt.Sprintf("--source h=%s,P=%d,u=%s,p=%s,D=%s,t=%s,A=%s %s --where=\"%s\" --statistics --why-quit %s",
						Host, s.Port, API.DbSuperUser, API.DbSuperPass, s.Db, srcTable, s.Charset, _targetName, s.QueryStr, actions[runMode])
				case 2:
					if runMode != 3 {
						//Copy Table Schema
						dstTable = fmt.Sprintf("%s_%s", srcTable, timeNow)
						s.TargetName = fmt.Sprintf("%s.%s", s.Db, dstTable)
						_targetName = fmt.Sprintf("--dest h=%s,P=%d,u=%s,p=%s,D=%s,t=%s,A=%s", Host, s.Port, API.DbSuperUser, API.DbSuperPass, s.Db, dstTable, s.Charset)
						q = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s LIKE %s;", dstTable, srcTable)
						_, _, err = jobdb.Query(q)
						if err != nil {
							ErrorLog(fmt.Sprintf("Copy table schema error: %s\n", err))
							return err
						}
					}
					ptArgv = fmt.Sprintf("--source h=%s,P=%d,u=%s,p=%s,D=%s,t=%s,A=%s %s --where=\"%s\" --statistics --why-quit %s",
						Host, s.Port, API.DbSuperUser, API.DbSuperPass, s.Db, srcTable, s.Charset, _targetName, s.QueryStr, actions[runMode])
				default:
					ptArgv = "--help"
				}

				if runMode > 1 {
					targetLog = fmt.Sprintf("%s/%s", s.LogPath, targetName)
					j.Id, err = s.AddJob(targetLog)
					if err != nil {
						ErrorLog(fmt.Sprintf("Add initializing jobs record failed: %s\n", err))
						return err
					}
					err = j.Get()
					if err != nil {
						ErrorLog(fmt.Sprintf("Get jobs info failed: %s\n", err))
						return err
					}
				} else {
					targetLog = fmt.Sprintf("%s/%s/%s", API.DataPath, Host, targetName)
				}

				// run pt-archiver
				c := &Cmd{
					Command: "pt-archiver",
					Argv:    ptArgv,
					LogPath: targetLog,
				}
				ErrorLog(fmt.Sprintf("%v", c))

				err = c.Start()
				if err != nil {
					ErrorLog(fmt.Sprintf("Start pt-archiver Error: %s\n", err))
					return fmt.Errorf("Start pt-archiver Error: %s\n", err)
				}

				if runMode > 1 {
					err = s.AddRunningJobs(c.Cmd.Process.Pid)
					if err != nil {
						ErrorLog(fmt.Sprintf("Add running jobs record failed: %s\n", err))
						return err
					}
					err = j.UpdateStatus2(c.Cmd.Process.Pid)
					if err != nil {
						ErrorLog(fmt.Sprintf("Update jobs status=2 failed: %s\n", err))
						return err
					}
				}

				log, err := c.WriteLog()
				if err != nil {
					return fmt.Errorf("Write log Error: %s\n", log)
				}
				//go c.WriteRunLog()

				_, err = c.Wait()
				if err != nil {
					ErrorLog(fmt.Sprintf("Wait pt-archiver Error: %s\n", log))
					return fmt.Errorf("Wait pt-archiver Error: %s\n", log)
				}

				if runMode == 1 {
					err := c.RemoveLog()
					if err != nil {
						return fmt.Errorf("Remove log Error: %s\n", err)
					}
				}
				if runMode > 1 {
					err = j.DelRunningJobs()
					if err != nil {
						ErrorLog(fmt.Sprintf("Delete running jobs failed: %s\n", err))
						return err
					}
					err = j.UpdateStatus3()
					if err != nil {
						ErrorLog(fmt.Sprintf("Update jobs status=3 failed: %s\n", err))
						return err
					}
				}

				//clear dry-run table
				if runMode == 1 {
					if s.TargetType == 2 {
						// Remove Table Schema
						q = fmt.Sprintf("DROP TABLE IF EXISTS %s;", dstTable)
						ErrorLog(q)
						_, _, err = jobdb.Query(q)
						if err != nil {
							ErrorLog(fmt.Sprintf("Remove dry-run table error: %s\n", err))
							return err
						}
					}
				}
			}
		}

		// push sql file to xbm backup system
		if runMode == 2 && s.TargetType == 1 {
			x := XbmInfo{}

			xbmdata, err := HttpPost(XBM.BackupUrl, []byte(fmt.Sprintf("{\"Appname\":\"Archiver\", \"Master\":\"%s\"}", Host)), "", "")
			if err != nil {
				return fmt.Errorf("Get xbm info Error: %s\n", err)
			}
			ErrorLog(fmt.Sprintf("%s", xbmdata))

			err = json.Unmarshal(xbmdata, &x)
			if err != nil {
				return fmt.Errorf("Parse xbm info Error: %s\n", err)
			}
			ErrorLog(fmt.Sprintf("%s", x))

			j.UpdateXbm1(x)

			// exec tar && nc to xbm_host
			dataDir, dataName := path.Split(s.DataPath)
			err = os.Chdir(dataDir)
			if err != nil {
				return err
			}
			c := &Cmd{Command: "tar", Argv: fmt.Sprintf("-zcvf - %s | nc %s %s", dataName, x.XbmHost, x.NcPort)}
			ErrorLog(fmt.Sprintf("%v", c))
			err = c.Start()
			if err != nil {
				ErrorLog(fmt.Sprintf("Start push Error: %s\n", err))
				return fmt.Errorf("Start push Error: %s\n", err)
			}
			_, err = c.Wait()
			if err != nil {
				ErrorLog(fmt.Sprintf("Wait push Error: %s\n", err))
				return fmt.Errorf("Wait push Error: %s\n", err)
			}

			// check && update push status
			xbmdata, err = HttpPost(XBM.CheckUrl, []byte(fmt.Sprintf("{\"Task_id\":\"%s\"}", x.TaskId)), "", "")
			if err != nil {
				return fmt.Errorf("Get xbm status Error: %s\n", err)
			}
			x1 := XbmInfo{}
			err = json.Unmarshal(xbmdata, &x1)
			if err != nil {
				return fmt.Errorf("Parse xbm status Error: %s\n", err)
			}
			ErrorLog(fmt.Sprintf("%s", x1))
			x.XbmDir, x.Status = x1.XbmDir, x1.Status
			ErrorLog(fmt.Sprintf("%s", x))

			j.UpdateXbm2(x)
		}
	}

	return err
}

//
func (s *Scheduled) Update(schd Scheduled) (string, error) {
	var q string

	active := -1
	if schd.Active == 1 {
		active = 1
	} else {
		active = 0
	}
	if len(schd.Name) > 0 {
		active = 0
		schd.Name = strings.Replace(strings.Replace(schd.Name, "\"", "''", -1), "'", "''", -1)
		q = fmt.Sprintf("%s name='%s',", q, schd.Name)
	}
	if len(schd.Deadline) > 0 {
		active = 0
		schd.Deadline = strings.Replace(strings.Replace(schd.Deadline, "\"", "''", -1), "'", "''", -1)
		q = fmt.Sprintf("%s deadline='%s',", q, schd.Deadline)
	}
	if len(schd.Cron) > 0 {
		active = 0
		schd.Cron = strings.Replace(strings.Replace(schd.Cron, "\"", "''", -1), "'", "''", -1)
		q = fmt.Sprintf("%s cron='%s',", q, schd.Cron)
	}
	if len(schd.QueryStr) > 0 {
		active = 0
		schd.QueryStr = strings.Replace(strings.Replace(schd.QueryStr, "\"", "''", -1), "'", "''", -1)
		q = fmt.Sprintf("%s querystr='%s',", q, schd.QueryStr)
	}
	if len(schd.Charset) > 0 {
		active = 0
		schd.Charset = strings.Replace(strings.Replace(schd.Charset, "\"", "''", -1), "'", "''", -1)
		q = fmt.Sprintf("%s charset='%s',", q, schd.Charset)
	}
	if schd.TargetType > 0 {
		active = 0
		q = fmt.Sprintf("%s target_type='%d',", q, schd.TargetType)
	}
	if schd.Weight > 0 {
		active = 0
		q = fmt.Sprintf("%s weight_id='%d',", q, schd.Weight)
	}

	q = fmt.Sprintf("UPDATE scheduled SET %s active = %d WHERE scheduled_id = %d;", q, active, s.Id)
	ErrorLog(q)
	dbadmin, err := ADMIN_DB.Connect()
	defer dbadmin.Close()
	if err != nil {
		return "Db connect error", err
	}
	_, res, err := dbadmin.Query(q)
	if err != nil {
		return "", err
	}

	return res.Message(), nil
}

func listSchd(id int) ([]Scheduled, error) {
	var s Scheduled
	var schds []Scheduled

	dbadmin, err := ADMIN_DB.Connect()
	if err != nil {
		return schds, err
	}
	if id > 0 {
		s.Id = id
		s.Get()
		schds = append(schds, s)
		return schds, nil
	}
	q := "SELECT * FROM scheduled ORDER BY name"
	rows, res, err := dbadmin.Query(q)
	if err != nil {
		return schds, nil
	}
	for _, row := range rows {
		s = Scheduled{}
		s.Id = row.Int(res.Map("scheduled_id"))
		s.Name = row.Str(res.Map("name"))
		s.XboxTag = row.Str(res.Map("xboxtag"))
		s.Port = row.Int(res.Map("port"))
		s.Db = row.Str(res.Map("db"))
		s.Table = row.Str(res.Map("tbl"))
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
			ErrorLog(fmt.Sprintf("Error: listSchd() get onload status fail, schdid:%d, %s", s.Id, err))
		}
		s.Hosts, err = s.getMasterHost()
		if err != nil {
			ErrorLog(fmt.Sprintf("Error: ListSchd getMasterinfo fail, schdid:%d, %s", s.Id, err))
			//continue
		}
		s.Host = strings.Join(s.Hosts, ",")
		s.Dbs = s.getDbs()
		s.Tables = s.getTables()
		schds = append(schds, s)
	}

	return schds, nil
}
