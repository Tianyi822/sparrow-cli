package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sparrow-cli/client"
	"sparrow-cli/config"
	"sparrow-cli/env"
	"sparrow-cli/global"
	"sparrow-cli/logger"
	"time"
)

func initProjEnv() {
	// 判断环境变量是否有 SparrowCliHome
	//	- 若有，则从 SparrowCliHome 中加载配置文件并将该路径保存到全局变量 env.SparrowCliHome 中
	// 	- 若没有，则指定默认路径 ~/.sparrow-cli 为 HOME_PATH，并保存在 env.SparrowCliHome 中
	homePath := os.Getenv("SparrowCliHome")
	if homePath == "" {
		homePath = os.Getenv("HOME") + "/.sparrow-cli"
	}
	env.SparrowCliHome = homePath
}

// initComponents 初始化组件
func initComponents(ctx context.Context) {
	// 初始化日志组件
	if err := logger.InitLogger(ctx); err != nil {
		log.Fatalf("初始化日志组件失败: %v", err)
	}
}

func main() {
	// 初始化项目家目录
	initProjEnv()

	// 加载配置文件
	config.LoadConfig()

	// 加载组件
	initializationCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	initComponents(initializationCtx)
	cancel()

	// 构建请求体
	requestBody := client.RequestBody{
		Model: global.CurrentModel.Name,
		Messages: []client.Message{
			{
				Role:    "system",
				Content: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。",
			},
			{
				Role:    "user",
				Content: "你好，我叫李雷，1+1等于多少？",
			},
		},
		Temperature: 0.6,
	}

	// 将请求体序列化为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("JSON编码失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", global.CurrentModel.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+global.CurrentModel.ApiKey)

	// 创建HTTP客户端
	c := &http.Client{}

	// 发送请求
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			fmt.Println(closeErr)
		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}

	// 打印响应结果
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", body)
}
