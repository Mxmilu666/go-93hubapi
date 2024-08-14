package main

import (
	"fmt"
	"net/http"
	"os"

	"go-93hubapi/source"
)

func main() {
	fmt.Printf("Go-93HubAPI v0.0.1 \n")
	configFile := "config.yml"

	// 检查配置文件是否存在，如果不存在则创建默认配置文件
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err := source.CreateDefaultConfig(configFile)
		if err != nil {
			fmt.Printf("Error creating default config file: %v\n", err)
			return
		}
		fmt.Printf("Created default config file: %s. Please edit it with your configuration.\n", configFile)
		return
	}

	// 读取配置文件
	config, err := source.ReadConfig(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	// 下载仓库文件
	err = source.DownloadRepoFiles(config)
	if err != nil {
		fmt.Printf("Error downloading repository files: %v\n", err)
		return
	}

	// 启动 HTTP 服务器
	address := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	fmt.Printf("Starting server at %s\n", address)
	http.ListenAndServe(address, source.FileServer(config.Dest))
}
