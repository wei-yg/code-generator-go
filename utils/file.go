package utils

import (
	"fmt"
	"io"
	"os"
)

// GenerateFile 生成文件，参数为：文件路径名，文件内容，原有的是否覆盖（false为不覆盖）
func GenerateFile(folder string, fileName string, content string, cover bool) {
	filePath := folder + fileName
	// 如果不存在就递归创建目录
	if !checkPathIsExist(folder) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var f *os.File
	var err error
	if checkPathIsExist(filePath) {
		if !cover {
			fmt.Println("文件: " + filePath + " 已存在(未覆盖)")
			return
		}
		f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.Create(filePath)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	}
	fmt.Println("文件: ", filePath, " 已生成！")
}

// 检查文件是否存在
func checkPathIsExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false // 不存在
	}
	return true // 存在
}
