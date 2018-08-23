package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func assertContainsString(value string, substr string) {
	if !strings.Contains(value, substr) {
		panic(fmt.Sprintf("Expected \"%s\", in \"%s\"", substr, value))
	}
}

func TestItWorks(t *testing.T) {
	result, err := Render("test.tpl", "test.json")
	if err != nil {
		panic(fmt.Sprintf("Error %s", err))
	}
	assertContainsString(result, "25")
	assertContainsString(result, "message in text file")
	assertContainsString(result, "A string from inner json")
}

func TestEnv(t *testing.T) {
	expectedUserValue := os.Getenv("USER")
	result, err := Render("env-test.tpl", "")
	if err != nil {
		panic(fmt.Sprintf("Error %s", err))
	}
	assertContainsString(result, fmt.Sprintf("User: %s", expectedUserValue))
}
