package utils

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
)

func SetRedis(key string, value string) error{

	if key=="" || value==""{
		return errors.New("SetRedis key or value is empty")
	}

	c, err := redis.Dial("tcp", "47.106.174.32:6379")
	if err != nil {
		return errors.New(fmt.Sprintf("SetRedis redis.Dial, err=%v",err))
	}

	_,err=c.Do("auth","dachuang123")
	if err!= nil{
		return errors.New(fmt.Sprintf("SetRedis c.Send error, err=%v",err))
	}

	reply, err := c.Do("set", key,value,"EX", 60*60)
	if err != nil{
		return errors.New(fmt.Sprintf("SetToken c.Send set key error, key=%v, value=%v, err=%v",key,value,err))
	}
	log.Printf("SetToken c.Do set reply=%v",reply)

	defer CloseRedis(c)
	return nil

}

func SetRedisByEx(key string, value string, ex int) error{
	if key=="" || value==""{
		return errors.New("SetRedis key or value is empty")
	}

	c, err := redis.Dial("tcp", "47.106.174.32:6379")
	if err != nil {
		return errors.New(fmt.Sprintf("SetRedis redis.Dial, err=%v",err))
	}

	_,err=c.Do("auth","dachuang123")
	if err!= nil{
		return errors.New(fmt.Sprintf("SetRedis c.Send error, err=%v",err))
	}

	reply, err := c.Do("set", key,value,"EX", ex)
	if err != nil{
		return errors.New(fmt.Sprintf("SetToken c.Send set key error, key=%v, value=%v, err=%v",key,value,err))
	}
	log.Printf("SetToken c.Do set reply=%v",reply)

	defer CloseRedis(c)
	return nil
}


func GetRedis(key string) (string, error){
	if key==""{
		return "",errors.New("GetRedis key or value is empty")
	}
	c, err := redis.Dial("tcp", "47.106.174.32:6379")
	if err != nil {
		return "",errors.New(fmt.Sprintf("GetRedis redis.Dial, err=%v",err))
	}

	_,err=c.Do("auth","dachuang123")
	if err!= nil{
		return "",errors.New(fmt.Sprintf("GetRedis c.Send error, err=%v",err))
	}

	value,err:=redis.String(c.Do("get",key))
	if err != nil{
		return "",errors.New(fmt.Sprintf("GetRedis c.Send get key error, key=%v, value=%v, err=%v",key,value,err))
	}
	defer CloseRedis(c)
	return value,nil
}

func DeleteRedis(key string) error{
	if key==""{
		return errors.New("DeleteRedis key or value is empty")
	}

	c, err := redis.Dial("tcp", "47.106.174.32:6379")
	if err != nil {
		return errors.New(fmt.Sprintf("DeleteRedis redis.Dial, err=%v",err))
	}

	_,err=c.Do("auth","dachuang123")
	if err!= nil{
		return errors.New(fmt.Sprintf("DeleteRedis c.Send error, err=%v",err))
	}

	_, err = c.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}



func CloseRedis(conn redis.Conn){
	err :=conn.Close()
	if err != nil{
		log.Printf("redis close err, err=%v", err)
	}
}