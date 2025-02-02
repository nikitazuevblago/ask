package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func getHelp() {
	fmt.Println("Ask is a CLI helper with Groq API under the hood")
	fmt.Println("Usage: ask [options] OR ask \"prompt\"")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func getConfigPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(execPath), "ask_config.json"), nil
}

func loadConfig(config Config) (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		//log.Println("File does not exist")
		// If the file does not exist, create a default config
		config = Config{"", "gemma2-9b-it", "cmd"}
		// Save the default config immediately
		if jsonData, err := json.Marshal(config); err == nil {
			os.WriteFile(configPath, jsonData, 0644)
		}
	} else {
		//log.Println("File exists:", file.Name())
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return Config{}, err
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			// If we can't decode the config, create a new default one
			config = Config{"", "gemma2-9b-it", "cmd"}
			if jsonData, err := json.Marshal(config); err == nil {
				os.WriteFile(configPath, jsonData, 0644)
			}
		}
	}
	return config, nil
}

func saveConfig(config Config) {
	configPath, err := getConfigPath()
	if err != nil {
		fmt.Println("Error getting config path:", err)
		return
	}

	json, err := json.Marshal(config)
	if err != nil {
		//log.Println("Error saving config:", err)
		return
	}
	os.WriteFile(configPath, json, 0644)
}
