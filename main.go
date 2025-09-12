package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sparrow-cli/client"
	"sparrow-cli/config"
	"sparrow-cli/env"
	"sparrow-cli/global"
	"sparrow-cli/logger"
	"strings"
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

	// 初始化系统提示词
	global.InitSystemPrompt()

	// 构建请求体
	var messages []client.Message

	// 使用全局管理的系统消息
	messages = append(messages, client.Message{
		Role:    client.SysRole,
		Content: global.GetSystemPrompt(),
	})

	// 创建标准输入扫描器
	scanner := bufio.NewScanner(os.Stdin)

	// 创建HTTP客户端
	// 9.9 和 9.11 哪个大，这个问题为什么通常用来测试大模型
	c := &http.Client{}

	for {
		fmt.Print("请输入问题：")
		// 用户输入的问题
		if !scanner.Scan() {
			break
		}
		msg := strings.TrimSpace(scanner.Text())
		if msg == "" {
			continue
		}
		if msg == "!quit" {
			break
		}
		messages = append(messages, client.Message{
			Role:    client.UserRole,
			Content: msg,
		})

		req := client.BuildStreamRequest(messages, 0.6)

		// 发送请求
		resp, err := c.Do(req)
		if err != nil {
			logger.Fatal("请求失败: %v", err)
		}

		// 解析响应数据
		responseBody, err := client.ParseStreamResponseWithCallback(resp, printContent)
		if err != nil {
			logger.Fatal("解析响应失败: %v", err)
		}

		// 打印响应结果
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("模型: %s\n", responseBody.Model)

		fmt.Printf("Token使用: 输入=%d, 输出=%d, 总计=%d\n",
			responseBody.Usage.PromptTokens,
			responseBody.Usage.CompletionTokens,
			responseBody.Usage.TotalTokens)

		// 将AI的回复添加到对话历史中
		if len(responseBody.Choices) > 0 {
			messages = append(messages, client.Message{
				Role:    client.AssistantRole,
				Content: responseBody.Choices[0].Message.Content,
			})
		}
	}
}

func printContent(content string, isFinished bool) {
	fmt.Print(content)
}
