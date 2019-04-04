package apis

import (
	"NetAnalysisBackend/call"
	"NetAnalysisBackend/utils"
	"log"
	"net/http"
	"strconv"
)

func GetRecommendMusicVideo(w http.ResponseWriter, r *http.Request)  {
	defer utils.RecoverResolve()
	retMap := make(map[string]interface{})

	err := r.ParseForm()
	if err!= nil{
		log.Printf("GetRecommendMusicVideo r.ParseForm() error, err=%v",err)
		retMap["code"] = utils.FORMPARSERR
		utils.SendJson(&retMap,w)
		return
	}

	platformIdList, fail:=r.Form["platform_id"]
	if !fail || len(platformIdList) == 0{
		log.Printf("GetRecommendMusicVideo r.Form[\"platform_id\"] error, fail=%v, len=%v",fail,len(platformIdList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	codeList, fail:=r.Form["uid"]
	if !fail || len(codeList) == 0 {
		log.Printf("GetRecommendMusicVideo r.Form[\"uid\"] error, fail=%v, len=%v",fail,len(codeList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	uid := codeList[0]					//string类型
	platformId := platformIdList[0]

	//获取用户关注的列表
	accessToken,err :=utils.GetRedis(platformId+"_access_token")
	if err != nil{
		log.Printf("GetRecommendMusicVideo GetRedis error,key=%v, err=%v",platformId+"_access_token",err)
		retMap["code"] = utils.REDISERR
		utils.SendJson(&retMap,w)
		return
	}

	//int64位
	uidInt64,err := strconv.ParseInt(uid,10,64)
	if err != nil {
		log.Printf("GetRecommendMusicVideo uid ParseInt error, value=%v,err=%v",uid,err)
		retMap["code"] = utils.MATHERR
		utils.SendJson(&retMap,w)
		return
	}

	uidList, err := call.GetFollowingUidListByUid(uidInt64, accessToken)
	if err != nil{
		log.Printf("GetRecommendMusicVideo GetFollowingUidListByUid error, err=%v, uidInt64=%v, accessToken=%v",err,uidInt64,accessToken)
		if err.Error() == "not Login"{
			retMap["code"] = utils.NOTLOGIN
		}else {
			retMap["code"] = utils.PARAMERR
		}
		utils.SendJson(&retMap,w)
		return
	}

	weiboList, err := call.GetWeiboListByUid(uidInt64, accessToken)
	if err != nil {
		log.Printf("GetRecommendMusicVideo GetWeiboListByUid error, err=%v",err)
		if err.Error() == "not Login"{
			retMap["code"] = utils.NOTLOGIN
		}else{
			retMap["code"] = utils.PARAMERR
		}
		utils.SendJson(&retMap,w)
		return
	}

	videoVidList, err := call.GetRecommendMusicVideoCall(uidInt64, uidList,weiboList)
	if err != nil {
		log.Printf("GetRecommendMusicVideo call.GetRecommendMusicVideoCall error, err=%v", err)
		retMap["code"] = utils.CALLERR
		utils.SendJson(&retMap,w)
		return
	}

	videoDataList,err :=utils.GetVideoDataListByVidList(videoVidList)
	if err != nil{
		log.Printf("GetRecommendMusicVideo GetVideoDataListByVidList error, err=%v", err)
		retMap["code"] = utils.DATABASEERR
		utils.SendJson(&retMap,w)
		return
	}

	retMap["code"] = utils.OK
	retMap["data"] = videoDataList
	utils.SendJson(&retMap,w)
	return
}


func GetRecommendMusicUp(w http.ResponseWriter, r *http.Request){
	defer utils.RecoverResolve()
	retMap := make(map[string]interface{})

	err := r.ParseForm()
	if err!= nil{
		log.Printf("GetRecommendMusicUp r.ParseForm() error, err=%v",err)
		retMap["code"] = utils.FORMPARSERR
		utils.SendJson(&retMap,w)
		return
	}

	platformIdList, fail:=r.Form["platform_id"]
	if !fail || len(platformIdList) == 0{
		log.Printf("GetRecommendMusicUp r.Form[\"platform_id\"] error, fail=%v, len=%v",fail,len(platformIdList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	codeList, fail:=r.Form["uid"]
	if !fail || len(codeList) == 0 {
		log.Printf("GetRecommendMusicUp r.Form[\"uid\"] error, fail=%v, len=%v",fail,len(codeList))
		retMap["code"] = utils.PARAMERR
		utils.SendJson(&retMap,w)
		return
	}

	uid := codeList[0]					//string类型
	platformId := platformIdList[0]

	//获取用户关注的列表
	accessToken,err :=utils.GetRedis(platformId+"_access_token")
	if err != nil{
		log.Printf("GetRecommendMusicUp GetRedis error,key=%v, err=%v",platformId+"_access_token",err)
		retMap["code"] = utils.REDISERR
		utils.SendJson(&retMap,w)
		return
	}

	//int64位
	uidInt64,err := strconv.ParseInt(uid,10,64)
	if err != nil {
		log.Printf("GetRecommendMusicUp uid ParseInt error, value=%v,err=%v",uid,err)
		retMap["code"] = utils.MATHERR
		utils.SendJson(&retMap,w)
		return
	}

	uidList, err := call.GetFollowingUidListByUid(uidInt64, accessToken)
	if err != nil{
		log.Printf("GetRecommendMusicUp GetFollowingUidListByUid error, err=%v, uidInt64=%v, accessToken=%v",err,uidInt64,accessToken)
		if err.Error() == "not Login"{
			retMap["code"] = utils.NOTLOGIN
		}else {
			retMap["code"] = utils.PARAMERR
		}
		utils.SendJson(&retMap,w)
		return
	}

	weiboList, err := call.GetWeiboListByUid(uidInt64, accessToken)
	if err != nil {
		log.Printf("GetRecommendMusicUp GetWeiboListByUid error, err=%v",err)
		if err.Error() == "not Login"{
			retMap["code"] = utils.NOTLOGIN
		}else{
			retMap["code"] = utils.PARAMERR
		}
		utils.SendJson(&retMap,w)
		return
	}

	upMasterList, err := call.GetRecommendMusicUp(uidInt64, uidList,weiboList)
	if err != nil {
		log.Printf("GetRecommendMusicUp call.GetRecommendMusicUp error, err=%v", err)
		retMap["code"] = utils.CALLERR
		utils.SendJson(&retMap,w)
		return
	}

	upMasterDataList,err := utils.GetUpMasterDataListByUpIDList(upMasterList)
	if err != nil{
		log.Printf("upMasterDataList GetUpMasterDataListByUpIDList error, err=%v", err)
		retMap["code"] = utils.DATABASEERR
		utils.SendJson(&retMap,w)
		return
	}

	retMap["code"] = utils.OK
	retMap["data"] = upMasterDataList
	utils.SendJson(&retMap,w)
	return

}

