package main

import (
	"./controllers"
	"github.com/gin-gonic/gin"
	"os"
)


var ff *os.File

func main() {

	e := gin.Default()

	e.Static("/static", "./static")
	e.LoadHTMLGlob("views/*")

	e.GET("/",controllers.Index)
	e.GET("/api/getNews",controllers.GetNews)
	e.POST("/api/submit",controllers.Submit)
	e.GET("/ipfsObject/:hash",controllers.GetIpfsObject)
	// Start server
	e.Run(":80")


}
