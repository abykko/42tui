package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"42cli/conf"
)

func DoRequest(endpoint string) (map[string]interface{}, int, error) {

	protocolPrefix, err := conf.GetString("protocol_prefix")
	if err != nil {
		return nil, 0, fmt.Errorf("getting protocol_prefix: %w", err)
	}

	addr, err := conf.GetString("server_addr")
	if err != nil {
		return nil, 0, fmt.Errorf("getting server_addr: %w", err)
	}

	port, err := conf.GetInt("server_port")
	if err != nil {
		return nil, 0, fmt.Errorf("getting server_port: %w", err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	url := fmt.Sprintf("%s%s:%d%s", protocolPrefix, addr, port, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("decoding response: %w", err)
	}

	return data, resp.StatusCode, nil
}

func DoSignedRequest(endpoint string) (map[string]interface{}, int, error) {

	protocolPrefix, err := conf.GetString("protocol_prefix")
	if err != nil {
		return nil, 0, fmt.Errorf("getting protocol_prefix: %w", err)
	}

	addr, err := conf.GetString("server_addr")
	if err != nil {
		return nil, 0, fmt.Errorf("getting server_addr: %w", err)
	}

	port, err := conf.GetInt("server_port")
	if err != nil {
		return nil, 0, fmt.Errorf("getting server_port: %w", err)
	}

	secretEnv, err := conf.GetString("secret_env_var_name")
	if err != nil {
		return nil, 0, fmt.Errorf("getting secret_env_var_name: %w", err)
	}

	secret := os.Getenv(secretEnv)
	if secret == "" {
		return nil, 0, fmt.Errorf("environment variable %q is empty", secretEnv)
	}

	// Timestamp
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Signature HMAC(timestamp)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(timestamp))
	signature := hex.EncodeToString(h.Sum(nil))

	// HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	// URL final
	url := fmt.Sprintf("%s%s:%d%s", protocolPrefix, addr, port, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Close = true

	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", signature)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	// Decode JSON
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("decoding response: %w", err)
	}

	return data, resp.StatusCode, nil
}