package fibo

import (
	"log"
	"os"
	"testing"
)

func TestTableGetEnv(t *testing.T) {
	// prepare env
	var envs = []struct {
		name  string
		value string
	}{
		{"TEST-VALUE", "value"},
		{"TEST-EMPTY", ""},
	}

	for _, env := range envs {
		err := os.Setenv(env.name, env.value)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := os.Unsetenv("TEST-EXIST")
	if err != nil {
		log.Fatal(err)
	}

	// run tests
	var tests = []struct {
		name     string
		fallback string
		expected string
	}{
		{"TEST-VALUE", "fallback", "value"},
		{"TEST-EMPTY", "fallback", ""},
		{"TEST-EXIST", "fallback", "fallback"},
	}

	for _, test := range tests {
		output := getEnv(test.name, test.fallback)
		if output != test.expected {
			t.Error("Test Failed: {} name, {} fallback, {} expected, recieved: {}",
				test.name, test.fallback, test.expected, output)
		}
	}

}
