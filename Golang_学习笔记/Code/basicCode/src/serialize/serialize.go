package serialize

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Address struct{
	Province string
	City string
	Town string
	Road string

}

type IDCard struct{
	Name      string
	Age       int
	Addresses []*Address
}

type Serialize struct{

}

/**
注意：转换成JSON格式的数据，其对应的struct的属性名应为public（即：首字母大写）
否则，无法完成数据转换
 */
func (s *Serialize) Encoded(){
	chris_fj := &Address{"FuJian","quanzhou","jinjiang","189"}


	chris_sh := &Address{"shanghai","shanghai","xujiahui","130"}

	idInfo := IDCard{"chris", 18,[]*Address{chris_fj, chris_sh}}
	//进行数据转换
	js, _ := json.Marshal(idInfo)
	fmt.Printf("JSON format: %s", js)

	//将数据写入文件
	file, _ := os.OpenFile("./data/idCard.json", os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	enc := json.NewEncoder(file)
	err := enc.Encode(idInfo)
	if err != nil {
		log.Println("Error in encoding json")
	}
}
