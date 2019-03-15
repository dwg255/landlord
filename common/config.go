package common

import (
	"database/sql"
)

var GameConfInfo GameConf

type GameConf struct {
	HttpPort  int
	LogPath   string
	LogLevel  string
	AppSecret string
	DbPath    string
	Db        *sql.DB
}
