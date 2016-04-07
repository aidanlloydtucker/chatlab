package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/billybobjoeaglt/chatlab/common"
)

var config *Config
var configPath string

type Config struct {
	Username        string
	PrivateKey      string
	AnsweredStorePK bool
	Password        string
	ShouldSavePass  bool
	FirstTime       bool
}

var Password string

// Creates a new config struct
func MakeConfig() (*Config, error) {
	var newConfig = &Config{
		Username:        "",
		PrivateKey:      "./key.key",
		AnsweredStorePK: false,
		Password:        "",
		ShouldSavePass:  true,
		FirstTime:       true,
	}
	return newConfig, nil
}

func GetConfig() *Config {
	return config
}

// Checks if config is in file. If not, make one; If so, open the file.
func LoadConfig() error {
	configPath = filepath.Join(common.ProgramDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		var newConfig *Config
		newConfig, err = MakeConfig()
		if err != nil {
			return err
		}
		config = newConfig
	} else {
		configFile, err := os.Open(configPath)
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

	Password = config.Password

	return nil
}

func SaveConfig() error {
	file, err := json.Marshal(&config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, file, 0777)
	return err
}
