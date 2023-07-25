package config

import (
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

const configFile = "/XueLang/MaShen/conf/conf.ini"

var File *goconfig.ConfigFile

func init() {
	// 接受当前文件地址
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile

	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	// 带命令行参数的读取文件
	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}

	// 文件系统的读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("配置文件出错：", err)
	}
	if File == nil {
		log.Fatal("File 为空")
	}
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func A() {
	fmt.Println(File)
}
