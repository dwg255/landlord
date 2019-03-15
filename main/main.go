package main // import "landlord"

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "landlord/router"
	"net/http"
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
