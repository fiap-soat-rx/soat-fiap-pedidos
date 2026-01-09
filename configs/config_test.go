package configs

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	os.Unsetenv("SERVER_PORT")
	c := LoadConfig()
	if c.ServerPort == "" {
		t.Fatalf("expected default server port set")
	}
}

func TestLoadConfig_WithEnvOverrides(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	c := LoadConfig()
	if c.ServerPort != "9090" {
		t.Fatalf("expected server port 9090, got %s", c.ServerPort)
	}
	if c.LogLevel != "debug" {
		t.Fatalf("expected log level debug, got %s", c.LogLevel)
	}
}

func TestGetEnvOrPanic_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when default empty and env not set")
		}
	}()
	os.Unsetenv("SOME_ENV_VAR")
	_ = getEnvOrPanic("SOME_ENV_VAR", "")
}

func TestGetEnvAsBool_Behavior(t *testing.T) {
	os.Setenv("BOOL_TRUE", "true")
	if !getEnvAsBool("BOOL_TRUE", false) {
		t.Fatalf("expected true from BOOL_TRUE")
	}
	os.Setenv("BOOL_INVALID", "notbool")
	if !getEnvAsBool("BOOL_INVALID", true) {
		t.Fatalf("expected default true on invalid parse")
	}
}
