package configuration

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

// Configuration struct
type Configuration struct {
	Service   string `yaml:"SERVICE"`
	Model     string `yaml:"MODEL"`
	GeminiKey string `yaml:"GEMINI_KEY"`
	OpenAIKey string `yaml:"OPENAI_KEY"`
}

var defaultConfig = Configuration{
	Service:   "",
	Model:     "",
	GeminiKey: "",
	OpenAIKey: "",
}

func getConfigFilePath() string {
	var configDir string
	var configFileName string

	switch OS := runtime.GOOS; OS {
	case "windows":
		// On Windows, use the APPDATA directory
		configDir = filepath.Join(os.Getenv("APPDATA"), "TerminalAI")
		configFileName = "config.yaml"
	case "darwin":
		// On macOS, use the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configDir = filepath.Join(home, "Library", "Application Support", "TerminalAI")
		configFileName = "config.yaml"
	default:
		// On Unix-like systems, use the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configDir = filepath.Join(home, ".terminalAI")
		configFileName = "config.yaml"
	}

	// Create configuration directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	return filepath.Join(configDir, configFileName)
}

func readConfig(filePath string) Configuration {
	// Read configuration from file
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		// If configuration file doesn't exist or cannot be read, return default configuration
		fmt.Println("\x1b[31mConfiguration file not found or cannot be read.\x1b[37m")
		fmt.Printf("Please run '\x1b[32mterminalai config init\x1b[37m' to set up the configuration.\n\n")

		initConfig()
	}

	// Unmarshal configuration from YAML
	var config Configuration
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		fmt.Println("\x1b[31mError unmarshalling configuration.\x1b[37m")
		fmt.Println("Please run '\x1b[32mterminalai config init\x1b[37m' to reset the configuration.")
		// panic(err)
		os.Exit(1)
	}

	return config
}

func saveConfig(config Configuration, filePath string) {
	// Marshal configuration to YAML
	configData, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	// Write configuration to file
	if err := os.WriteFile(filePath, configData, 0644); err != nil {
		fmt.Println("\x1b[31mError saving configuration.\x1b[37m")
		panic(err)
	}

	fmt.Println("Configuration saved successfully.")
}

func Config() {

	if len(os.Args) < 3 {
		// Usage
		fmt.Println("Usage: terminalai config <command>")
		fmt.Println("Commands:")
		fmt.Println("  show\t\t\t\tShow current configuration")

		fmt.Println("  set\t <key> <value>\t\tSet configuration")
		fmt.Println("  init\t\t\t\tInitialize configuration")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "show":
		showConfig()
	case "set":
		if len(os.Args) < 4 {
			// Usage
			fmt.Println("Usage: terminalai config set <key> <value>")
			os.Exit(1)
		}
		if len(os.Args) == 4 {
			// fmt.Println("Value cannot be empty.")
			setConfig(os.Args[3], "")

		} else {
			setConfig(os.Args[3], os.Args[4])
		}

	// case "reset":
	// 	resetConfig()

	case "init":
		initConfig()
	}

}

func GetConfig(key string) string {
	configFile := getConfigFilePath()
	config := readConfig(configFile)

	// service and model cannot be empty
	if config.Service == "" || config.Model == "" {
		fmt.Println("\x1b[31mService and Model configuration cannot be empty.\x1b[37m")
		fmt.Println("Please run '\x1b[32mterminalai config set <key> <value>\x1b[37m' to set the configuration.")
		os.Exit(1)
	}

	switch key {
	case "service":
		// if empty, return ""
		return config.Service
	case "model":
		return config.Model
	case "gemini-key":
		return config.GeminiKey
	case "openai-key":
		return config.OpenAIKey
	}

	return ""
}

func showConfig() {
	configFile := getConfigFilePath()
	config := readConfig(configFile)

	// print whole configuration with key value pairs
	fmt.Println("Sevice:", config.Service)
	fmt.Println("model:", config.Model)
	fmt.Println("Gemini API Key:", config.GeminiKey)
	fmt.Println("OpenAI API Key:", config.OpenAIKey)
	os.Exit(0)
}

func setConfig(key string, value string) {
	configFile := getConfigFilePath()
	config := readConfig(configFile)

	switch key {
	case "service":
		config.Service = value
	case "model":
		config.Model = value
	case "gemini-key":
		config.GeminiKey = value
	case "openai-key":
		config.OpenAIKey = value
	}

	saveConfig(config, configFile)

	fmt.Println("Configuration updated successfully.")
	os.Exit(0)

}

func initConfig() {

	reader := bufio.NewReader(os.Stdin)

	configFile := getConfigFilePath()
	config := defaultConfig

	fmt.Println("Welcome to Terminal AI Configuration Setup")
	fmt.Println("Please set the following configuration options:")
	fmt.Printf("1. Terminal AI Service (e.g., openai, gemini):\x1b[32m")
	config.Service = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())
	fmt.Printf("\x1b[37m2. Terminal AI Model (e.g., gpt-3, gpt-4-turbo, etc.):\x1b[32m")
	config.Model = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())
	fmt.Printf("\x1b[37m3. Google API Key (optional):\x1b[32m")
	config.GeminiKey = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())
	fmt.Printf("\x1b[37m4. OpenAI API Key (optional):\x1b[32m")
	config.OpenAIKey = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())

	fmt.Printf("\x1b[37m")
	saveConfig(config, configFile)
	fmt.Println("Configuration setup completed successfully.")
	os.Exit(0)

}
