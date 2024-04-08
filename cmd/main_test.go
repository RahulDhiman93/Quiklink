package main

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	err := os.Chdir("..")
	if err != nil {
		t.Fatal("Failed to change directory:", err)
	}

	err = godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file for testing:", err)
	}
	_, err = run()
	if err != nil {
		t.Error("Failed Run()")
	}
}
