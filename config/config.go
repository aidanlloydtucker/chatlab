package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var config *Config

type Config struct {
	Username   string
	PrivateKey string
	Passphrase string
	FirstTime  bool
}

// Creates a new config struct
func MakeConfig() (*Config, error) {
	var newConfig = &Config{
		Username:   "",
		PrivateKey: "./key.key",
		Passphrase: "./pass.key",
		FirstTime:  true,
	}
	return newConfig, nil
}

func GetConfig() *Config {
	return config
}

// Checks if config is in file. If not, make one; If so, open the file.
func LoadConfig() error {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		var newConfig *Config
		newConfig, err = MakeConfig()
		if err != nil {
			return err
		}
		config = newConfig
	} else {
		configFile, err := os.Open("./config.json")
		if err != nil {
			return err
		}
		defer configFile.Close()
		byteArrFile, err := ioutil.ReadAll(configFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(byteArrFile, &config)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveConfig() error {
	file, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./config.json", file, 0777)
	return err
}
