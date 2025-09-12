package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sparrow-cli/logger"
	"strings"
)

// ResponseBody AI API 响应的主体结构
type ResponseBody struct {
	ID      string   `json:"id"`      // 请求的唯一标识符
	Object  string   `json:"object"`  // 响应对象类型，通常为 "chat.completion"
	Created int64    `json:"created"` // 响应创建的时间戳（Unix 时间戳）
	Model   string   `json:"model"`   // 使用的AI模型名称
	Choices []Choice `json:"choices"` // 响应选择列表，包含AI生成的回复
	Usage   Usage    `json:"usage"`   // Token使用情况统计
}

// Choice 单个响应选择项
type Choice struct {
	Index        int     `json:"index"`         // 选择项的索引位置
	Message      Message `json:"message"`       // AI助手返回的消息内容
	FinishReason string  `json:"finish_reason"` // 响应结束的原因（如 "stop", "length", "content_filter" 等）
}

// Usage Token使用情况统计信息
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 输入提示消耗的token数量
	CompletionTokens int `json:"completion_tokens"` // 生成回复消耗的token数量
	TotalTokens      int `json:"total_tokens"`      // 总共消耗的token数量（输入+输出）
}

// StreamChunk 流式响应的单个数据块结构
type StreamChunk struct {
	ID      string              `json:"id"`      // 请求的唯一标识符
	Object  string              `json:"object"`  // 响应对象类型，通常为 "chat.completion.chunk"
	Created int64               `json:"created"` // 响应创建的时间戳（Unix 时间戳）
	Model   string              `json:"model"`   // 使用的AI模型名称
	Choices []StreamChunkChoice `json:"choices"` // 流式响应选择列表
}

// StreamChunkChoice 流式响应中的单个选择项
type StreamChunkChoice struct {
	Index        int              `json:"index"`         // 选择项的索引位置
	Delta        StreamChunkDelta `json:"delta"`         // 增量数据
	FinishReason *string          `json:"finish_reason"` // 响应结束的原因（可为 null）
	Usage        *Usage           `json:"usage"`         // Token使用情况（仅在最后一个块中显示）
}

// StreamChunkDelta 流式响应中的增量数据
type StreamChunkDelta struct {
	Role    string `json:"role,omitempty"`    // 消息发送者角色（只在第一个块中显示）
	Content string `json:"content,omitempty"` // 增量消息内容
}

// ParseResponse 解析 HTTP 响应并返回 ResponseBody 结构体
// 参数:
//   - resp: HTTP 响应对象
//
// 返回:
//   - *ResponseBody: 解析后的响应数据结构
//   - error: 解析过程中的错误
func ParseResponse(resp *http.Response) (*ResponseBody, error) {
	// 确保响应体在函数结束时关闭
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Warn("关闭响应体失败: %v", closeErr)
		}
	}()

	// 读取响应体数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 数据到 ResponseBody 结构体
	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// ParseStreamResponse 解析流式 HTTP 响应并返回完整的 ResponseBody 结构体
// 参数:
//   - resp: HTTP 响应对象（text/event-stream 格式）
//
// 返回:
//   - *ResponseBody: 拼接后的完整响应数据结构
//   - error: 解析过程中的错误
func ParseStreamResponse(resp *http.Response) (*ResponseBody, error) {
	// 确保响应体在函数结束时关闭
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Warn("关闭响应体失败: %v", closeErr)
		}
	}()

	// 初始化结果结构体
	result := &ResponseBody{
		Choices: make([]Choice, 1), // 初始化一个选择项
	}
	result.Choices[0].Message.Role = AssistantRole

	// 创建扫描器按行读取
	scanner := bufio.NewScanner(resp.Body)
	var contentBuilder strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		// 检查是否是结束标志
		if line == "data: [DONE]" {
			break
		}

		// 解析 data: 开头的行
		if strings.HasPrefix(line, "data: ") {
			jsonData := line[6:] // 移除 "data: " 前缀

			// 解析 JSON 数据块
			var chunk StreamChunk
			if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
				logger.Warn("解析流式数据块失败: %v, 数据: %s", err, jsonData)
				continue
			}

			// 填充基本信息（只在第一次时填充）
			if result.ID == "" {
				result.ID = chunk.ID
				result.Object = "chat.completion" // 转换为非流式的对象类型
				result.Created = chunk.Created
				result.Model = chunk.Model
			}

			// 处理选择项
			if len(chunk.Choices) > 0 {
				choice := chunk.Choices[0]

				// 拼接内容
				if choice.Delta.Content != "" {
					contentBuilder.WriteString(choice.Delta.Content)
				}

				// 检查结束原因
				if choice.FinishReason != nil {
					result.Choices[0].FinishReason = *choice.FinishReason
				}

				// 获取 Token 使用情况（通常在最后一个块中）
				if choice.Usage != nil {
					result.Usage = *choice.Usage
				}
			}
		}
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取流式响应失败: %v", err)
	}

	// 设置最终内容
	result.Choices[0].Message.Content = contentBuilder.String()
	result.Choices[0].Index = 0

	return result, nil
}

// ParseStreamResponseWithCallback 解析流式 HTTP 响应并在每个数据块到达时调用回调函数
// 参数:
//   - resp: HTTP 响应对象（text/event-stream 格式）
//   - callback: 每个数据块的回调函数（参数: 增量内容, 是否结束）
//
// 返回:
//   - *ResponseBody: 拼接后的完整响应数据结构
//   - error: 解析过程中的错误
func ParseStreamResponseWithCallback(resp *http.Response, callback func(content string, isFinished bool)) (*ResponseBody, error) {
	// 确保响应体在函数结束时关闭
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.Warn("关闭响应体失败: %v", closeErr)
		}
	}()

	// 初始化结果结构体
	result := &ResponseBody{
		Choices: make([]Choice, 1), // 初始化一个选择项
	}
	result.Choices[0].Message.Role = AssistantRole

	// 创建扫描器按行读取
	scanner := bufio.NewScanner(resp.Body)
	var contentBuilder strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		// 检查是否是结束标志
		if line == "data: [DONE]" {
			if callback != nil {
				callback("", true) // 通知结束
			}
			break
		}

		// 解析 data: 开头的行
		if strings.HasPrefix(line, "data: ") {
			jsonData := line[6:] // 移除 "data: " 前缀

			// 解析 JSON 数据块
			var chunk StreamChunk
			if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
				logger.Warn("解析流式数据块失败: %v, 数据: %s", err, jsonData)
				continue
			}

			// 填充基本信息（只在第一次时填充）
			if result.ID == "" {
				result.ID = chunk.ID
				result.Object = "chat.completion" // 转换为非流式的对象类型
				result.Created = chunk.Created
				result.Model = chunk.Model
			}

			// 处理选择项
			if len(chunk.Choices) > 0 {
				choice := chunk.Choices[0]

				// 拼接内容并调用回调
				if choice.Delta.Content != "" {
					contentBuilder.WriteString(choice.Delta.Content)
					if callback != nil {
						callback(choice.Delta.Content, false)
					}
				}

				// 检查结束原因
				if choice.FinishReason != nil {
					result.Choices[0].FinishReason = *choice.FinishReason
				}

				// 获取 Token 使用情况（通常在最后一个块中）
				if choice.Usage != nil {
					result.Usage = *choice.Usage
				}
			}
		}
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取流式响应失败: %v", err)
	}

	// 设置最终内容
	result.Choices[0].Message.Content = contentBuilder.String()
	result.Choices[0].Index = 0

	return result, nil
}
