package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"self-signed-certificate/controllers"
)

const outPath = "./out/"
const resultFileName = "./result.log"

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/submit", func(c *gin.Context) {
		controllers.ExecCommand(c, resultFileName)
	})

	r.GET("/result", func(c *gin.Context) {
		controllers.ViewResult(c, resultFileName)
	})

	r.GET("/files", func(c *gin.Context) {
		controllers.GetFiles(c, outPath)
	})
	r.GET("/view/*filename", func(c *gin.Context) {
		controllers.ViewFiles(c, outPath)
	})
	r.POST("/delete/*filename", func(c *gin.Context) {
		controllers.DeleteFiles(c, outPath)
	})
	r.GET("/download/*filename", func(c *gin.Context) {
		controllers.DownloadFiles(c, outPath)
	})

	err := r.Run(":8000")
	if err != nil {
		return
	}
}
