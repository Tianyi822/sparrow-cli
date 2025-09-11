package config

// ProjectConfig 项目配置
type ProjectConfig struct {
	Models []ModelConfig `yaml:"models"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	Model  string `yaml:"model"`
	ApiKey string `yaml:"api_key"`
	URL    string `yaml:"url"`
}
