package lib

import (
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
