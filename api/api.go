package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"42cli/conf"
)

const logFileName = "api_logs.json"

var logger *slog.Logger

func init() {
	f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
		logger.Error("could not open log file, using stderr", "error", err)
	} else {
		logger = slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
}

func DoRequest(endpoint string) (map[string]interface{}, int, error) {
	protocolPrefix, _ := conf.GetString("protocol_prefix")
	addr, _ := conf.GetString("server_addr")
	port, _ := conf.GetInt("server_port")

	url := fmt.Sprintf("%s%s:%d%s", protocolPrefix, addr, port, endpoint)
	start := time.Now()

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("error creating the request",
			"url", url,
			"error", err)
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("request failed (network/timeout)",
			"url", url,
			"duration_ms", time.Since(start).Milliseconds(),
			"error", err)
		return nil, 0, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Error("error decoding JSON response",
			"url", url,
			"status", resp.StatusCode,
			"error", err)
		return nil, resp.StatusCode, fmt.Errorf("decoding response: %w", err)
	}

	duration := time.Since(start).Milliseconds()

	if resp.StatusCode >= 400 {
		logger.Warn("API returned an HTTP error",
			"url", url,
			"status", resp.StatusCode,
			"duration_ms", duration,
			"response_body", data)
	} else {
		logger.Info("request completed",
			"url", url,
			"status", resp.StatusCode,
			"duration_ms", duration)
	}

	return data, resp.StatusCode, nil
}

func WaitForRequestTo(
	endpoint string,
	check func(map[string]interface{}) bool,
	timeout time.Duration,
	interval time.Duration,
) error {
	start := time.Now()
	attempt := 0

	for {
		attempt++
		resp, _, err := DoRequest(endpoint)

		if err == nil && check(resp) {
			return nil
		}

		if time.Since(start) > timeout {
			logger.Error("timeout waiting for condition",
				"endpoint", endpoint,
				"attempts", attempt,
				"total_time_sec", time.Since(start).Seconds())
			return fmt.Errorf("timeout waiting for %s", endpoint)
		}

		time.Sleep(interval)
	}
}
