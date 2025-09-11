package config

import (
	"fmt"
	"os"
	"sparrow-cli/env"
	"sparrow-cli/fileutils"
	"sync"

	"gopkg.in/yaml.v3"
)

var loadConfigOnce sync.Once

var (
	models []ModelConfig
)

func LoadConfig() {
	loadConfigOnce.Do(func() {
		// 1. 判断环境变量是否有 SPARROW_CLI_HOME
		//	1.1 若有，则从 SPARROW_CLI_HOME 中加载配置文件并将该路径保存到全局变量 env.SPARROW_CLI_HOME 中
		// 	1.2 若没有，则指定默认路径 ~/.sparrow-cli 为 HOME_PATH，并保存在 env.SPARROW_CLI_HOME 中
		homePath := os.Getenv("SPARROW_CLI_HOME")
		if homePath == "" {
			homePath = os.Getenv("HOME") + "/.sparrow-cli"
		}
		env.SPARROW_CLI_HOME = homePath

		// 2. 判断 SPARROW_CLI_HOME 是否有 config.yaml 文件
		// 	2.1 若有，则加载该文件
		// 	2.2 若没有，则按照 config.items 结构创建 config.yaml 文件并保存在 HOME_PATH 中
		configFilePath := env.SPARROW_CLI_HOME + "/config/sparrow_cli_config.yaml"
		if !fileutils.IsExist(configFilePath) {
			// 新建文件并保存空配置
			file, createErr := fileutils.CreateFile(configFilePath)
			if createErr != nil {
				panic(createErr)
			}
			pc := &ProjectConfig{}
			yamlData, err := yaml.Marshal(pc)
			if err != nil {
				panic(fmt.Errorf("将配置数据转换为 YAML 失败: %w", err))
			}
			_, wErr := file.Write(yamlData)
			if wErr != nil {
				panic(fmt.Errorf("写入 YAML 数据到文件失败: %w", wErr))
			}
			closeErr := file.Close()
			if closeErr != nil {
				panic(fmt.Errorf("关闭文件失败: %w", closeErr))
			}
		}

		// 3. 从配置文件中加载配置
		// 读取配置文件内容
		data, readErr := os.ReadFile(configFilePath)
		if readErr != nil {
			// 如果读取配置文件时发生错误，返回错误信息
			panic(fmt.Errorf("load config file error: %w", readErr))
		}
		// 解析配置数据
		conf := &ProjectConfig{}
		if unErr := yaml.Unmarshal(data, &conf); unErr != nil {
			// 如果解析配置到结构体时发生错误，返回错误信息
			panic(fmt.Errorf("reflect config to struct error: %w", unErr))
		}

		// 设置全局配置
		models = conf.Models
	})
}
