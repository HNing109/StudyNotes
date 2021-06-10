package main

import (
	"factory_struct"
	"fmt"
)

func main() {
	object := factory_struct.NewFactoryStruct("chris", 18, "shanghai")
	fmt.Println(*object)
	object.SetName("FYJ")
	fmt.Println(object.GetName())
}

