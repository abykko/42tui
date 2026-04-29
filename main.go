// main.go

package main

import (
	"fmt"
	"strconv"
	"42cli/conf"
	"42cli/server"
	"42cli/deployment"
)

func main() {

	// Load .conf file settings
	cfg := conf.LoadConfig()

	port, err := strconv.Atoi(cfg["server_port"])
	if err != nil {
		fmt.Println("Conversion error:", err)
		return
	}

	// Generate secret key
	secret := server.GenerateSecret()
	// secret = "foo"

	fmt.Println(port)
	fmt.Println(secret)

	deployment.BuildServerDocker(port)

	// Simple request to server
	resp, status, err := server.DoSignedRequest(
		"http://127.0.0.1:6742/status",
		secret,
	)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	fmt.Println("Status:", status)
	fmt.Println("Response:", resp)

	fmt.Println("Press enter to finish...")
	fmt.Scanln()
}