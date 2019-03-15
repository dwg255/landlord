package controllers

import (
	"net/http"
	"github.com/astaxie/beego/logs"
)

func LoginOut(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("LoginOut panic:%v ", r)
		}
	}()
	cookie := http.Cookie{Name: "user", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)

	_,err := w.Write([]byte{'1'})
	if err != nil {
		logs.Error("LoginOut err: %v",err)
	}
}
