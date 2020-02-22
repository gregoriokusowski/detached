package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	CONFIG_FOLDER = ".detached/default"
	CONFIG_FILE   = "config"
)

func Exists() bool {
	_, err := os.Stat(absConfigPath())
	return err == nil
}

func Save(instance interface{}) error {
	if _, err := os.Stat(AbsConfigFolder()); os.IsNotExist(err) {
		err = os.MkdirAll(AbsConfigFolder(), 0755)
		if err != nil {
			return fmt.Errorf("Failed to create ~/.detached folder: %s", err)
		}
	}
	bytes, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to convert config to json: %s", err)
	}

	err = ioutil.WriteFile(absConfigPath(), bytes, 0755)
	if err != nil {
		return fmt.Errorf("Failed to persist config: %s", err)
	}
	return nil
}

func Load(c interface{}) error {
	raw, err := ioutil.ReadFile(absConfigPath())
	if err != nil {
		return fmt.Errorf("Failed to load config: %s", err)
	}

	err = json.Unmarshal(raw, &c)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %s", err)
	}

	return nil
}

func AbsConfigFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, CONFIG_FOLDER)
}

func absConfigPath() string {
	return filepath.Join(AbsConfigFolder(), CONFIG_FILE)
}

// Creates a file with the config file
func AddConfig(filename, content string) error {
	err := ioutil.WriteFile(filepath.Join(AbsConfigFolder(), filename), []byte(content), 0755)
	if err != nil {
		return fmt.Errorf("Failed to persist %s config: %s", filename, err)
	}
	return nil
}

// Retrieves a config file content
func GetConfig(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filepath.Join(AbsConfigFolder(), filename))
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}
