#!/bin/bash

set -e

echo "=== Checking and installing dependencies ==="

# Detect OS
OS_TYPE="$(uname -s)"

if [ "$OS_TYPE" = "Darwin" ]; then
    echo "macOS detected."
    
    # Check for Homebrew
    if ! command -v brew &> /dev/null; then
        echo "Homebrew is required on macOS. Please install it from https://brew.sh"
        exit 1
    fi

    # Install Podman
    if command -v podman &> /dev/null; then
        echo "Podman is already installed:" $(podman --version)
    else
        echo "Installing Podman via Homebrew..."
        brew install podman
        echo "Initializing Podman machine..."
        podman machine init || true
        podman machine start || true
    fi

    # Install Go
    if command -v go &> /dev/null; then
        echo "Go is already installed:" $(go version)
    else
        echo "Installing Go via Homebrew..."
        brew install go
    fi

elif [ "$OS_TYPE" = "Linux" ]; then
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VARIANT=${VARIANT_ID:-""}
    else
        echo "Could not detect Linux distribution."
        exit 1
    fi

    # Check Podman
    if command -v podman &> /dev/null; then
        echo "Podman is already installed:" $(podman --version)
    else
        echo "Installing Podman..."
        case $OS in
            ubuntu|debian)
                sudo apt update && sudo apt install -y podman
                ;;
            fedora)
                if [ "$VARIANT" = "silverblue" ]; then
                    echo "Podman should be preinstalled on Silverblue."
                    exit 1
                else
                    sudo dnf install -y podman
                fi
                ;;
        esac
    fi

    # Check Go
    if command -v go &> /dev/null; then
        echo "Go is already installed:" $(go version)
    else
        echo "Installing Go..."
        case $OS in
            ubuntu|debian)
                sudo apt update && sudo apt install -y golang
                ;;
            fedora)
                if [ "$VARIANT" = "silverblue" ]; then
                    echo "Fedora Silverblue detected."
                    echo "Use Toolbx to set up a dev environment:"
                    echo "  toolbox create && toolbox enter"
                    exit 0
                else
                    sudo dnf install -y golang
                fi
                ;;
            *)
                echo "Distribution '$OS' is not supported."
                exit 1
                ;;
        esac
    fi
else
    echo "Unsupported Operating System: $OS_TYPE"
    exit 1
fi

echo "All dependencies are ready!"

# --- Build and Run Steps ---

echo "=== Downloading Go modules ==="
go mod download

echo "=== Building application ==="
go build -o 42tui .

echo "=== Starting application ==="
./42tui