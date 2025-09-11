package global

// Model 全局模型配置
type Model struct {
	Name   string
	ApiKey string
	URL    string
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