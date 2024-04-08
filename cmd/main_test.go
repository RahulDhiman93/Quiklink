package main

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	err := os.Chdir("..")
	if err != nil {
		t.Fatal("Failed to change directory:", err)
	}

	_, err = run()
	if err != nil {
		t.Error("Failed Run()")
	}
}
