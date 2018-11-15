package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"fmt"
	"github.com/dwg255/landlord/common"
)

var (
	gameConf = &common.GameConfInfo
)

func initConf() (err error) {
	gameConf.HttpPort,err = beego.AppConfig.Int("http_port")
	if err != nil {
		logs.Error("init http_port failed,err:%v",err)
		return
	}
	logs.Debug("read conf succ,http port %v", gameConf.HttpPort)

	//todo mysql相关配置
	gameConf.MysqlConf.MysqlAddr = beego.AppConfig.String("mysql_addr")
	if len(gameConf.MysqlConf.MysqlAddr) == 0 {
		err = fmt.Errorf("init config failed, mysql_addr [%s]", gameConf.MysqlConf.MysqlAddr)
		return
	}
	gameConf.MysqlConf.MysqlUser = beego.AppConfig.String("mysql_user")
	if len(gameConf.MysqlConf.MysqlUser) == 0 {
		err = fmt.Errorf("init config failed, mysql_user [%s]", gameConf.MysqlConf.MysqlUser)
		return
	}
	gameConf.MysqlConf.MysqlPassword = beego.AppConfig.String("mysql_password")
	if len(gameConf.MysqlConf.MysqlPassword) == 0 {
		err = fmt.Errorf("init config failed, mysql_password [%s]", gameConf.MysqlConf.MysqlPassword)
		return
	}
	gameConf.MysqlConf.MysqlDatabase = beego.AppConfig.String("mysql_db")
	if len(gameConf.MysqlConf.MysqlDatabase) == 0 {
		err = fmt.Errorf("init config failed, mysql_password [%s]", gameConf.MysqlConf.MysqlDatabase)
		return
	}

	//todo 密钥
	gameConf.AppSecret = beego.AppConfig.String("app_secret")
	if len(gameConf.AppSecret) == 0 {
		err = fmt.Errorf("init config failed, app_secret [%s]", gameConf.AppSecret)
		return
	}

	//todo 日志配置
	gameConf.LogPath = beego.AppConfig.String("log_path")
	gameConf.LogLevel = beego.AppConfig.String("log_level")

	return
}
