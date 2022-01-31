package utils

type section map[string]string

//读写配置信息
type Config struct {
	Filename string `ini:"filename"`
	Method   string `ini:"method"`
}

var defaultConfig *Config
