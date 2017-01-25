package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Server serverConfig
}

type serverConfig struct {
	Port int    `json:"port"`
	Name string `json:"name"`
}

func LoadConfig(path string) Configuration {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}

	var config Configuration
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}
	return config
}
