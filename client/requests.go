package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sparrow-cli/global"
	"sparrow-cli/logger"
)

// Role AI 对话中的角色类型定义
type Role string

const (
	SysRole       Role = "system"    // 系统角色，用于设置 AI 助手的行为和指令
	UserRole      Role = "user"      // 用户角色，表示来自用户的消息
	AssistantRole Role = "assistant" // 助手角色，表示 AI 助手的回复消息
)

// RequestBody AI API 请求体结构
type RequestBody struct {
	Model       string    `json:"model"`       // 使用的AI模型名称
	Messages    []Message `json:"messages"`    // 对话消息列表
	Temperature float64   `json:"temperature"` // 生成文本的随机性控制参数（0.0-2.0）
	Stream      bool      `json:"stream"`      // 是否启用流式响应
}

// Message 单条对话消息结构
type Message struct {
	Role    Role   `json:"role"`    // 消息发送者角色
	Content string `json:"content"` // 消息内容文本
}

// BuildRequest 构建 AI API 的 HTTP 请求（向后兼容，默认非流式）
// 参数:
//   - messages: 对话消息列表
//   - temperature: 生成文本的随机性控制参数
//
// 返回:
//   - *http.Request: 构建完成的 HTTP 请求对象
//
// 已弃用: 推荐使用 BuildNonStreamRequest 或 BuildStreamRequest
func BuildRequest(messages []Message, temperature float64) *http.Request {
	// 构建请求体数据（非流式）
	reqBody := &RequestBody{
		Model:       global.CurrentModel.Name, // 从全局配置获取模型名称
		Messages:    messages,                 // 传入的对话消息
		Temperature: temperature,              // 设置温度参数
		Stream:      false,                    // 设置为非流式响应
	}

	return buildHTTPRequest(reqBody)
}

// BuildStreamRequest 构建流式 AI API 的 HTTP 请求
// 参数:
//   - messages: 对话消息列表
//   - temperature: 生成文本的随机性控制参数
//
// 返回:
//   - *http.Request: 构建完成的流式 HTTP 请求对象
func BuildStreamRequest(messages []Message, temperature float64) *http.Request {
	// 构建请求体数据（流式）
	reqBody := &RequestBody{
		Model:       global.CurrentModel.Name, // 从全局配置获取模型名称
		Messages:    messages,                 // 传入的对话消息
		Temperature: temperature,              // 设置温度参数
		Stream:      true,                     // 设置为流式响应
	}

	return buildHTTPRequest(reqBody)
}

// buildHTTPRequest 构建 HTTP 请求的内部方法
// 参数:
//   - reqBody: 请求体数据结构
//
// 返回:
//   - *http.Request: 构建完成的 HTTP 请求对象
func buildHTTPRequest(reqBody *RequestBody) *http.Request {
	// 将请求体序列化为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Fatal("JSON编码失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", global.CurrentModel.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Fatal("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+global.CurrentModel.ApiKey)

	return req
}
