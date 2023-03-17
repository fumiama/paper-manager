package utils

import "os"

// IsExist 文件/路径存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsNotExist 文件/路径不存在
func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	return err != nil && os.IsNotExist(err)
}

// FileSize 获取文件大小
func FileSize(path string) (n int64) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	n = stat.Size()
	return
}
