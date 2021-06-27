package serialize

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
序列化
注意：转换成JSON格式的数据，其对应的struct的属性名应为public（即：首字母大写）
否则，无法完成数据转换
 */
func (s *Serialize) Encoded(){
	chris_fj := &Address{"FuJian","quanzhou","jinjiang","189"}


	chris_sh := &Address{"shanghai","shanghai","xujiahui","130"}

	idInfo := IDCard{"chris", 18,[]*Address{chris_fj, chris_sh}}
	//进行数据转换
	js, _ := json.Marshal(idInfo)
	fmt.Printf("JSON format: %s \n", js)

	//将数据写入文件
	file, _ := os.OpenFile("./data/idCard.json", os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	enc := json.NewEncoder(file)
	err := enc.Encode(idInfo)
	if err != nil {
		log.Println("Error in encoding json")
	}
}

/**
反序列化
 */
func(s *Serialize) Decode() interface{}{
	//读取文件中的数据
	var filePath = "./data/idCard.json"
	file, errFile := os.Open(filePath)
	if errFile != nil{
		fmt.Println("Error: open file failure!!, error", errFile)
		return nil
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)
	var datas string
	for {
		data, readErr := fileReader.ReadString('\n')
		if readErr == io.EOF{
			break
		}
		datas += data
	}
	fmt.Println("READ File: ",datas)


	//创建存放数据的对象
	var idCard IDCard
	//反序列化
	errUnmarshal := json.Unmarshal([]byte(datas), &idCard)

	if errUnmarshal != nil{
		fmt.Println("Error: Decode failure! ", errUnmarshal)
	}
	fmt.Print("DECODE JSON: ", idCard.Name, " ", idCard.Age, " ")
	for _, val := range idCard.Addresses{
		//打印指针中的数据
		fmt.Printf("%v  ", *val)
	}
	return idCard
}
