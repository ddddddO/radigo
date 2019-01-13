package lib

import (
	//"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// TODO:.m4aファイルは専用ディレクトリに生成する
func Ffmpeg(m3u8, stationId string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fileName := "%s.m4a"
	dest := filepath.Join(cwd, fmt.Sprintf(fileName, stationId))
	if IsExist(dest) {
		err = os.Remove(dest)
		if err != nil {
			return "", err
		}
	}

	ffmpeg := `ffmpeg -i %s -t 17 -c copy %s`
	cmd := exec.Command(
		"sh", "-c",
		fmt.Sprintf(ffmpeg, m3u8, dest),
	)

	/*
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	*/

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	/*
		fmt.Println("-stdout-")
		fmt.Println(stdout.String())

		fmt.Println("-stderr-")
		fmt.Println(stderr.String())
	*/

	return dest, nil
}
