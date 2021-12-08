package config

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Files struct {
		Name string `yaml:"name"`
	} `yaml:"files"`
	ExternalApis struct {
		Url string `yaml:"url"`
	} `yaml:"externalApis"`
}
