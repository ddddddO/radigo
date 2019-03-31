package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	//"github.com/gin-gonic/gin"
)

// TODO:.m4aファイルは専用ディレクトリに生成する
func Ffmpeg(m3u8, stationId string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fileName := "%s.m4a"
	destPath := filepath.Join(cwd, fmt.Sprintf(fileName, stationId))
	if IsExist(destPath) {
		if err := os.Remove(destPath); err != nil {
			return "", err
		}
	}

	ffmpeg := `ffmpeg -i %s -t 17 -c copy %s`
	cmd := exec.Command(
		"sh", "-c",
		fmt.Sprintf(ffmpeg, m3u8, destPath),
	)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return destPath, nil
}

/*
// NOTE:client側へ.m4aのアップロードを試みる関数。ginのみでは無理そう
// 参考：https://medium.com/eureka-engineering/multipart-file-upload-in-golang-c4a8eb15a3ee
func Upload(ctx *gin.Context, path string) error {
	// debug
	//	tmp := "/mnt/c/DEV/workspace/GO/src/github.com/ddddddO/radigo/cmd/api/tmp.txt"
	//	path = tmp

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	finf, err := f.Stat()
	if err != nil {
		return err
	}

	buf := make([]byte, finf.Size())
	_, err = f.Read(buf)
	if err != nil {
		return err
	}

	fmt.Printf("CliantIP: %s\n", ctx.ClientIP())

	// Content-type
	// とは :https://wa3.i-3-i.info/word15787.html
	// 種類 :https://developer.mozilla.org/ja/docs/Web/HTTP/Basics_of_HTTP/MIME_types#Important_MIME_types_for_Web_developers
	// 種類 :https://www.tagindex.com/html5/basic/mimetype.html
	ctx.Data(200, "audio/aac", buf)
	//ctx.Data(200, "text/plain", buf)

	return nil
}
*/
