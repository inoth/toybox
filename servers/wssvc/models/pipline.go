package models

type HubConfig struct {
	MaxMsgChan int      `yaml:"MaxMsgChan"`
	Parser     string   `yaml:"Parser"`  // 解析管道仅限一个
	Process    []string `yaml:"Process"` // 处理管道可配置多个
}
