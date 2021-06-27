package main

import (
	"serialize"
)
func main() {
	seri := serialize.Serialize{}
	//序列化、保存
	seri.Encoded()

	//反序列化
	seri.Decode()

}
