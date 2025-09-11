package env

var SparrowCliHome = ""

// Model 环境中的模型配置（避免循环引用）
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
