package utils

import (
	"log"
	"os"
)

func CreateLogger(filename string) (*log.Logger, error) {
	var L *log.Logger
	FH, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return L, err
	} else {
		L = log.New(FH, "", log.Ldate|log.Ltime)
	}
	return L, nil
}
