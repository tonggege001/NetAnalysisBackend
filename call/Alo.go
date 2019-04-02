package call

import (
	"NetAnalysisBackend/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const ALODOMAIN = "http://127.0.0.1:2048"
//const ALODOMAIN = "http://10.15.150.82:2048"
func GetRecommendMusicVideoCall(uid int64, fowList []int64, weiBoList[]map[string]interface{}) ([]int, error){
	defer utils.RecoverResolve()
	retList := make([]int,0)

	sendMap := make(map[string]interface{})
	sendMap["user_uid"] = uid
	sendMap["following"] = fowList
	sendMap["weibo"] = weiBoList
	byteJson,err := json.Marshal(sendMap)
	if err != nil{
		log.Printf("GetRecommendMusicVideoCall Marshal sendMap=%v",sendMap)
		return retList,err
	}

	body := bytes.NewReader(byteJson)
	client := &http.Client{}
	request, err := http.NewRequest("POST", ALODOMAIN, body) //建立一个请求
	if err!= nil{
		log.Printf("GetRecommendMusicVideoCall http.NewRequest error body=%v, err=%v",body,err)
		return retList, err
	}

	response, err := client.Do(request) //提交
	if err!= nil {
		log.Printf("GetRecommendMusicVideoCall client.Do(request) error, err=%v",err)
		return retList,err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("GetRecommendMusicVideoCall ioutil.ReadAll(response.Body) error, err=%v", err)
		return retList, err
	}

	respMap := make(map[string]interface{})
	err =json.Unmarshal(respBody, &respMap)
	if err != nil {
		log.Printf("GetRecommendMusicVideoCall json.Unmarshal(respBody, &respMap) error")
		return retList, err
	}

	if respMap["code"].(int) == 0{
		log.Printf("GetRecommendMusicVideoCall call get code not 0")
		return retList,errors.New(fmt.Sprintf("getError code, code=%v",respMap["code"].(int)))
	}

	retList = respMap["video_vid"].([]int)
	return retList,nil
}



