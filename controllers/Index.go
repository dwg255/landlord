package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"html/template"
	"landlord/common"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	ret := make(map[string]interface{})
	t, err := template.ParseFiles("templates/poker.html")
	if err != nil {
		logs.Error("user request Index - can't find template file %s", "templates/poker.html")
		return
	}
	userId := ""
	cookie, err := r.Cookie("userid")
	if err != nil {
		logs.Error("user request Index - get cookie err: %v", err)
	} else {
		userId = cookie.Value
	}
	username := ""
	cookie, err = r.Cookie("username")
	if err != nil {
		logs.Error("user request Index - get cookie err: %v", err)
	} else {
		username = cookie.Value
	}
	logs.Debug("user request Index - userId [%v]", userId)
	logs.Debug("user request Index - username [%v]", username)
	user := make(map[string]interface{})
	user["uid"] = userId
	user["username"] = username
	res, err := json.Marshal(user)
	if err != nil {
		logs.Error("user request Index - json marsha1 user %v ,err:%v", user, err)
		return
	}
	ret["user"] = string(res)
	ret["port"] = strconv.Itoa(common.GameConfInfo.HttpPort)
	err = t.Execute(w, ret)
	if err != nil {
		logs.Error("user request Index - template execute err : %v", err)
	}
}
