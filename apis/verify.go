package apis

import (
	"NetAnalysisBackend/call"
	"NetAnalysisBackend/utils"
	"log"
	"net/http"
)


func WeiboCall(w http.ResponseWriter, r *http.Request){
	defer utils.RecoverResolve()
	retMap := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil{
		log.Printf("WeiboCall r.ParseForm() error, err=%v",err)
		retMap["code"] = utils.FORMPARSERR
		utils.SendJson(&retMap,w)
		return
	}

	codeList, fail:=r.Form["code"]
	if !fail || len(codeList) == 0{
		log.Printf("WeiboCall r.Form[\"codeList\"] error, fail=%v, len=%v",fail,len(codeList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	platformIdList, fail := r.Form["state"]
	if !fail || len(platformIdList) == 0{
		log.Printf("WeiboCall r.Form[\"platform_id\"] error, fail=%v, len=%v",fail,len(platformIdList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return

	}

	code := codeList[0]
	platformId := platformIdList[0]

	err = utils.SetRedis(platformId+"_code",code)
	if err != nil{
		log.Printf("WeiboCall SetRedis error, err=%v", err)
		retMap["code"] = utils.REDISERR
		utils.SendJson(&retMap,w)
		return
	}

	retMap["code"] = utils.OK
	utils.SendJson(&retMap, w)
	return

}



func VerifyWeibo(w http.ResponseWriter, r *http.Request){
	defer utils.RecoverResolve()
	retMap := make(map[string]interface{})

	err := r.ParseForm()
	if err!= nil{
		log.Printf("VerifyWeibo r.ParseForm() error, err=%v",err)
		retMap["code"] = utils.FORMPARSERR
		utils.SendJson(&retMap,w)
		return
	}

	platformIdList, fail:=r.Form["platform_id"]
	if !fail || len(platformIdList) == 0{
		log.Printf("VerifyWeibo r.Form[\"platform_id\"] error, fail=%v, len=%v",fail,len(platformIdList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return

	}
	platformId := platformIdList[0]
	redirectUrl,err := call.VerifyWeiboCall(platformId)

	if err !=nil {
		log.Printf("VerifyWeibo VerifyWeiboCall error, err=%v",err)
		retMap["code"] = utils.CODEERR
		utils.SendJson(&retMap,w)
		return
	}

	retMap["code"] = utils.OK
	retMap["url"] = redirectUrl
	utils.SendJson(&retMap,w)
	return
}


func GetVerifiedAndUserInfo(w http.ResponseWriter, r *http.Request){
	defer utils.RecoverResolve()
	retMap := make(map[string]interface{})

	err := r.ParseForm()
	if err!= nil{
		log.Printf("GetVerifiedAndUserInfo r.ParseForm() error, err=%v",err)
		retMap["code"] = utils.FORMPARSERR
		utils.SendJson(&retMap,w)
		return
	}

	platformIdList, fail:=r.Form["platform_id"]
	if !fail || len(platformIdList) == 0{
		log.Printf("GetVerifiedAndUserInfo r.Form[\"platform_id\"] error, fail=%v, len=%v",fail,len(platformIdList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	codeList, fail:=r.Form["code"]
	if !fail || len(codeList) == 0 {
		log.Printf("GetVerifiedAndUserInfo r.Form[\"codeList\"] error, fail=%v, len=%v",fail,len(codeList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	code := codeList[0]
	platformId := platformIdList[0]

	if code != "0"{
		log.Printf("GetVerifiedAndUserInfo code != 0, code=%v",code)
		retMap["code"] = utils.CODEERR
		utils.SendJson(&retMap,w)
		return
	}

	authMap, err := call.GetWeiboVerified(platformId, utils.APPID, utils.APPSECRET, "authorization_code",)
	if err!= nil{
		log.Printf("GetVerifiedAndUserInfo GetWeiboVerified error, err=%v", err)
		retMap["code"] = utils.NOTLOGIN
		utils.SendJson(&retMap,w)
		return
	}

	err = utils.SetRedisByEx(platformId+"_access_token",authMap["access_token"].(string),(int)(authMap["expires_in"].(float64)))
	if err != nil{
		log.Printf("GetVerifiedAndUserInfo SetRedis error, err=%v", err)
		retMap["code"] = utils.REDISERR
		utils.SendJson(&retMap, w)
		return
	}

	uid := authMap["uid"].(string)
	userInfoMap,err := call.GetUserInfo(authMap["access_token"].(string), uid)
	if err != nil{
		log.Printf("GetVerifiedAndUserInfo GetUserInfo error, err=%v", err)
		retMap["code"] = utils.APICALLERR
		utils.SendJson(&retMap,w)
		return
	}
	err = utils.SetRedisByEx(platformId+"_uid",uid,60*60)
	if err!= nil{
		log.Printf("GetVerifiedAndUserInfo SetRedis error, err=%v", err)
		retMap["code"] = utils.REDISERR
		utils.SendJson(&retMap, w)
		return
	}

	imageAvater := userInfoMap["profile_image_url"].(string)
	userName := userInfoMap["name"].(string)

	retMap["code"] = utils.OK
	retMap["uid"] = uid
	retMap["image_avater"] = imageAvater
	retMap["user_name"] = userName
	utils.SendJson(&retMap,w)
	return
}

func QuitLogin(w http.ResponseWriter, r *http.Request){
	defer utils.RecoverResolve()
	retMap := make(map[string]interface{})

	err := r.ParseForm()
	if err!= nil{
		log.Printf("QuitLogin r.ParseForm() error, err=%v",err)
		retMap["code"] = utils.FORMPARSERR
		utils.SendJson(&retMap,w)
		return
	}

	platformIdList, fail:=r.Form["platform_id"]
	if !fail || len(platformIdList) == 0{
		log.Printf("QuitLogin r.Form[\"platform_id\"] error, fail=%v, len=%v",fail,len(platformIdList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	platformId := platformIdList[0]
	err = utils.DeleteRedis(platformId+"_access_token")
	if err!= nil{
		log.Printf("QuitLogin utils.DeleteRedis, key=%v, err=%v",platformId+"_access_token",err)
		retMap["code"] = utils.REDISERR
		utils.SendJson(&retMap,w)
		return
	}

	err = utils.DeleteRedis(platformId+"_uid")
	if err!= nil {
		log.Printf("QuitLogin utils.DeleteRedis, key=%v, err=%v",platformId+"_uid",err)
		retMap["code"] = utils.REDISERR
		utils.SendJson(&retMap,w)
		return
	}

	retMap["code"] = utils.OK
	utils.SendJson(&retMap,w)
	return
}












