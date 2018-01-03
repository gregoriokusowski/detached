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
	CONFIG_FOLDER = ".detached"
	CONFIG_FILE   = "default"
)

func Exists() bool {
	_, err := os.Stat(absConfigPath())
	return err == nil
}

func Save(c struct{}) error {

}

func Load(c *struct{}) error {
	raw, err := ioutil.ReadFile(absConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(CONFIG_FOLDER, 0644)
			if err != nil {
				return nil, err
			}
			instance, err := buildInstance(ctx)
			if err != nil {
				return nil, err
			}
			err = persist(instance)
			if err != nil {
				return nil, err
			}
			return instance, nil
		}
		fmt.Printf("%+w", err)
		return nil, err
	}
	var aws Aws
	json.Unmarshal(raw, &aws)
	return &aws, nil
}

func absConfigFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, CONFIG_FOLDER)
}

func absConfigPath() string {
	return filepath.Join(configFolder(), CONFIG_FILE)
}
