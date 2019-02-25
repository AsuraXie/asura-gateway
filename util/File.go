package util

import (
	"fmt"
	"io/ioutil"
)

func ParseConfFile(path string){
	content,err := ioutil.ReadFile(path)
	if err !=nil{
		panic("config file read error :" + err.Error())
	}
	var temp []byte
	var result []string
	for i:=0;i<len(content);i++{
 		if content[i] == '\n'{
			result = append(result,string(temp))
			temp = []byte{}
		}else{
			temp = append(temp, content[i])
		}
	}
	fmt.Println(result)
}