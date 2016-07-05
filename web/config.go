package web

import (
	//	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	. "archiver/utils"

	"github.com/msbranco/goconfig"
)

var (
	WEB_CONFIG *AppConf

	WEB        *webInfo
	LOG        *logInfo
	BASIC_AUTH *authInfo
)

// http args
type webInfo struct {
	Port         string
	Uri          string
	TimeOut      int64
	ApiUrl       string
	XboxUrl      string
	TemplatePath string
}

// log args
type logInfo struct {
	LogLevel  string
	AccessLog string
	ErrorLog  string
	LogPath   string
}

// basic auth args
type authInfo struct {
	State bool
	User  string
	Pass  string
}

type AppConf struct {
	Web  webInfo
	Log  logInfo
	Auth authInfo
}

// read conf file
func NewConfig() *AppConf {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Web getConf Error: %s\n", err)
			os.Exit(-1)
		}
	}()

	webConfig := GetConfigName()
	//fmt.Println(*webConfig)
	c, err := goconfig.ReadConfigFile(*webConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	port, err := c.GetInt64("web", "http_port")
	CheckErr(err)
	Port := strconv.FormatInt(int64(port), 10)
	Uri, err := c.GetString("web", "uri_path")
	CheckErr(err)
	timeout, err := c.GetInt64("web", "timeout")
	CheckErr(err)
	ApiUrl, err := c.GetString("web", "api_url")
	CheckErr(err)
	XboxUrl, err := c.GetString("web", "xbox_url")
	CheckErr(err)
	templatePath, err := c.GetString("web", "template_path")
	CheckErr(err)

	web := &webInfo{Port, Uri, timeout, ApiUrl, XboxUrl, templatePath}

	logLevel, err := c.GetString("default", "log_level")
	CheckErr(err)
	logPath, err := c.GetString("default", "log_path")
	CheckErr(err)
	accessLog, err := c.GetString("web", "access_log")
	CheckErr(err)
	errorLog, err := c.GetString("web", "error_log")
	CheckErr(err)

	log := &logInfo{logLevel, accessLog, errorLog, logPath}

	authState, err := c.GetBool("http_auth", "auth_state")
	CheckErr(err)
	authBasicUser, err := c.GetString("http_auth", "auth_basic_user")
	CheckErr(err)
	authBasicPass, err := c.GetString("http_auth", "auth_basic_pass")
	CheckErr(err)

	auth := &authInfo{authState, authBasicUser, authBasicPass}

	return &AppConf{*web, *log, *auth}
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
	WEB_CONFIG = NewConfig()

	WEB = &WEB_CONFIG.Web
	LOG = &WEB_CONFIG.Log
	BASIC_AUTH = &WEB_CONFIG.Auth
}
