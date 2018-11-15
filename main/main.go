package main

import (
	"net/http"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	_"github.com/dwg255/landlord/router"
)

func main()  {
	err := initConf()
	if err != nil {
		logs.Error("init conf err:%v",err)
		return
	}

	err = initSec()
	if err != nil {
		logs.Error("init sec err:%v",err)
		return
	}

	var addr = flag.String("addr", fmt.Sprintf(":%d",gameConf.HttpPort), "http service address")
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		logs.Error("ListenAndServe: err:%v", err)
	}
}