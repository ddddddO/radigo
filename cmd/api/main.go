package main

import (
	"github.com/gin-gonic/gin"

	rg "github.com/ddddddO/radigo/radigo"
)

var (
	stationId string
	ft        string
	to        string
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
	ctx.String(200, "health cheack ok!\n")
	return
}

func handler(ctx *gin.Context) {
	stationId = "MBS"
	ft = "20190108050000"
	to = "20190108060000"

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

	m3u8, err := c.GetTimeFreeM3U8(stationId, ft, to, token)
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	dest, err := rg.Ffmpeg(m3u8, stationId)
	if err != nil {
		ctx.String(404, err.Error()+"\n")
		return
	}

	ctx.String(200, dest+"\n")

	return
}
