package deployment

import (
	"os"
	"log"
	"os/exec"
	"crypto/rand"
	"encoding/hex"

	"42cli/conf"
)

func generateSecret() string {
	b := make([]byte, 32) // 256 bits

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func Build() error {

	secret := generateSecret()
	log.Println("[build] generated secret")

	secretEnv, err := conf.GetString("secret_env_var_name")
	if err != nil {
		log.Println("[build] error getting secret env name:", err)
		return err
	}

	if err := os.Setenv(secretEnv, secret); err != nil {
		log.Println("[build] error setting SECRET env var:", err)
		return err
	}
	log.Println("[build] set env var %s\n", secretEnv)

	servDir, err := conf.GetString("server_dir")
	if err != nil {
		log.Println("[build] error getting server_dir:", err)
		return err
	}

	imageName, err := conf.GetString("podman_image_name")
	if err != nil {
		log.Println("[build] error getting podman_image_name:", err)
		return err
	}

	log.Println("[build] building podman image %s in %s\n", imageName, servDir)

	StopByImage()

	cmd := exec.Command("podman", "build", "-t", imageName, ".")
	cmd.Dir = servDir

	if err := cmd.Run(); err != nil {
		log.Println("[build] error building podman image:", err)
		return err
	}

	log.Println("[build] podman build completed successfully")

	return nil
}
