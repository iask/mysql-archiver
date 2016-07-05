package utils

import (
	"io/ioutil"
	"os"
)

// exception process
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

//
func CheckFileExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

//
func ReadFile(f string) (string, error) {
	var fd []byte

	fh, err := os.Open(f)
	if err != nil {
		return string(fd), err
	}
	defer fh.Close()
	fd, err = ioutil.ReadAll(fh)
	if err != nil {
		return string(fd), err
	}

	return string(fd), nil
}
