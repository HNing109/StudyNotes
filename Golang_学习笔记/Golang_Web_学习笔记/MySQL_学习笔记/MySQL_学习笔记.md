# 1、









# 示例：Go连接MySQL数据库

- 初始化函数

  ```go
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
  ```

  

- 模型

  ```go
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
  ```

  

- 工具包

  ```go
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
  ```

  

- main函数

  ```go
  package main
  
  import (
  	"db"
  )
  
  /**
   * Created by Chris on 2021/7/4.
   */
  
  func main() {
  	var user01 = db.User{
  		ID:   0,
  		Name: "chris",
  		Age:  19,
  	}
  	var user02 = db.User{
  		ID:   1,
  		Name: "fyj",
  		Age:  15,
  	}
      //insert 数据
  	user01.AddUser(user01)
  	user02.AddUser2(user02)
  }
  ```







