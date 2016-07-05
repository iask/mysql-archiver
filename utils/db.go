package utils

import (
	"fmt"
	"time"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

// db info
type DbInstance struct {
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
}

// init db obj
func (d DbInstance) Connect() (mysql.Conn, error) {
	db := mysql.New("tcp", "", fmt.Sprintf("%s:%s", d.DbHost, d.DbPort), d.DbUser, d.DbPass, d.DbName)
	db.SetTimeout(2 * time.Second)
	err := db.Connect()

	return db, err
}
