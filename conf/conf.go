package conf

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	// Configuration file path
	configPath   = "conf/.conf"

	// This variable stores de settings
	cachedConfig map[string]string
	configHash   [32]byte

	mu sync.RWMutex
)

func calculateHash(path string) ([32]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return [32]byte{}, fmt.Errorf("read config file: %w", err)
	}

	return sha256.Sum256(data), nil
}

func ReadConfig(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file %s: %w", path, err)
	}
	defer file.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

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
		return nil, fmt.Errorf("scan config: %w", err)
	}

	return config, nil
}

func loadConfig() (map[string]string, error) {
	mu.Lock()
	defer mu.Unlock()

	newHash, err := calculateHash(configPath)
	if err != nil {
		return nil, err
	}

	// We compare the cached settings hash with the new hash
	// if they're different we read and store the settings
	// file again.
	if cachedConfig == nil || newHash != configHash {
		cfg, err := ReadConfig(configPath)
		if err != nil {
			return nil, err
		}

		cachedConfig = cfg
		configHash = newHash
	}

	return cachedConfig, nil
}

// Get setting as a string. Get() is an alias method refering to GetString()
func GetString(key string) (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", fmt.Errorf("config error: %w", err)
	}

	value, ok := cfg[key]
	if !ok {
		return "", fmt.Errorf("config key not found: %s", key)
	}

	return value, nil
}
func Get(key string) (string, error) {
	return GetString(key)
}

func GetInt(key string) (int, error) {
	value, err := GetString(key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid int for %s: %w", key, err)
	}

	return i, nil
}

func GetBool(key string) (bool, error) {
	value, err := GetString(key)
	if err != nil {
		return false, err
	}

	// Support different binary selection indicator
	switch strings.ToLower(value) {
	case "true", "1", "yes", "y", "on":
		return true, nil
	case "false", "0", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool for %s: %s", key, value)
	}
}
