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
	bytes, err := json.Marshal(instance)
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
