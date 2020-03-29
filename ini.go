package main

import (
	"fmt"
	"log"
	"gopkg.in/gcfg.v1"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Conf struct {
	Mysql struct {
		      Host     string
		      Port     string
		      Username string
		      Password string
		      Database string
	      }
	Other struct {
		      Savedir string
	      }
}

func readConfig() (Config Conf) {
	err := gcfg.ReadFileInto(&Config, "config.ini")
	if err != nil {
		fmt.Println("Failed to parse config file: %s", err)
	}
	return Config
}

func connectSql(DbConfig Conf) (db *sql.DB) {
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		DbConfig.Mysql.Username,
		DbConfig.Mysql.Password,
		DbConfig.Mysql.Host,
		DbConfig.Mysql.Port,
		DbConfig.Mysql.Database)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		log.Fatal("请检查数据库配置", err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}