package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	
	"42tui/tui"
	"42tui/conf"
)

/*
	Verifica la existencia de una utilidad dada
	por su nombre.
*/
func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

/*
	Verifica la integridad de un archivo dado
	por su nombre.
*/
func checkFile(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func checkSetting(key string, value string) bool {
		result, err := conf.GetString(key)
	if err != nil {
		return false
	}

	cleanResult := strings.TrimSpace(result)

	if cleanResult == "" {
		return false
	}

	if cleanResult != strings.TrimSpace(value) {
		return false
	}

	return true
}

func checkIsEmpty(key string) bool {
	result, err := conf.GetString(key)
	if err != nil {
		return false
	}

	cleanResult := strings.TrimSpace(result)

	if cleanResult == "" {
		return true
	}

	return false
}

func checker() error {
	// 1. Verificar si Podman está instalado
	if !checkCommand("podman") {
		return fmt.Errorf("podman no está instalado o no se encuentra en el PATH")
	}

	// 2. Verificar si Go está instalado en el sistema
	if !checkCommand("go") {
		return fmt.Errorf("golang (go) no está instalado o no se encuentra en el PATH")
	}

	// 3. Verificar si existe el archivo de dependencias del proyecto
	if !checkFile("go.mod") {
		return fmt.Errorf("no se encontró el archivo go.mod en el directorio actual")
	}

	// 4. Verificar si autologin está activado y las credenciales están a disposición
	if checkSetting("autologin", "yes") {
		if (checkIsEmpty("user_login") || checkIsEmpty("password_login")) {
			return fmt.Errorf("autologin está activado y las credenciales están vacias: /conf/.conf")
		}
	}
	
	// Queda pendiente añadir mas comprobaciones si fuera necesario.

	return nil
}

func main() {
	/*
		Preparamos el archivo de logs.
	*/
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	/*
		Realizamos comprobaciones previas al arranque.
	*/
	if err := checker(); err != nil {
		log.Printf("Error de dependencias: %v\n", err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	/*
		Arranque de la interfaz.
	*/
	tui.Tui()
}