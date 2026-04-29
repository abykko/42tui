package main

import (
	"fmt"
	"strconv"

	"42cli/conf"
	"42cli/server"
	"42cli/deployment"
)

func main() {

	// Generate secret key
	secret := server.GenerateSecret()
	fmt.Println(secret)

	// Load .conf file settings
	cfg := conf.LoadConfig(true)

	port, err := strconv.Atoi(cfg["server_port"])
	if err != nil {
		fmt.Println("Conversion error:", err)
		return
	}

	imageName := cfg["podman_image_name"]
	serverDir := cfg["server_dir"]

	// Build podman server image
	deployment.BuildServerPodman(imageName, serverDir)

	// Run podman server
	podmanId, err := deployment.RunServerPodman(imageName, secret, port)
	if err != nil {
		fmt.Println("Error running podman server.")
		return
	}

	fmt.Println("Container ID:", podmanId)

	fmt.Println("Press ENTER to stop server...")

	// Espera Enter
	fmt.Scanln()

	// STOP container
	err = deployment.StopServerPodman(podmanId)
	if err != nil {
		fmt.Println("Error stopping container:", err)
		return
	}

	fmt.Println("Server stopped cleanly")
}