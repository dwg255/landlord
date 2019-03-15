package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/mattn/go-sqlite3"
)

func conversionLogLevel(logLevel string) int {
	switch logLevel {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = gameConf.LogPath
	config["level"] = conversionLogLevel(gameConf.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marsha1 faild,err", err)
		return
	}
	err = logs.SetLogger(logs.AdapterFile, string(configStr))
	if err != nil {
		logs.Error("init logger :%v",err)
	}
	return
}

func initSqlite() (err error) {
	gameConf.Db, err = sql.Open("sqlite3", gameConf.DbPath)
	if err != nil {
		logs.Error("initSqlite err : %v", err)
		return
	}
	var stmt *sql.Stmt
	stmt, err = gameConf.Db.Prepare(`CREATE TABLE IF NOT EXISTS "account" ("id" INTEGER NOT NULL,"email" text(32),"username" TEXT(16),"password" TEXT(32),"coin" integer,"created_date" TEXT(32),"updated_date" TEXT(32),PRIMARY KEY ("id"))`)
	if err != nil {
		logs.Error("initSqlite err : %v", err)
		return
	}
	_, err = stmt.Exec()
	if err != nil {
		logs.Error("create table err:", err)
		return
	}
	return
}

func initSec() (err error) {
	err = initLogger()
	if err != nil {
		logs.Error("init logger failed,err:%v", err)
		return
	}
	err = initSqlite()
	if err != nil {
		logs.Error("init sqlite failed,err:%v", err)
		return
	}
	logs.Info("init sec succ")
	return
}
