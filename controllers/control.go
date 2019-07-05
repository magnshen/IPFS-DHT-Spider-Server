package controllers

import (
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)
type configuration struct {
	USERNAME string
	PASSWORD string
	NETWORK string
	SERVER  string
	PORT    int
	DATABASE string
}
var Setting configuration
var dsn string

func init(){
	file, _ := os.Open("./config.cnf")
	defer file.Close()
	decoder := json.NewDecoder(file)
	Setting = configuration{}
	err := decoder.Decode(&Setting)
	if err != nil {
		fmt.Println("config file load error")
		os.Exit(2)
	}
	dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s",Setting.USERNAME,Setting.PASSWORD,Setting.NETWORK,Setting.SERVER,Setting.PORT,Setting.DATABASE)
}

func Index(c *gin.Context){

	dbwk ,err := models.NewDbWorker(dsn)
	if err != nil{
	 	c.String(http.StatusOK,"database connect error")
		return
	}
	news ,days,err := dbwk.GetAll()
	dbwk.CloseConnect()
	index_data := map[string]interface{}{
		"total_hashs": news.Total_hashs,
		"yesterday_hashs":news.Yesterday_hashs,
		"today_hashs" : news.Today_hashs,
		"spiders" : news.Spiders,
	}

	jsdata_date := [][]string{}
	jsdata_hash := [][]string{}
	i := 0
	for i < len(days.Days_line){
		num := strconv.Itoa(i+1)
		jsdata_hash = append(jsdata_hash,[]string{num,days.Days_line[i].Hashs})
		jsdata_date = append(jsdata_date,[]string{num,days.Days_line[i].Date_time})
		i++
	}

	index_data["js1_data"] = jsdata_hash
	index_data["js2_data"] = jsdata_date
	news_data := [][]string{}
	//index_data["New_hash"] = [][]string{}
	i = 0
	for i < len(news.New_hash){
		news_data = append(news_data,[]string{news.New_hash[i].Hash,news.New_hash[i].CreateTime})
		i++
	}
	index_data["New_hash"] = news_data
	history_heats := [][]string{}
	lastweek_heats:= [][]string{}
	i = 0
	for i < len(days.History_heats){
		history_heats = append(history_heats,[]string{days.History_heats[i].Hash,days.History_heats[i].Heat})
		i++
	}
	i = 0
	for i < len(days.Lastweek_heats){
		lastweek_heats = append(lastweek_heats,[]string{days.Lastweek_heats[i].Hash,days.Lastweek_heats[i].Heat})
		i++
	}
	index_data["History_heats"] = history_heats
	index_data["Lastweek_heats"] = lastweek_heats
	c.HTML(http.StatusOK, "pc.html", index_data)
}

func GetNews(c *gin.Context){
	dbwk ,err := models.NewDbWorker(dsn)
	if err != nil{
		c.String(http.StatusOK,"database connect error")
		return
	}
	news ,err := dbwk.GetNewsString()
	dbwk.CloseConnect()
	c.String(http.StatusOK,news)
}

func Submit(c *gin.Context){
	hashs := c.PostForm("hashs")
	nodeId := c.PostForm("nodeId")
	spiderName := c.PostForm("spiderName")
	hashList := []string{}
	json.Unmarshal([]byte(hashs),&hashList)
	//dsn :=  fmt.Sprintf("%s:%s@%s(%s:%d)/%s",Setting.USERNAME,Setting.PASSWORD,Setting.NETWORK,Setting.SERVER,Setting.PORT,Setting.DATABASE)
	dbwk ,err:= models.NewDbWorker(dsn)
	if err != nil{
		c.String(http.StatusNotImplemented, "Open mysql failed,err:%v\n")
		return
	}
	for i:=0;i< len(hashList);i++{
		dbwk.InsertData(hashList[i])
	}
	e := dbwk.UpdateSpider(nodeId,spiderName)
	if e != nil{
		fmt.Println(e.Error())
	}
	dbwk.CloseConnect()
	c.String(http.StatusOK, "Success")
}
func GetIpfsObject(c *gin.Context){
	hash := c.Param("hash")
	dbwk ,err := models.NewDbWorker(dsn)
	if err != nil{
		c.String(http.StatusOK,"database connect error")
		return
	}
	obj ,err := dbwk.GetIpfsObject(hash)
	if err != nil{
		c.String(http.StatusOK,"NULL")
		return
	}
	dbwk.CloseConnect()
	c.String(http.StatusOK,obj)
}














