package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

type Config struct {
	GroqApiKey   string `json:"groq_api_key"`
	Model        string `json:"model"`
	TerminalType string `json:"terminal_type"`
}

func main() {
	// Loading config
	var config Config
	// Check if config.json exists
	config, err := loadConfig(config)
	if err != nil {
		//log.Println("Error loading config:", err)
		return
	}

	// Define flags for the program
	model := flag.String("model", "gemma2-9b-it", "Set the model to use for the request, check https://console.groq.com/docs/models for more information")
	groqApiKey := flag.String("api_key", "", "Set the Groq API key to use for the request, check https://console.groq.com/keys for more information")
	terminalType := flag.String("terminal", "cmd", "Set the terminal type you're using, such as cmd, bash, powershell, etc.")
	help := flag.Bool("help", false, "Show help for the program")

	flag.Parse()
	if *help {
		getHelp()
		return
	}

	// Checking if the flags are provided
	argsString := strings.Join(os.Args[1:], " ")

	if strings.Contains(argsString, "-apiKey") {
		fmt.Println("API key changed from", config.GroqApiKey, "to", *groqApiKey)
		config.GroqApiKey = *groqApiKey
	}

	if strings.Contains(argsString, "-model") {
		fmt.Println("Model changed from", config.Model, "to", *model)
		config.Model = *model
	}

	if strings.Contains(argsString, "-terminal") {
		if *terminalType == "cmd" || *terminalType == "bash" || *terminalType == "powershell" || *terminalType == "zsh" || *terminalType == "fish" {
			fmt.Println("Terminal type changed from", config.TerminalType, "to", *terminalType)
			config.TerminalType = *terminalType
		} else {
			fmt.Println("Error: invalid terminal type!")
			return
		}
	}

	// Getting the API key from the flag, environment variable or json file
	if *groqApiKey == "" {
		if config.GroqApiKey == "" {
			if os.Getenv("GROQ_API_KEY") == "" {
				fmt.Println("Error: groq-api-key is required, check https://console.groq.com/keys for more information\n. And use the -api flag to set the GROQ API key.")
				return
			} else {
				config.GroqApiKey = os.Getenv("GROQ_API_KEY")
			}
		}
	} else {
		config.GroqApiKey = *groqApiKey
	}

	// Check if no flags were provided and only a prompt is given
	if flag.NFlag() == 0 && len(os.Args) == 2 {

		if len(os.Args) != 2 {
			fmt.Println("Only one argument is required, such as: ask \"List all directories?\"")
			return
		}

		// Parse the prompt from user
		prompt := os.Args[1]

		// Make the API request
		response, err := makeRequest(prompt, *model, *terminalType)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Copy the response to clipboard
		if err := clipboard.WriteAll(response); err != nil {
			fmt.Printf("Error copying to clipboard: %v\n", err)
			return
		}

		fmt.Println("Response copied to clipboard: ", response)
	} else if flag.NFlag() == 0 && len(os.Args) != 2 {
		fmt.Println("Error: no such command!")
		getHelp()
		return
	}

	// Save the config to ask_config.json
	saveConfig(config)
}

// go build -o ask.exe .
// ask "Create a new directory hello and remove it right after?"
