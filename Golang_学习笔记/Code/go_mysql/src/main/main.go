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
	user01.AddUser(user01)
	user02.AddUser2(user02)
}







 
