package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	. "archiver/utils"

	"github.com/msbranco/goconfig"
)

var (
	MAX_PAGE   = 50
	MAX_OFFSET = 50

	API_CONFIG *AppConf

	API        *apiInfo
	LOG        *logInfo
	XBM        *xbmInfo
	DBMAP      *dbmapInfo
	ADMIN_DB   *DbInstance
	BASIC_AUTH *authInfo
)

// api http args
type apiInfo struct {
	Port         string
	Uri          string
	TimeOut      int64
	DataPath     string
	CronCheckSec int64
	DbSuperUser  string
	DbSuperPass  string
}

type logInfo struct {
	LogLevel  string
	AccessLog string
	ErrorLog  string
	LogPath   string
}

type authInfo struct {
	State bool
	User  string
	Pass  string
}

type xbmInfo struct {
	BackupUrl string
	CheckUrl  string
}

// dbmap info
type dbmapInfo struct {
	Url  string
	User string
	Pass string
}

type AppConf struct {
	Api   apiInfo
	Log   logInfo
	Auth  authInfo
	DbMap dbmapInfo
	Xbm   xbmInfo
	DbIns DbInstance
}

// read conf file
func NewConfig() *AppConf {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Defer getConf Error: %s\n", err)
			os.Exit(-1)
		}
	}()

	apiConfig := GetConfigName()
	//fmt.Println(*apiConfig)
	c, err := goconfig.ReadConfigFile(*apiConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	port, err := c.GetInt64("api", "http_port")
	CheckErr(err)
	Port := strconv.FormatInt(int64(port), 10)
	timeout, err := c.GetInt64("api", "timeout")
	CheckErr(err)
	uri, err := c.GetString("api", "uri_path")
	CheckErr(err)
	dataPath, err := c.GetString("default", "data_path")
	CheckErr(err)
	cronCheckSec, err := c.GetInt64("default", "cron_check_sec")
	CheckErr(err)
	dbSuperUser, err := c.GetString("default", "db_super_user")
	CheckErr(err)
	dbSuperPass, err := c.GetString("default", "db_super_pass")
	CheckErr(err)
	api := &apiInfo{Port, uri, timeout, dataPath, cronCheckSec, dbSuperUser, dbSuperPass}

	logLevel, err := c.GetString("default", "log_level")
	CheckErr(err)
	accessLog, err := c.GetString("default", "access_log")
	CheckErr(err)
	errorLog, err := c.GetString("default", "error_log")
	CheckErr(err)
	logPath, err := c.GetString("default", "log_path")
	CheckErr(err)
	log := &logInfo{logLevel, accessLog, errorLog, logPath}

	authState, err := c.GetBool("http_auth", "auth_state")
	CheckErr(err)
	authBasicUser, err := c.GetString("http_auth", "auth_basic_user")
	CheckErr(err)
	authBasicPass, err := c.GetString("http_auth", "auth_basic_pass")
	CheckErr(err)
	auth := &authInfo{authState, authBasicUser, authBasicPass}

	dbHost, err := c.GetString("database", "dbhost")
	CheckErr(err)
	dbPort, err := c.GetString("database", "dbport")
	CheckErr(err)
	dbUser, err := c.GetString("database", "dbuser")
	CheckErr(err)
	dbPass, err := c.GetString("database", "dbpass")
	CheckErr(err)
	dbName, err := c.GetString("database", "dbname")
	CheckErr(err)
	dbins := &DbInstance{dbHost, dbPort, dbUser, dbPass, dbName}

	dbmapUrl, err := c.GetString("dbmap", "url")
	CheckErr(err)
	dbmapUser, err := c.GetString("dbmap", "dbmap_user")
	CheckErr(err)
	dbmapPass, err := c.GetString("dbmap", "dbmap_pass")
	CheckErr(err)
	dbmap := &dbmapInfo{dbmapUrl, dbmapUser, dbmapPass}

	xbmBackupUrl, err := c.GetString("xbm_backup", "url_backup")
	CheckErr(err)
	xbmCheckUrl, err := c.GetString("xbm_backup", "url_check")
	CheckErr(err)
	xbm := &xbmInfo{xbmBackupUrl, xbmCheckUrl}

	return &AppConf{*api, *log, *auth, *dbmap, *xbm, *dbins}
}

// log process
func AccessLog(r *http.Request, status int) {
	AccessLogger, err := CreateLogger(LOG.LogPath + "/" + LOG.AccessLog)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	LogFormat := fmt.Sprintf("%s \"%s %s %s\" %d %d \"%s\" \"%v\"",
		r.RemoteAddr, r.Method, r.URL.Path, r.Proto, status, r.ContentLength, r.UserAgent(), r.Referer())
	AccessLogger.Printf(LogFormat)
}

func ErrorLog(v interface{}) {
	DebugLogger, err := CreateLogger(LOG.LogPath + "/" + LOG.ErrorLog)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	DebugLogger.Println(v)
}

func init() {
	API_CONFIG = NewConfig()

	API = &API_CONFIG.Api
	LOG = &API_CONFIG.Log
	XBM = &API_CONFIG.Xbm
	DBMAP = &API_CONFIG.DbMap
	ADMIN_DB = &API_CONFIG.DbIns
	BASIC_AUTH = &API_CONFIG.Auth
}
