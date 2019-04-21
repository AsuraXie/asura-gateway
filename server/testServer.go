package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TestServer(){
	engin := gin.Default()
	engin.POST("/postmethod",PostHandler)
	engin.GET("/getmethod",GetHandler)

	engin.Run("localhost:8089")
}

func PostHandler(c *gin.Context){
	time.Sleep(2*time.Second)
	c.JSON(http.StatusOK,c.Request.Header)
}

func GetHandler(c *gin.Context){
	c.JSON(http.StatusOK,c.Request.Header)
}