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
	TerminalAIBackend string `yaml:"TERMINAL_AI_BACKEND"`
	TerminalAIModel   string `yaml:"TERMINAL_AI_MODEL"`
	GoogleAPIKey      string `yaml:"GOOGLE_AI_API_KEY"`
	OpenAPIKey        string `yaml:"OPENAI_API_KEY"`
	// Add more configuration options as needed
}

// func main() {
// 	configFile := getConfigFilePath()
// 	config := readConfig(configFile)

// 	// Use the configuration
// 	fmt.Println("Option1:", config.Option1)
// 	fmt.Println("Option2:", config.Option2)

// 	// Modify configuration
// 	config.Option2 = 42

// 	// Save configuration
// 	saveConfig(config, configFile)
// }

func getConfigFilePath() string {
	var configDir string
	var configFileName string

	switch OS := runtime.GOOS; OS {
	case "windows":
		// On Windows, use the APPDATA directory
		configDir = filepath.Join(os.Getenv("APPDATA"), "MyTool")
		configFileName = "config.yaml"
	case "darwin":
		// On macOS, use the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configDir = filepath.Join(home, "Library", "Application Support", "MyTool")
		configFileName = "config.yaml"
	default:
		// On Unix-like systems, use the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configDir = filepath.Join(home, ".mytool")
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
		fmt.Println("Configuration file not found or cannot be read. Using default configuration.")
		return Configuration{
			TerminalAIBackend: "",
			TerminalAIModel:   "",
			GoogleAPIKey:      "",
			OpenAPIKey:        "",
		}
	}

	// Unmarshal configuration from YAML
	var config Configuration
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		panic(err)
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
		panic(err)
	}

	fmt.Println("Configuration saved successfully.")
}

func Config() {

	// if os.Args[2] == "show" {
	// 	showConfig()

	// }
	// else if os.Args[]

	switch os.Args[2] {
	case "show":
		showConfig()
	case "set":
		if len(os.Args) < 4 {
			// Usage
			fmt.Println("Usage: terminalAI config set <key> <value>")
			os.Exit(1)
		}
		if len(os.Args) == 4 {
			os.Args[4] = ""
		}
		setConfig(os.Args[3], os.Args[4])

	case "reset":
		resetConfig()

	case "init":
		initConfig()
	}

}

func GetConfig(key string) string {
	configFile := getConfigFilePath()
	config := readConfig(configFile)

	switch key {
	case "service":
		// if empty, return ""
		return config.TerminalAIBackend
	case "model":
		return config.TerminalAIModel
	case "gemini-key":
		return config.GoogleAPIKey
	case "openai-key":
		return config.OpenAPIKey
	}

	return ""
}

func showConfig() {
	configFile := getConfigFilePath()
	config := readConfig(configFile)

	// print whole configuration with key value pairs
	fmt.Println("Sevice:", config.TerminalAIBackend)
	fmt.Println("model:", config.TerminalAIModel)
	fmt.Println("Gemini API Key:", config.GoogleAPIKey)
	fmt.Println("OpenAI API Key:", config.OpenAPIKey)
}

func setConfig(key string, value string) {
	configFile := getConfigFilePath()
	config := readConfig(configFile)

	switch key {
	case "service":
		config.TerminalAIBackend = value
	case "model":
		config.TerminalAIModel = value
	case "gemini-key":
		config.GoogleAPIKey = value
	case "openai-key":
		config.OpenAPIKey = value
	}

	saveConfig(config, configFile)

	fmt.Println("Configuration updated successfully.")

}

func resetConfig() {
	configFile := getConfigFilePath()
	config := Configuration{
		TerminalAIBackend: "",
		TerminalAIModel:   "",
		GoogleAPIKey:      "",
		OpenAPIKey:        "",
	}

	saveConfig(config, configFile)

	fmt.Println("Configuration reset successfully.")
}

func initConfig() {
	// Prompt user to set configuration

	reader := bufio.NewReader(os.Stdin)

	configFile := getConfigFilePath()
	config := readConfig(configFile)

	fmt.Println("Welcome to Terminal AI Configuration Setup")
	fmt.Println("Please set the following configuration options:")
	fmt.Println("1. Terminal AI Service (e.g., openai, gemini)")
	config.TerminalAIBackend = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())
	fmt.Println("2. Terminal AI Model (e.g., gpt-3, gpt-4-turbo, etc.)")
	config.TerminalAIModel = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())
	fmt.Println("3. Google API Key (optional)")
	config.GoogleAPIKey = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())
	fmt.Println("4. OpenAI API Key (optional)")
	config.OpenAPIKey = strings.TrimSpace(func() string { input, _ := reader.ReadString('\n'); return input }())

	saveConfig(config, configFile)

}
