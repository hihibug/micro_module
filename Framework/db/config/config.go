package config

type Config struct {
	DbType      string `json:"db-type" yaml:"dbType"`
	Path        string `json:"path" yaml:"path"`
	Config      string `json:"config" yaml:"config"`
	Dbname      string `json:"dbname" yaml:"dbName"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	MaxIdleCons int    `json:"maxIdleCons" yaml:"maxIdleCons"`
	MaxOpenCons int    `json:"maxOpenCons" yaml:"maxOpenCons"`
	LogMode     string `json:"logMode" yaml:"logMode"`
}
