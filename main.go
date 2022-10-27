package main

import (
	"fmt"
	"generate/config"
	"generate/utils"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("加载配置文件失败")
		return
	}
	fmt.Println("读取配置文件成功,开始生成文件!")

	// 生成model文件
	utils.CreateFiles(config.YamlConfig)
}
