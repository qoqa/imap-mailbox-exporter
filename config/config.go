package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigAcccount struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ConfigServer struct {
	Host string `yaml:"hostname"`
	Port string `yaml:"port"`

	Account []ConfigAcccount `yaml:"accounts"`
}

func (configServer ConfigServer) HostPort() string {
	return configServer.Host + ":" + configServer.Port
}

type Config struct {
	Server []ConfigServer `yaml:"server"`
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := replaceEnvPlaceholder(string(configBytes))

	err = yaml.Unmarshal([]byte(configString), &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func replaceEnvPlaceholder(data string) string {
	expr := regexp.MustCompile("env:([A-Z_]+)")
	matches := expr.FindAll([]byte(data), -1)

	for _, match := range matches {
		variable := string(match)
		variable = strings.TrimLeft(variable, "env:")

		env := os.Getenv(variable)
		if env == "" {
			log.Printf("Environment variable %s is empty. Skipping replacement.", variable)
			continue
		}

		data = strings.ReplaceAll(data, fmt.Sprintf("env:%s", variable), env)
	}

	return data
}

// Find the account and server from the given hostname and username
func (config Config) FindAccountInServer(hostname, username string) (*ConfigServer, *ConfigAcccount, error) {
	for _, server := range config.Server {
		if server.Host != hostname {
			continue
		}

		for _, account := range server.Account {
			if account.Username != username {
				continue
			}

			return &server, &account, nil
		}
	}

	return nil, nil, errors.New("cound not find user on given server in configuration")
}
