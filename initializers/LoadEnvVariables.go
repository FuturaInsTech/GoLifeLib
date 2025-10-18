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

func LoadAndExpandEnvVariables() {
	localFile := ".env_local"
	mainFile := ".env"

	envMap := make(map[string]string)

	if _, err := os.Stat(localFile); os.IsNotExist(err) {
		log.Fatal(".env_local file not found")
	}
	loadEnvFile(localFile, envMap)

	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		log.Fatal(".env file not found")
	}
	loadEnvFile(mainFile, envMap)

	for key, val := range envMap {
		val = expandInline(val, envMap)
		envMap[key] = val
		_ = os.Setenv(key, val)
	}
}

func loadEnvFile(filename string, envMap map[string]string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading .env or .env_local file")
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		envMap[key] = val
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
		} else {
			repl = os.Getenv(varName) // fallback to OS environment
		}
		value = value[:start] + repl + value[end+1:]
	}
	return value
}
