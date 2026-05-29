package config

type Config struct {
	Database struct {
		Sqlite struct {
			Used bool   `yaml:"used"`
			Path string `yaml:"path"`
		} `yaml:"sqlite"`
		Postgres struct {
			Used     bool   `yaml:"used"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			DbName   string `yaml:"dbName"`
		} `yaml:"postgres"`
	} `yaml:"database"`
	Server struct {
		Port    string `yaml:"port"`
		Logging struct {
			Logger []string `yaml:"loggers"`
			Level  string   `yaml:"level"`
		} `yaml:"logging"`
	} `yaml:"server"`
}
