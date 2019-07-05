package models

import (
	"database/sql"
	"encoding/json"
	"errors"
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

type hashInfo struct {
	ID int
	Hash string
	Consult bool
	ObjGetCount int
	Hits int
	CreateTime string
}

type spiderInfo struct {
	ID int
	NodeId string
	UpdateTime string
	SpiderName string
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

func (dbw *DbWorker) InsertData(hash string) (err error){
	if(len(hash)<40 ||len(hash)>64){
		fmt.Println(len(hash))
		err = errors.New("hash is wrong!")
		return
	}
	hashObj := new(hashInfo)
	row := dbw.DB.QueryRow("SELECT * FROM Hash_List WHERE Hash=?",hash)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	err =row.Scan(&hashObj.ID,&hashObj.Hash,&hashObj.Consult,&hashObj.ObjGetCount,&hashObj.Hits,&hashObj.CreateTime);
	if err != nil{
		_, err = dbw.DB.Exec("INSERT INTO Hash_List (Hash, Consult,ObjGetCount,CreateTime) VALUES (?, false ,0,Now())",hash)

	}else{
		_, err = dbw.DB.Exec("UPDATE Hash_List SET Hits = Hits+1 WHERE Hash = ? ",hash)
	}
	if err != nil {
		fmt.Printf("insert data to Hash_List error: %v\n", err)
		return
	}
	return
}

func (dbw *DbWorker) UpdateSpider(nodeId string,spiderName string) (err error){
	if(len(nodeId)<40 ||len(nodeId)>64){
		fmt.Println(len(nodeId))
		err = errors.New("Nodeid is wrong!")
		return
	}
	spiderObj := new(spiderInfo)
	row := dbw.DB.QueryRow("SELECT * FROM Spider_List WHERE NodeId=?",nodeId)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	err =row.Scan(&spiderObj.ID,&spiderObj.NodeId,&spiderObj.UpdateTime,&spiderObj.SpiderName);
	if err != nil{
		_, err = dbw.DB.Exec("INSERT INTO Spider_List (NodeId,UpdateTime,SpiderName) VALUES (?, Now(),?)",nodeId,spiderName)

	}else{
		_, err = dbw.DB.Exec("UPDATE Spider_List SET UpdateTime = Now() WHERE ID = ? ",spiderObj.ID)
	}
	if err != nil {
		fmt.Printf("insert data to Spider_List error: %v\n", err)
		return
	}
	return
}

func (dbw *DbWorker) CloseConnect(){
	dbw.DB.Close()
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