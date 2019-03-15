package main

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"landlord/common"
	"os"
)

var (
	gameConf = &common.GameConfInfo
)

func initConf() (err error) {
	environment := os.Getenv("ENV")
	if environment != "dev" && environment != "testing" && environment != "product" {
		environment = "product"
	}
	logs.Info("the running environment is : %s", environment)
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		logs.Error("new conf failed ,err : %v", err)
		return
	}

	environment += "::"
	gameConf.HttpPort, err = conf.Int(environment + "http_port")
	if err != nil {
		logs.Error("init http_port failed,err: %v", err)
		return
	}

	logs.Debug("read conf succ , http port : %v", gameConf.HttpPort)

	//todo 日志配置
	gameConf.LogPath = conf.String(environment + "log_path")
	if len(gameConf.LogPath) == 0 {
		gameConf.LogPath = "./logs/game.log"
	}

	logs.Debug("read conf succ , LogPath :  %v", gameConf.LogPath)
	gameConf.LogLevel = conf.String(environment + "log_level")
	if len(gameConf.LogLevel) == 0 {
		gameConf.LogLevel = "debug"
	}
	logs.Debug("read conf succ , LogLevel :  %v", gameConf.LogLevel)

	//todo sqlite配置
	gameConf.DbPath = conf.String(environment + "db_path")
	if len(gameConf.DbPath) == 0 {
		gameConf.DbPath = "./db/landlord.db"
	}
	logs.Debug("read conf succ , DbPath :  %v", gameConf.DbPath)
	return
}
