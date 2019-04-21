package main

import (
	//"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("DELIVERY")
	//copyM4AatLocal()

	r := gin.Default()
	r.GET("/m4a", returnM4AFunc)
	r.Run(":8888")

}

func returnM4AFunc(ctx *gin.Context) {
	m4a := "./new.m4a"
	f, err := os.Open(m4a)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	fInf, err := f.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	ctx.DataFromReader(200, fInf.Size(), "audio/mp4", f, map[string]string{"Content-Disposition": `attachment; filename="downloaded.m4a"`})
	//ctx.DataFromReader(200, fInf.Size(), "application/octet-stream", f, map[string]string{"Content-Disposition": `attachment; filename="downlorded.m4a"`})
	return
}

func copyM4AatLocal() {
	m4a := "../api/MBS.m4a"
	f, err := os.Open(m4a)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// コピーできた。ので、ファイルサイズ大きければここで分割できる？
	newF, err := os.Create("./new.m4a")
	if err != nil {
		log.Fatalln(err)
	}
	defer newF.Close()

	io.Copy(newF, f)

}
