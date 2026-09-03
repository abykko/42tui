package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"42cli/conf"
)

func DoRequest(endpoint string) (map[string]interface{}, int, error) {
	protocolPrefix, _ := conf.GetString("protocol_prefix")
	addr, _ := conf.GetString("server_addr")
	port, _ := conf.GetInt("server_port")

	url := fmt.Sprintf("%s%s:%d%s", protocolPrefix, addr, port, endpoint)

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

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

func WaitForRequestTo(endpoint string, check func(map[string]interface{}) bool, timeout time.Duration, interval time.Duration) error {
	start := time.Now()

	for {
		resp, _, err := DoRequest(endpoint)
		if err == nil && check(resp) {
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for %s", endpoint)
		}

		time.Sleep(interval)
	}
}