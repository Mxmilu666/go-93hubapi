package source

import (
	"os"

	"gopkg.in/yaml.v2"
)

type GitHubConfig struct {
	Token string `yaml:"token"`
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Config struct {
	GitHub        GitHubConfig `yaml:"github"`
	Server        ServerConfig `yaml:"server"`
	Dest          string       `yaml:"dest"`
	ProxyURL      string       `yaml:"ProxyURL"`
	MaxConcurrent int          `yaml:"maxconcurrent"`
}

func CreateDefaultConfig(filename string) error {
	defaultConfig := Config{
		GitHub: GitHubConfig{
			Token: "your_token",
			Owner: "Mxmilu666",
			Repo:  "Bangbang93Hub",
		},
		Server: ServerConfig{
			Address: "localhost",
			Port:    8080,
		},
		Dest:          "your_dest",
		ProxyURL:      "",
		MaxConcurrent: 4,
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	err = encoder.Encode(defaultConfig)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
