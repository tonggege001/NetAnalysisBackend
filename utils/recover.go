package utils

import "log"

func RecoverResolve(){
	if err:=recover();err!=nil{
		log.Println("A panic happened!")
		log.Println(err) // 这里的err其实就是panic传入的内容，55
	}
}

func Recover2(string2 string)  {
	if err:=recover();err!=nil{
		log.Printf(string2+" A panic happened!")
		log.Println(err) // 这里的err其实就是panic传入的内容，55
	}
}