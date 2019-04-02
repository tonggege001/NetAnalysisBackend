package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var MysqlDB *sql.DB

type VideoInfo struct {
	VideoUrl			string		`json:"video_url"`
	UpName				string		`json:"up_name"`
	VideoName			string		`json:"video_name"`
	ImageUrl			string		`json:"image_url"`
	TimeLength			int64		`json:"time_length"`
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
	ret := make([]VideoInfo,0)
	tx := MysqlDB
	MetaRows, err :=tx.Query("select vid,video_url,image_url from BStation.video_meta where vid in (?)",vidList)
	if err != nil{
		log.Printf("GetVideoDataListByVidList tx.Query MetaRows error,  err=%v",err)
		return ret,err
	}


	for MetaRows.Next(){
		meta := VideoInfo{}
		vid := 0
		err := MetaRows.Scan(&vid, &meta.VideoUrl, &meta.ImageUrl)

		if err != nil{
			log.Printf("VerifyLogin Rows.Scan error, Rows=%v, err=%v",MetaRows,err)
			return ret,err
		}

		DetailRows, err := tx.Query("select video_name,up_name,time_length from BStation.video_meta where vid = ?",vid)
		if err != nil{
			return ret, err
		}

		for DetailRows.Next(){
			err := DetailRows.Scan(&meta.VideoName, &meta.UpName, &meta.TimeLength)
			if err != nil{
				return  ret, err
			}
		}

		ret = append(ret,meta)
	}
	return ret,nil
}







