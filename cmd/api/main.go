package main

import (
	"github.com/gin-gonic/gin"

	"github.com/ddddddO/radigo/lib"
)

var (
	stationId string
	ft        string
	to        string
)

func main() {
	r := gin.Default()

	r.GET("/health", health)
	r.POST("/get_m3u8", handler)

	r.Run(":8888")

}

func health(ctx *gin.Context) {
	ctx.String(200, "health cheack ok!\n")
	return
}

func handler(ctx *gin.Context) {
	stationId = "MBS"
	ft = "20190108050000"
	to = "20190108060000"

	email := ctx.PostForm("email")
	pass := ctx.PostForm("pass")

	c := lib.NewClient()

	err := c.Login(email, pass)
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

	m3u8, err := c.GetTimeFreeM3U8(stationId, ft, to, token)
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	dest, err := lib.Ffmpeg(m3u8, stationId)
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	ctx.String(200, dest+"\n")

	return
}
