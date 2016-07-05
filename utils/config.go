package utils

import (
	"flag"
)

var CONFIG *string

// for get global config file name
func GetConfigName() *string {
	return CONFIG
}

func init() {
	CONFIG = flag.String("conf", "conf/app.conf", " config file path")
	flag.Parse()
}
