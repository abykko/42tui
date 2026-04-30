package main

import (
	"testing"
	"42cli/deployment"
)

func TestDeployAndStop(t *testing.T) {
	id, err := deployment.DeployServer()
	if err != nil {
		t.Fatal(err)
	}

	err = deployment.StopServerPodman(id)
}