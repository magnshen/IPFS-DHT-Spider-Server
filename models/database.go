package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
type Hash struct {
	Hash string
	Heat string
	CreateTime string
}

type NewsInfo struct {
	Total_hashs string
	Yesterday_hashs string
	Today_hashs string
	Spiders int
	New_hash []Hash
}
type DaysLinePoint struct {
	Date_time string
	Hashs string
}
type DaysInfo struct {
	Days_line []DaysLinePoint
	History_heats []Hash
	Lastweek_heats []Hash
}
type DbWorker struct {
	DB       *sql.DB
}
func (self *DbWorker)GetAll()(new NewsInfo,days DaysInfo,err error){
	row := self.DB.QueryRow("SELECT * FROM Web_Data WHERE ID = 1")
	var id int
	var news_info_str,days_info_str string
	err = row.Scan(&id,&news_info_str,&days_info_str);
	if err != nil{
		return
	}

	json.Unmarshal([]byte(news_info_str),&new)
	json.Unmarshal([]byte(days_info_str),&days)
	return
}

func (self *DbWorker)GetNewsString()(news string,err error){
	row := self.DB.QueryRow("SELECT * FROM Web_Data WHERE ID = 1")
	var id int
	var days_info_str string
	err = row.Scan(&id,&news,&days_info_str);
	return
}

func (self *DbWorker)GetIpfsObject(hash string)(obj string,err error){
	row := self.DB.QueryRow("SELECT * FROM Hash_Obj WHERE Hash = ?",hash)
	var id int
	var hash_ string
	err = row.Scan(&id,&hash_,&obj);
	return
}

func NewDbWorker(dsn string)(*DbWorker,error){
	//dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s","root","123456","tcp","192.168.2.56",3306,"Spider")
	dbwk := &DbWorker{}
	DB,err := sql.Open("mysql",dsn)
	if err != nil{
		fmt.Printf("Open mysql failed,err:%v\n",err)
		return dbwk ,err
	}
	DB.SetConnMaxLifetime(100*time.Second)  //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(100)//设置最大连接数
	DB.SetMaxIdleConns(16) //设置闲置连接数
	dbwk.DB=DB
	return dbwk ,err
}