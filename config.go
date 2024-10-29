// config.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	RPCUser     string `json:"rpcuser"`
	RPCPassword string `json:"rpcpassword"`
	RPCHost     string `json:"rpchost"`
	RPCPort     string `json:"rpcport"`
}

func LoadConfig() (*Config, error) {
	var config Config

	// First, try to load from environment variables
	config.RPCUser = os.Getenv("RPCUSER")
	config.RPCPassword = os.Getenv("RPCPASSWORD")
	config.RPCHost = os.Getenv("RPCHOST")
	config.RPCPort = os.Getenv("RPCPORT")

	if config.RPCUser != "" && config.RPCPassword != "" {
		// Use config from environment variables
		if config.RPCHost == "" {
			config.RPCHost = "localhost"
		}
		if config.RPCPort == "" {
			config.RPCPort = "8332"
		}
		return &config, nil
	}

	// If not found, try to load from config.json
	file, err := os.Open("config.json")
	if err != nil {
		return nil, fmt.Errorf("Error opening config file: %v", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config file: %v", err)
	}

	return &config, nil
}
