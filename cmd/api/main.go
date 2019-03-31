package main

import (
	"github.com/gin-gonic/gin"

	rg "github.com/ddddddO/radigo/radigo"
)

var (
	stationId string
	start     string
	end       string
)

type Auth struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

func main() {
	r := gin.Default()

	r.GET("/health", health)
	r.POST("/get_m4a", handler)

	r.Run(":8888")

}

func health(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*") // No 'Access-Control-Allow-Origin' header is present on the requested resource. のエラー対策
	ctx.JSON(200, gin.H{"type": "get", "message": "health cheack ok!"})
	return
}

func handler(ctx *gin.Context) {
	stationId = "MBS"
	start = "20190330053000"
	end = "20190330053100"

	var auth Auth
	ctx.BindJSON(&auth)

	c := rg.NewClient()

	err := c.Login(auth.Email, auth.Pass)
	if err != nil {
		ctx.String(403, err.Error()+"\n")
		return
	}

	token, partialKey, err := c.Auth1()
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	err = c.Auth2(token, partialKey)
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	m3u8, err := c.GetTimeFreeM3U8(stationId, start, end, token)
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	destPath, err := rg.Ffmpeg(m3u8, stationId)
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	ctx.String(200, destPath+"\n")

	return
}
