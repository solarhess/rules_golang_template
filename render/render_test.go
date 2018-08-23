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
func assertEquals(actual string, expected string) {
	if actual != expected {
		panic(fmt.Sprintf("Expected \"%s\", was \"%s\"", expected, actual))
	}
}

func TestItWorks(t *testing.T) {
	result, data, err := Render("test.tpl", "test.json")
	if err != nil {
		panic(fmt.Sprintf("Error %s", err))
	}
	assertEquals(data["Type"].(string), "Sweaters")
	assertEquals(data["File"].(map[string]interface{})["Message"].(string), "A string from inner json")
	assertEquals(data["Text"].(string), "message in text file")
	assertContainsString(result, "25")
	assertContainsString(result, "message in text file")
	assertContainsString(result, "A string from inner json")
}

func TestEnv(t *testing.T) {
	expectedUserValue := os.Getenv("USER")
	result, data, err := Render("env-test.tpl", "")
	if err != nil {
		panic(fmt.Sprintf("Error %s", err))
	}
	assertContainsString(result, fmt.Sprintf("User: %s", expectedUserValue))
	assertEquals(data["USER"].(string), expectedUserValue)
}
