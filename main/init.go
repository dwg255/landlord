package main

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
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

	configStr,err := json.Marshal(config)
	if err != nil {
		fmt.Println("marsha1 faild,err",err)
		return
	}
	logs.SetLogger(logs.AdapterFile,string(configStr))
	return
}

func initMysql() (err error) {
	conf := gameConf.MysqlConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",conf.MysqlUser,conf.MysqlPassword,conf.MysqlAddr,conf.MysqlDatabase)
	logs.Debug(dsn)
	database, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return
	}

	gameConf.MysqlConf.Pool = database
	return
}

func initSec() (err error) {
	err = initLogger()
	if err != nil {
		logs.Error("init logger failed,err:%v",err)
		return
	}

	err = initMysql()
	if err != nil {
		logs.Error("init mysql failed,err :%v",err)
		return
	}
	logs.Info("init sec succ")
	return
}

