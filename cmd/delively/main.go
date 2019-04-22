package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

func main() {
	fmt.Println("DELIVERY")
	//copyM4AatLocal()
	/*
		r := gin.Default()

		r.LoadHTMLGlob("./templates/*")
		r.GET("/index", renderHTMLFunc)

		r.GET("/m4a", returnM4AFunc)

		r.Run(":8888")
	*/

	// NOTE: https://github.com/gin-gonic/gin#run-multiple-service-using-gin
	serverHTML := &http.Server{
		Addr:         ":9999",
		Handler:      routerHTML(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverAPI := &http.Server{
		Addr:         ":8888",
		Handler:      routerAPI(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return serverHTML.ListenAndServe()
	})

	g.Go(func() error {
		return serverAPI.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatalln(err)
	}

}

func routerHTML() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	e.LoadHTMLGlob("./templates/*")
	e.GET("/index", renderHTMLFunc)

	return e
}

func routerAPI() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/m4a", returnM4AFunc)

	return e
}

var dlPath = "pathpath"

func renderHTMLFunc(ctx *gin.Context) {
	ctx.HTML(200, "index.tmpl", gin.H{
		"title": "速報とかラジオで",
		// TODO: サーバー側のm4aのダウンロードパスをレンダリングする
		// "dlPath": cmd/api/main.goのdestPath,
		"dlPath": dlPath,
	})
}

// NOTE: ブラウザからlocalhost:8888/m4a でDLすること成功し、再生もできた。
//       ターミナル・コマンドプロンプトで curl -i localhost:8888/m4a --output downloaded.m4a だと、DLできても、ファイルが壊れているらしく再生できなかった。
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

	dlPath += m4a

	// 参考：https://qiita.com/yuji38kwmt/items/9edb4b17768d112ae43b
	//     ：https://github.com/gin-gonic/gin
	ctx.DataFromReader(200, fInf.Size(), "audio/mp4", f, map[string]string{"Content-Disposition": `attachment; filename="downloaded.m4a"`})
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
