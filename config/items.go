package config

// ProjectConfig 项目配置
type ProjectConfig struct {
	Models []ModelConfig    `yaml:"models"`
	Logger LoggerConfigData `yaml:"logger"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	Model  string `yaml:"model"`
	ApiKey string `yaml:"api_key"`
	URL    string `yaml:"url"`
}

// LoggerConfigData 定义了日志配置
type LoggerConfigData struct {
	Level      string `yaml:"level"`       // 日志级别
	MaxAge     uint16 `yaml:"max_age"`     // 日志文件保留最大天数
	MaxSize    uint16 `yaml:"max_size"`    // 日志文件最大大小(MB)
	MaxBackups uint16 `yaml:"max_backups"` // 日志备份文件最大数量
	Compress   bool   `yaml:"compress"`    // 是否压缩日志文件
}
