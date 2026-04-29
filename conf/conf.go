// conf.go

package conf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readConfig(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		config[key] = value
	}

	return config, scanner.Err()
}

func LoadConfig() (map[string]string){
	cfg, err := readConfig("conf/.conf")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return nil
	}
	return cfg
}