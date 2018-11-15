package controllers

import (
	"net/http"
	"html/template"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

func Index(w http.ResponseWriter, r *http.Request)  {
	ret := make(map[string]interface{})
	t, err :=template.ParseFiles("templates/poker.html")
	if err != nil {
		logs.Error("can't find template file %s","templates/poker.html")
		return
	}
	userId := ""
	cookie, err := r.Cookie("userid")
	if err != nil {
		logs.Error("get cookie err: %v",err)
	} else {
		userId = cookie.Value
	}
	username := ""
	cookie, err = r.Cookie("username")
	if err != nil {
		logs.Error("get cookie err: %v",err)
	} else {
		username = cookie.Value
	}
	logs.Debug("userId [%v]",userId)
	logs.Debug("username [%v]",username)
	user := make(map[string]interface{})
	user["uid"] = userId
	user["username"] = username
	res,err := json.Marshal(user)
	if err != nil {
		logs.Error("json marsha1 user %v ,err:%v",user,err)
		return
	}
	ret["user"] = string(res)
	t.Execute(w, ret)
}
