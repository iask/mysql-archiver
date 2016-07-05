package web

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
