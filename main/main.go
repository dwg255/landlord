package main // import "landlord"

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	_ "landlord/router"
	"net/http"
	"os"
	"strconv"
)

func main() {
	err := initConf()
	if err != nil {
		logs.Error("init conf err:%v", err)
		return
	}
	defer func() {
		if gameConf.Db != nil {
			err = gameConf.Db.Close()
			if err != nil {
				logs.Error("main close sqllite db err :%v", err)
			}
		}
	}()
	err = initSec()
	if err != nil {
		logs.Error("init sec err:%v", err)
		return
	}

	var addr = flag.String("addr", fmt.Sprintf(":%d", gameConf.HttpPort), "http service address")
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		logs.Error("ListenAndServe: err:%v", err)
	}
}

func init() {  	//生成pid文件，保存pid
	pidFileName := "pid"
	fileInfo, err := os.Stat(pidFileName)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(pidFileName, os.ModePerm)
			fileInfo, _ = os.Stat(pidFileName)
		}
	}
	if fileInfo.IsDir() {
		pid := os.Getpid()
		pidFile, err := os.OpenFile(pidFileName+"/landlord.pid", os.O_RDWR|os.O_CREATE, 0766)
		if err != nil {
			logs.Error("open pidFile [%s] error :%v", pidFileName, err)
			return
		}
		err = pidFile.Truncate(0)  //清空数据

		_, err = io.WriteString(pidFile, strconv.Itoa(pid))
		if err != nil {
			logs.Error("write pid error :%v", err)
		}

		err = pidFile.Close()
		if err != nil {
			logs.Error("close pid file err: %v", err)
		}
	} else {
		logs.Error("pidFile [%s] is exists and not dir", pidFileName)
	}
}
