package call

import (
	"NetAnalysisBackend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func VerifyWeiboCall(platformId string)(string,error){
	defer utils.RecoverResolve()

	clientId := utils.APPID
	redirectUri := utils.DOMAIN

	apiUrl := "https://api.weibo.com/oauth2/authorize"
	Url, err := url.Parse(apiUrl)
	if err != nil {
		return "", err
	}

	params:=url.Values{}
	params.Set("client_id",clientId)
	params.Set("redirect_uri",redirectUri)
	params.Set("state",platformId)
	Url.RawQuery = params.Encode()				//如果参数中有中文参数,这个方法会进行URLEncode
	urlPath := Url.String()

	print(urlPath)
	return urlPath,nil
}

func GetWeiboVerified(platformID string, clientId string, clientSecret string, grantType string)(
	map[string]interface{},error){

	defer utils.RecoverResolve()
	apiUrl :="https://api.weibo.com/oauth2/access_token"
	redirectUri := utils.DOMAIN

	code, err := utils.GetRedis(platformID+"_code")
	if err!= nil{
		return nil,errors.New(fmt.Sprintf("GetWeiboVerified no code available, key=%v",platformID+"_code"))
	}

	params:=url.Values{}
	params.Set("client_id",clientId)  //这两种都可以
	params.Set("client_secret",clientSecret)
	params.Set("grant_type",grantType)
	params.Set("code",code)
	params.Set("redirect_uri",redirectUri)

	resp, err:= http.PostForm(apiUrl,params)
	if err != nil{
		return nil, err
	}

	defer resp.Body.Close()
	body, err:= ioutil.ReadAll(resp.Body)
	if err!= nil{
		return nil, err
	}

	retMap := make(map[string]interface{})
	err = json.Unmarshal(body[:], &retMap)
	if err != nil{
		return nil, err
	}
	return retMap,nil

}


func GetUserInfo(accessToken string, uid string)(map[string]interface{},error){
	defer utils.RecoverResolve()
	retMap := make(map[string]interface{})
	apiUrl := "https://api.weibo.com/2/users/show.json"

	Url, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}

	params:=url.Values{}
	params.Set("access_token",accessToken)
	params.Set("uid",uid)
	Url.RawQuery = params.Encode()				//如果参数中有中文参数,这个方法会进行URLEncode
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil{
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	err = json.Unmarshal(body[:], &retMap)
	if err != nil{
		return nil, err
	}

	errCode,ok := retMap["error_code"].(int)
	if ok {
		return nil, errors.New(fmt.Sprintf("GetUserInfo retMap get error code, code=%v",errCode))
	}

	return retMap,nil

}




