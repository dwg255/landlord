package common

import (
	"github.com/jmoiron/sqlx"
)

var GameConfInfo GameConf

type GameConf struct {
	HttpPort    int
	LogPath     string
	LogLevel    string
	AppSecret   string
	MysqlConf   MysqlConf
}

type MysqlConf struct {
	MysqlAddr     string
	MysqlUser     string
	MysqlPassword string
	MysqlDatabase string
	Pool          *sqlx.DB
}
