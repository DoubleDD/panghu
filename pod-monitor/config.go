package config

type Config struct {
	App struct {
		name string `yaml:"name"`
	} `yaml:"app"`

	Minio struct {
		accessKey string `yaml:"accessKey"`
		secretKey string `yaml:"secretKey"`
		bucket    string `yaml:"bucket"`
		endpoint  string `yaml:"endpoint"`
		prefix    string `yaml:"prefix"`
	} `yaml:"minio"`
}
