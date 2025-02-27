package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/doz-8108/api-gateway/config"
	"github.com/doz-8108/api-gateway/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func CreateMySqlConnection(envVars config.EnvVars, utils utils.Utils) *sqlx.DB {
	db := sqlx.MustConnect("mysql", envVars.MYSQL_DSN)

	err := db.Ping()
	if err != nil {
		panic(err)
	} else {
		utils.Logger.Info("MySQL connected")
	}

	file, err := os.ReadFile("./init.sql")
	if err != nil {
		utils.Logger.Error(err)
	}

	tx := db.MustBegin()
	for i, stm := range strings.Split(string(file), ";\n") {
		fmt.Println(i, stm)
		tx.MustExec(stm)
	}
	tx.Commit()
	return db
}
