package db

import (
	"fmt"
	"utils"
)

/**
 * Created by Chris on 2021/7/4.
 */

type User struct{
	ID int
	Name string
	Age int
}

/**
使用预编译添加
 */
func(u *User) AddUser(user User) error{
	sql := "insert into t_user_info (id, name, age) values(?, ?, ?)"
	//预编译
	inStmt, err := utils.Db.Prepare(sql)
	if err != nil{
		fmt.Println("error of prepare: ", err)
		return err
	}
	//执行
	_, errExe := inStmt.Exec(user.ID, user.Name, user.Age)
	if errExe != nil{
		fmt.Println("error of executive", errExe)
		return errExe
	}
	return nil
}

/**
直接添加
 */
func(u *User) AddUser2(user User) error{
	sql := "insert into t_user_info (id, name, age) values(?, ?, ?)"
	_, errExe := utils.Db.Exec(sql, user.ID, user.Name, user.Age)
	if errExe != nil{
		fmt.Println("error of executive: ", errExe)
		return errExe
	}
	return nil
}

