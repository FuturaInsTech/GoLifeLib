package initializers

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Expand all ${VAR} variable values and and set the env values
// Part of Standard usage of .env file within the application
func LoadAndExpandEnvVariables() {
	filename := ".env"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatal(".env file not found")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading .env file")
	}

	lines := strings.Split(string(data), "\n")
	envMap := make(map[string]string)

	// Step 1: parse all lines into key=value pairs
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split only at the first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		envMap[key] = val
	}

	for key, val := range envMap {
		val = expandInline(val, envMap)
		envMap[key] = val
		_ = os.Setenv(key, val)
	}
}

func expandInline(value string, envMap map[string]string) string {
	for {
		start := strings.Index(value, "${")
		if start == -1 {
			break
		}
		end := strings.Index(value[start:], "}")
		if end == -1 {
			break
		}
		end += start
		varName := value[start+2 : end]

		repl := ""
		if v, ok := envMap[varName]; ok {
			repl = v
		}
		value = value[:start] + repl + value[end+1:]
	}
	return value
}
