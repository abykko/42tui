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

	"42tui/conf"
)

// Genera la firma HMAC-SHA256 igual que la función validate_sign del servidor
func generateSignature(secret string, timestamp string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

func DoRequest(endpoint string) (map[string]interface{}, int, error) {
	protocolPrefix, _ := conf.GetString("protocol_prefix")
	addr, _ := conf.GetString("server_addr")
	port, _ := conf.GetInt("server_port")
	secretEnv, _ := conf.GetString("secret_env_var_name")

	// Obtenemos el secreto generado dinámicamente en el build/start
	secret := os.Getenv(secretEnv)

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

	// Inyectamos las cabeceras de seguridad requeridas por el servidor
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := generateSignature(secret, timestamp)

	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", signature)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("decoding response (status %d): %w", resp.StatusCode, err)
	}

	return data, resp.StatusCode, nil
}

func WaitForRequestTo(endpoint string, check func(map[string]interface{}) bool, timeout time.Duration, interval time.Duration) error {
	start := time.Now()

	for {
		resp, status, err := DoRequest(endpoint)
		// Verificamos que no haya error de red, que la respuesta sea HTTP 200 y que pase el check
		if err == nil && status == http.StatusOK && check(resp) {
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for %s", endpoint)
		}

		time.Sleep(interval)
	}
}