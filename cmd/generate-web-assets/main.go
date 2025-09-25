package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get the project root directory
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory:", err)
	}

	// Build the web app
	fmt.Println("Building web application...")
	webDir := filepath.Join(projectRoot, "apps", "web")
	
	// Install dependencies
	cmd := exec.Command("pnpm", "install")
	cmd.Dir = webDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to install web dependencies:", err)
	}

	// Build the web app
	cmd = exec.Command("pnpm", "build")
	cmd.Dir = webDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to build web app:", err)
	}

	// Copy dist to internal/web
	fmt.Println("Copying web assets to internal/web...")
	srcDir := filepath.Join(webDir, "dist")
	dstDir := filepath.Join(projectRoot, "internal", "web", "dist")

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dstDir), 0755); err != nil {
		log.Fatal("Failed to create destination directory:", err)
	}

	// Copy the dist folder
	cmd = exec.Command("cp", "-r", srcDir, dstDir)
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to copy web assets:", err)
	}

	fmt.Println("Web assets generated successfully!")
}