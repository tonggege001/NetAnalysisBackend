package call

import (
	"NetAnalysisBackend/utils"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetFollowingUidListByUid(uid int64, accessToken string) ([]int64,error) {
	defer utils.RecoverResolve()
	apiUrl := "https://api.weibo.com/2/friendships/friends/ids.json"
	retList := make([]int64,0)

	sendMap := make(map[string]interface{})
	curCursor := 0
	for {
		sendMap["access_token"] = accessToken
		sendMap["uid"] = uid
		sendMap["cursor"] = curCursor
		byteJson,err := json.Marshal(sendMap)
		if err != nil{
			log.Printf("GetFollowingUidListByUid Marshal sendMap=%v",sendMap)
			return retList,err
		}

		body := bytes.NewReader(byteJson)
		client := &http.Client{}
		request, err := http.NewRequest("GET", apiUrl, body) //建立一个请求
		if err!= nil{
			log.Printf("GetFollowingUidListByUid http.NewRequest error body=%v, err=%v",body,err)
			return retList, err
		}

		response, err := client.Do(request) //提交
		if err!= nil {
			log.Printf("GetFollowingUidListByUid client.Do(request) error, err=%v",err)
			return retList,err
		}

		respBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("GetFollowingUidListByUid ioutil.ReadAll(response.Body) error, err=%v", err)
			return retList, err
		}

		respMap := make(map[string]interface{})
		err =json.Unmarshal(respBody, &respMap)
		if err != nil {
			log.Printf("GetFollowingUidListByUid json.Unmarshal(respBody, &respMap) error")
			return retList, err
		}

		//判断登录态
		if respMap["error"] != nil{
			//登录态过期则
			return retList, errors.New("not Login")
		}
		err = response.Body.Close()
		if err != nil{
			log.Printf("GetFollowingUidListByUid response.Body.Close() error, err=%v",err)
			return retList,err
		}
		//跳出判断
		if respMap["previous_cursor"].(int) == respMap["next_cursor"].(int){
			break
		}

		for _,ids := range respMap["ids"].([]int64){
			retList = append(retList,ids)
		}

		curCursor = respMap["next_cursor"].(int)

	}

	return retList,nil
}


func GetWeiboListByUid(uid int64, accessToken string)([]map[string]interface{},error){
	defer utils.RecoverResolve()
	apiUrl := "https://api.weibo.com/2/statuses/home_timeline.json"
	retList := make([]map[string]interface{},0)

	page := 1
	sendMap := make(map[string]interface{})
	for {
		sendMap["access_token"] = accessToken
		sendMap["count"] = 100
		sendMap["page"] = 1
		byteJson,err := json.Marshal(sendMap)
		if err != nil{
			log.Printf("GetWeiboListByUid Marshal sendMap=%v",sendMap)
			return retList,err
		}

		body := bytes.NewReader(byteJson)
		client := &http.Client{}
		request, err := http.NewRequest("GET", apiUrl, body) //建立一个请求
		if err!= nil{
			log.Printf("GetWeiboListByUid http.NewRequest error body=%v, err=%v",body,err)
			return retList, err
		}

		response, err := client.Do(request) //提交
		if err!= nil {
			log.Printf("GetWeiboListByUid client.Do(request) error, err=%v",err)
			return retList,err
		}

		respBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("GetWeiboListByUid ioutil.ReadAll(response.Body) error, err=%v", err)
			return retList, err
		}

		respMap := make(map[string]interface{})
		err =json.Unmarshal(respBody, &respMap)
		if err != nil {
			log.Printf("GetWeiboListByUid json.Unmarshal(respBody, &respMap) error")
			return retList, err
		}

		//判断登录态
		if respMap["error"] != nil{
			//登录态过期则
			return retList, errors.New("not Login")
		}

		err = response.Body.Close()
		if err != nil{
			log.Printf("GetWeiboListByUid response.Body.Close() error, err=%v",err)
			return retList,err
		}

		//跳出判断
		statuses := respMap["statuses"].([]map[string]interface{})
		if len(statuses) == 0{
			break
		}else{
			page++
		}

		for _,value := range statuses{
			weiboTime,err := time.Parse("Mon Jan 02 03:04:05 +0800 2006",value["created_at"].(string))
			if err != nil {
				weiboTime = time.Now()
			}
			weiboDetail := make(map[string]interface{})
			weiboDetail["content"] = value["text"]
			weiboDetail["time"] = weiboTime.Unix()
			retList = append(retList, weiboDetail)
		}
	}

	return retList, nil
}






