package controllers

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"github.com/dwg255/landlord/common"
	"crypto/md5"
	"fmt"
	"time"
	"encoding/json"
	"strconv"
)

func Register(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("Register panic:%v ", r)
		}
	}()
	var ret = []byte{'1'}
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(ret)
	}()
	username := r.FormValue("username")
	if len(username) == 0 {
		username = r.PostFormValue("username")
		if username == "" {
			logs.Error("register err: username is empty")
			return
		}
	}
	password := r.FormValue("password")
	if len(password) == 0 {
		password = r.PostFormValue("password")
		if password == "" {
			logs.Error("register err: password is empty")
			return
		}
	}

	var account = common.Account{}
	err := common.GameConfInfo.MysqlConf.Pool.Get(&account, "select * from account where username=?", username)
	if err == nil {
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	md5Password := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	result, err := common.GameConfInfo.MysqlConf.Pool.Exec(`insert into account(email,username,password,coin,created_date,updated_date) values(?,?,?,?,?,?) `,
		username, username, md5Password, 10000, now, now,
	)
	if err != nil {
		logs.Error("insert new account [%v] err : %v", username, err)
		return
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		logs.Error("insert new account [%v] get last insert id err : %v", username, err)
		return
	}
	ret, err = json.Marshal(map[string]interface{}{"uid": lastInsertId, "username": username})
	if err != nil {
		logs.Error("json marsha1 failed err:%v", err)
		return
	}
	cookie := http.Cookie{Name: "userid", Value: strconv.Itoa(int(lastInsertId)), Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "username", Value: username, Path: "/", MaxAge: 86400}
	http.SetCookie(w, &cookie)
}
