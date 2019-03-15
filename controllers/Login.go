package controllers

import (
	"net/http"
	"github.com/astaxie/beego/logs"
	"fmt"
	"crypto/md5"
	"landlord/common"
)

func Login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r:= recover(); r != nil {
			logs.Error("user request Login - Login panic:%v ",r)
		}
	}()
	var ret = []byte{'1'}
	defer func() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_,err := w.Write(ret)
		if err != nil {
			logs.Error("")
		}
	}()
	email := r.FormValue("email")
	if len(email) == 0 {
		email = r.PostFormValue("email")
		if email == "" {
			logs.Error("user request Login - err: email is empty")
			return
		}
	}
	password := r.FormValue("password")
	if len(password) == 0 {
		password = r.PostFormValue("password")
		if password == "" {
			logs.Error("user request Login - err: password is empty")
			return
		}
	}
	md5Password := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	var account = common.Account{}
	//err := common.GameConfInfo.MysqlConf.Pool.Get(&account, "select * from account where username=? and password", email,md5Password)
	row := common.GameConfInfo.Db.QueryRow("SELECT * FROM `account` WHERE username=? AND password=?",email,md5Password)
	if row != nil {
		err := row.Scan(&account.Id,&account.Email,&account.Username,&account.Password,&account.Coin,&account.CreatedDate,&account.UpdateDate)
		if err != nil {
			cookie := http.Cookie{Name: "userid", Value: string(account.Id), Path: "/", MaxAge: 86400}
			http.SetCookie(w, &cookie)
			cookie = http.Cookie{Name: "username", Value: account.Username, Path: "/", MaxAge: 86400}
			http.SetCookie(w, &cookie)
		}
	}
}
