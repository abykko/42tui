package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func DoSignedRequest(url string, secret string) (map[string]interface{}, int, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(timestamp))
	signature := hex.EncodeToString(h.Sum(nil))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", signature)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, resp.StatusCode, nil
	}

	return data, resp.StatusCode, nil
}