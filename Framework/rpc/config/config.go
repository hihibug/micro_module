package config

type Config struct {
	Ip          string
	Port        uint64
	ServiceName string
	Weight      uint64
	Metadata    map[string]string
	ClusterName string // 默认值DEFAULT
	GroupName   string // 默认值DEFAULT_GROUP
}
