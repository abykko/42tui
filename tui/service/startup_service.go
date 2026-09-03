package service

import (
	"log"
	"time"

	"42cli/api"
	deployment "42cli/server-deployment"
)

func Startup() (bool, error) {

	// Build container with the server
	log.Println("[Startup] building container...")
	if err := deployment.Build(); err != nil {
		log.Println("[Startup] build error:", err)
		return false, err
	}

	// Run the server container once is ready
	log.Println("[Startup] running container...")
	if _, err := deployment.Run(); err != nil {
		log.Println("[Startup] run error:", err)
		return false, err
	}

	// Wait until the API is ready (has a timeout)
	log.Println("[Startup] waiting for API endpoint: /status")
	err := api.WaitForRequestTo(
		"/status",
		func(resp map[string]interface{}) bool {
			log.Println("[Startup] /status response:", resp)
			return resp["status"] == "ok"
		},
		10*time.Second,
		30*time.Millisecond,
	)
	if err != nil {
		log.Println("[Startup] API error:", err)
		return false, err
	}

	log.Println("[Startup] READY")
	return true, nil
}