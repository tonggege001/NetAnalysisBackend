package call

import (
	"NetAnalysisBackend/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetFollowingUidListByUid(uid int64, accessToken string) ([]int64,error) {
	defer utils.RecoverResolve()
	apiUrl := "https://api.weibo.com/2/friendships/friends/ids.json"
	retList := make([]int64,0)
	//totalNum := 0
	curCursor := 0
	for {
		Url, err := url.Parse(apiUrl)
		if err != nil{
			log.Printf("GetFollowingUidListByUid url.Parse error, err=%v",err)
			return retList,err
		}

		params:=url.Values{}
		params.Set("access_token",accessToken)  //这两种都可以
		params.Set("uid",strconv.FormatInt(uid,10))
		params.Set("cursor",strconv.FormatInt(int64(curCursor),10))
		Url.RawQuery = params.Encode()				//如果参数中有中文参数,这个方法会进行URLEncode
		urlPath := Url.String()
		response, err := http.Get(urlPath)
		if err!= nil {
			log.Printf("GetFollowingUidListByUid http.Get(urlPath) error, err=%v",err)
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
			print(retList)
			log.Printf("GetFollowingUidListByUid respMap[\"error\"], respMap=%v",respMap)
			return retList, errors.New("not Login")
		}

		err = response.Body.Close()
		if err != nil{
			log.Printf("GetFollowingUidListByUid response.Body.Close() error, err=%v",err)
			return retList,err
		}

		for _,ids := range respMap["ids"].([]interface{}){
			retList = append(retList,(int64)(ids.(float64)))
		}

		//totalNum = (int)(respMap["total_number"].(float64))
		curCursor = curCursor+5

		//跳出判断
		break	//只获取一次的uid（5个）

	}

	return retList,nil
}


func GetWeiboListByUid(uid int64, accessToken string)([]map[string]interface{},error){
	defer utils.RecoverResolve()
	apiUrl := "https://api.weibo.com/2/statuses/home_timeline.json"
	retList := make([]map[string]interface{},0)

	page := 1
	for {
		Url, err := url.Parse(apiUrl)
		if err != nil{
			log.Printf("GetFollowingUidListByUid url.Parse error, err=%v",err)
			return retList,err
		}

		params:=url.Values{}
		params.Set("access_token",accessToken)  //这两种都可以
		params.Set("uid",strconv.FormatInt(uid,10))
		params.Set("count","100")
		Url.RawQuery = params.Encode()				//如果参数中有中文参数,这个方法会进行URLEncode
		urlPath := Url.String()
		response, err := http.Get(urlPath)
		if err!= nil {
			log.Printf("GetWeiboListByUid http.Get(urlPath) error, err=%v",err)
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
		statuses := respMap["statuses"].([]interface{})
		if len(statuses) == 0{
			break
		}else{			//page暂时不管
			page++
		}

		for _,_value := range statuses{
			value := _value.(map[string]interface{})
			weiboTime,err := time.Parse("Mon Jan 02 03:04:05 +0800 2006",value["created_at"].(string))
			if err != nil {
				weiboTime = time.Now()
			}
			weiboDetail := make(map[string]interface{})
			weiboDetail["content"] = value["text"]
			weiboDetail["time"] = weiboTime.Unix()
			retList = append(retList, weiboDetail)
		}
		time.Sleep(time.Duration(500)*time.Microsecond)
		break

	}

	return retList, nil
}






