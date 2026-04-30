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
		return nil, fmt.Errorf("open config file %s: %w", path, err)
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
			return nil, fmt.Errorf("invalid config line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty config key in line: %s", line)
		}

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file %s: %w", path, err)
	}

	return config, nil
}

// LoadConfig loads configuration from the given path.
// If debug is true, it prints the loaded configuration.
func LoadConfig(path string, debug bool) (map[string]string, error) {
	cfg, err := readConfig(path)
	if err != nil {
		return nil, fmt.Errorf("error reading settings file: %w", err)
	}

	if debug {
		fmt.Println("== LOADED SETTINGS ==")
		for k, v := range cfg {
			fmt.Printf("%s = %s\n", k, v)
		}
		fmt.Println("== END SETTINGS ==")
	}

	return cfg, nil
}