package utils

import (
	//go官方的sql包：用于执行SQL语句
	"database/sql"
	//第三方SQL数据库驱动
	_ "github.com/go-sql-driver/mysql"
)

/**
 * Created by Chris on 2021/7/4.
 */

var (
	Db *sql.DB
	err error
)

func init() {
	//root:123456@tcp(192.168.154.19:5000)/go_db ： 用户名、密码、协议、IP：端口、连接的数据库
	Db, err = sql.Open("mysql", "root:123456@tcp(192.168.154.19:5000)/go_db")
	if err != nil{
		panic(err.Error())
	}
}

