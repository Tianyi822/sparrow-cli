package config

import (
	"fmt"
	"os"
	"sparrow-cli/env"
	"sparrow-cli/file"
	"sync"

	"gopkg.in/yaml.v3"
)

var loadConfigOnce sync.Once

var (
	Models []ModelConfig
	Logger LoggerConfigData
)

func LoadConfig() {
	loadConfigOnce.Do(func() {
		// 2. 判断 SparrowCliHome 是否有 config.yaml 文件
		// 	2.1 若有，则加载该文件
		// 	2.2 若没有，则按照 config.items 结构创建 config.yaml 文件并保存在 HOME_PATH 中
		configFilePath := env.SparrowCliHome + "/config/sparrow_cli_config.yaml"
		if !file.IsExist(configFilePath) {
			// 新建文件并保存空配置
			f, createErr := file.CreateFile(configFilePath)
			if createErr != nil {
				panic(createErr)
			}
			pc := &ProjectConfig{}
			yamlData, err := yaml.Marshal(pc)
			if err != nil {
				panic(fmt.Errorf("将配置数据转换为 YAML 失败: %w", err))
			}
			_, wErr := f.Write(yamlData)
			if wErr != nil {
				panic(fmt.Errorf("写入 YAML 数据到文件失败: %w", wErr))
			}
			closeErr := f.Close()
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
		Models = conf.Models
		Logger = conf.Logger

		// 设置环境中的默认模型
		if len(Models) > 0 {
			env.SetCurrentModel(Models[0].Model, Models[0].ApiKey, Models[0].URL)
		}
	})
}
