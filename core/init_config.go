package core

import (
	"cfst/config"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const ConfigFile = "application.yaml"

func LoadYmlConfig() {
	conf := &config.Config{}
	yamlConf, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yamlConf error: %s", err))
	}
	err = yaml.Unmarshal(yamlConf, conf)
	if err != nil {
		log.Fatalf("config Init Unmarshal: %v", err)
	}
	log.Println("config yamlFile load Init success.")
	Config = conf
}
