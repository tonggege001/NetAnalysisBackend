package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(data *map[string]interface{}) []byte{
	m := make(map[string]interface{})
	m["code"]= 10
	if data == nil{
		ret,_ :=json.Marshal(m)
		return ret
	}

	ret,err := json.Marshal(*data)
	if err!= nil{
		log.Printf("JSON marshall err, data=%v, err=%v",*data, err)
	}
	return ret
}

func SendJson(data * map[string]interface{}, w http.ResponseWriter){
	bytes :=JSON(data)
	_, err := w.Write(bytes)
	if err != nil{
		log.Printf("SendJson error, err=%v",err)
	}
	return
}
