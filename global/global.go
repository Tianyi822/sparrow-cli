package global

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// =============================================================================
// 模型配置管理
// =============================================================================

// Model 全局模型配置
type Model struct {
	Name   string // 模型名称
	ApiKey string // API密钥
	URL    string // API地址
}

// CurrentModel 当前使用的模型
var CurrentModel *Model

// SetCurrentModel 设置当前模型
func SetCurrentModel(name, apiKey, url string) {
	CurrentModel = &Model{
		Name:   name,
		ApiKey: apiKey,
		URL:    url,
	}
}

// =============================================================================
// 系统提示词管理
// =============================================================================

// SystemPrompt 全局系统提示词
var SystemPrompt string

// DefaultSystemPrompt 默认系统提示词
const DefaultSystemPrompt = "你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。专有名词不可翻译成其他语言。"

// InitSystemPrompt 初始化系统提示词，提示用户输入或使用默认值
func InitSystemPrompt() {
	fmt.Println("正在初始化系统提示词...")
	fmt.Printf("默认系统提示词: %s\n\n", DefaultSystemPrompt)

	fmt.Print("请输入自定义系统提示词（直接回车使用默认值）: ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			SystemPrompt = input
			fmt.Printf("✓ 已设置自定义系统提示词\n")
		} else {
			SystemPrompt = DefaultSystemPrompt
			fmt.Printf("✓ 使用默认系统提示词\n")
		}
	} else {
		// 如果读取失败，使用默认值
		SystemPrompt = DefaultSystemPrompt
		fmt.Printf("✓ 使用默认系统提示词\n")
	}

	fmt.Println("系统提示词初始化完成")
}

// GetSystemPrompt 获取当前系统提示词
func GetSystemPrompt() string {
	if SystemPrompt == "" {
		return DefaultSystemPrompt
	}
	return SystemPrompt
}

// SetSystemPrompt 设置系统提示词
func SetSystemPrompt(prompt string) {
	if prompt == "" {
		SystemPrompt = DefaultSystemPrompt
	} else {
		SystemPrompt = prompt
	}
}
