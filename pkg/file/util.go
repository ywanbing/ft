package file

import (
	"os"

	"github.com/google/uuid"
)

// PathExists 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func GenFileName() string {
	u := uuid.New()
	return u.String()

}
