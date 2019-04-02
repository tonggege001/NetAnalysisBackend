package main

import (
	"NetAnalysisBackend/apis"
	"NetAnalysisBackend/utils"
	"log"
	"net/http"
)

func main(){

	utils.InitMysql()
	http.HandleFunc("/login",apis.WeiboCall)								//微博的回调地址
	http.HandleFunc("/login_weibo",apis.VerifyWeibo)					//设置访问的路由
	http.HandleFunc("/login_page",apis.GetVerifiedAndUserInfo)		//用户认证登录
	http.HandleFunc("/quit",apis.QuitLogin)							//用户认证退出
	http.HandleFunc("/music_recommend_video",apis.GetRecommendMusicVideo)		//用户认证登录
	//http.HandleFunc("")
	err := http.ListenAndServe(":3224", nil) 					//设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}



