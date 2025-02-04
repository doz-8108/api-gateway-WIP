package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/doz-8108/api-gateway/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func CreateMySqlConnection(envVars config.EnvVars) *sqlx.DB {
	db := sqlx.MustConnect("mysql", envVars.MYSQL_DSN)

	err := db.Ping()
	if err != nil {
		panic(err)
	} else {
		println("MySQL connected")
	}

	file, err := os.ReadFile("./init.sql")
	if err != nil {
		println(err)
	}

	tx := db.MustBegin()
	for i, stm := range strings.Split(string(file), ";\n") {
		fmt.Println(i, stm)
		tx.MustExec(stm)
	}
	tx.Commit()
	return db
}
