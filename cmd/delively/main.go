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
