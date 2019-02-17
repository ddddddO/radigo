package lib

import (
	"os"
)

func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
