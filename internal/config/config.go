package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	path := getConfigPath()
	configContents, err := os.ReadFile(path)
	if err != nil {
		errMsg := fmt.Sprintf("reading config file from %s not possible", path)
		log.Fatal(errMsg)
	}

	config := Config{}
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		errMsg := fmt.Sprintf("reading config file from %s not possible", path)
		log.Fatal(errMsg)
	}

	return config
}

func getConfigPath() string {
	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("no config file found in user home directory named .gatorconfig.json")
	}
	return path + "/" + configFileName
}

func (config *Config) SetUser(userName string) error {
	config.CurrentUserName = userName
	return write(config)
}

func write(config *Config) error {
	path := getConfigPath()

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// perm 0600 = owner can read and write
	return os.WriteFile(path, data, 0600)
}
