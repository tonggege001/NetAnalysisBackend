package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

var MysqlDB *sql.DB

type VideoInfo struct {				//video的详细信息
	VideoUrl			string		`json:"video_url"`
	UpName				string		`json:"up_name"`
	VideoName			string		`json:"video_name"`
	ImageUrl			string		`json:"image_url"`
	TimeLength			int64		`json:"time_length"`
}

type UpMasterInfo struct {
	Uid					string		`json:"uid"`
	Name				string		`json:"name"`
	UrlAvater			string		`json:"url_avater"`
	FollowerCount		int			`json:"follower_count"`
	VideoCount			int			`json:"video_count"`
	UpText				string		`json:"up_text"`

}



func InitMysql(){
	var err error
	MysqlDB, err = sql.Open("mysql", "root:Dachuang123@tcp(47.106.174.32:3306)/?charset=utf8")

	if err != nil{
		log.Fatalf("InitMysql error, err=%v",err)
	}

	MysqlDB.SetMaxOpenConns(4)
	return
}


func GetVideoDataListByVidList(vidList []int)([]VideoInfo, error){
	defer Recover2("GetVideoDataListByVidList")
	ret := make([]VideoInfo,0)
	tx := MysqlDB
	for _,vid := range vidList{
		MetaRows, err :=tx.Query("select video_url,image_url from BStation.video_meta where vid = ?",vid)
		if err != nil{
			log.Printf("GetVideoDataListByVidList tx.Query MetaRows error,  err=%v",err)
			return ret,err
		}
		meta := VideoInfo{}
		for MetaRows.Next(){
			err = MetaRows.Scan(&meta.VideoUrl, &meta.ImageUrl)
			if err != nil{
				log.Printf("VerifyLogin Rows.Scan error, Rows=%v, err=%v",MetaRows,err)
				return ret,err
			}
		}

		_ =MetaRows.Close()
		DetailRows, err := tx.Query("select video_name,up_name,time_length from BStation.video_detail where vid = ?",vid)

		if err != nil{
			return ret, err
		}
		for DetailRows.Next(){
			err = DetailRows.Scan(&meta.VideoName, &meta.UpName, &meta.TimeLength)
			if err != nil{
				return  ret, err
			}
		}
		_ = DetailRows.Close()
		ret = append(ret,meta)
	}
	return ret,nil
}

func GetUpMasterDataListByUpIDList(upIdList []int)([]UpMasterInfo,error)  {
	defer Recover2("GetUpMasterDataListByUpIDList")
	ret := make([]UpMasterInfo,0)
	tx := MysqlDB
	for _,upid := range upIdList{

		UpRows, err :=tx.Query("select up_name from BStation.video_meta where up_id = ?",upid)
		if err != nil{
			log.Printf("GetUpMasterDataListByUpIDList tx.Query MetaRows error,  err=%v",err)
			return ret,err
		}
		upData := UpMasterInfo{
			Uid:				strconv.FormatInt(int64(upid),10),
			UrlAvater:			"",
			FollowerCount:		-1,
			VideoCount:			1,
			UpText: 			"",
		}

		for UpRows.Next(){
			err = UpRows.Scan(&upData.Name)
			if err != nil{
				log.Printf("VerifyLogin Rows.Scan error, Rows=%v, err=%v",UpRows,err)
				return ret,err
			}
			ret = append(ret,upData)
		}
	}

	return ret,nil
}



