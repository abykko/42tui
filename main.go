package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"42cli/tui"
)

// checkCommand verifica si una herramienta existe en el PATH del sistema
func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkFile verifica si un archivo existe
func checkFile(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func checkDependencies() error {
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

	return nil
}

func main() {
	// Configuración de logs en archivo
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// Comprobar dependencias antes de arrancar
	if err := checkDependencies(); err != nil {
		log.Printf("Error de dependencias: %v\n", err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Si todo está correcto, inicia la TUI
	tui.Tui()
}