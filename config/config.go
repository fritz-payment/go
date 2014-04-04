package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
)

// Config describes a config holder
type Config interface {
	ConfigFileName() string
	SetConfigFileName(string)
	IsCreated() bool
	SetIsCreated(bool)
	AppName() string
	DefaultFileName() string
}

// ReadConfig will read a JSON file into any structure
func ReadConfig(jsonFileName string, into interface{}) error {
	rv := reflect.ValueOf(into)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("cannot read into given value")
	}

	f, err := os.Open(jsonFileName)
	if err != nil {
		return fmt.Errorf("could not open JSON file %s: %v", jsonFileName, err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(into)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	return nil
}

// LoadConfig will load the configuration from the given config file name
// into the given cfg
//
// If the file does not exist, will attempt to create default config file.
// If the config file name is empty, will use default config file:
//   $HOME/.config/AppName.cfg.json
func LoadConfig(configFileName string, cfg Config) error {
	rv := reflect.ValueOf(cfg)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("config must be a pointer receiver")
	}
	var err error

	// create empty config
	if configFileName == "" {
		err = defaultConfigFileName(cfg)
		if err != nil {
			return fmt.Errorf("cannot determine default config file name: %v", err)
		}
	} else {
		cfg.SetConfigFileName(configFileName)
	}

	// check for existence
	_, err = os.Stat(cfg.ConfigFileName())
	if os.IsNotExist(err) {
		err = createDefaultConfig(cfg)
		if err != nil {
			return fmt.Errorf("cannot initialize empty config: %v", err)
		}
	}
	if err != nil {
		return fmt.Errorf("error on stat %s: %v", cfg.ConfigFileName(), err)
	}

	err = ReadConfig(cfg.ConfigFileName(), cfg)
	if err != nil {
		return fmt.Errorf("error on reading config: %v", err)
	}

	return nil
}

// Returns empty config object with default config file name set
func defaultConfigFileName(cfg Config) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot lookup current user: %v", err)
	}
	cfg.SetConfigFileName(filepath.Join(usr.HomeDir, ".config", cfg.AppName(), cfg.DefaultFileName()))
	return nil
}

// Create a default config at cfg.configFileName
//
// cfg will be populated with default values
func createDefaultConfig(cfg Config) error {
	cfgPath := filepath.Dir(cfg.ConfigFileName())
	err := os.MkdirAll(cfgPath, 0755)
	if err != nil {
		return fmt.Errorf("cannot create config dir %s: %v", cfg.ConfigFileName(), err)
	}
	cfgFile, err := os.OpenFile(cfg.ConfigFileName(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open config file %s for writing: %v", cfg.ConfigFileName(), err)
	}
	defer cfgFile.Close()
	jsonStr, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("cannot encode JSON: %v", err)
	}
	// JSON beautification
	buf := bytes.NewBuffer(nil)
	err = json.Indent(buf, jsonStr, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot indent JSON: %v", err)
	}
	_, err = io.Copy(cfgFile, buf)
	if err != nil {
		return fmt.Errorf("error writing buf to file: %v", err)
	}
	cfg.SetIsCreated(true)
	return nil
}
