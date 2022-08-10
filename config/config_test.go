package config

import (
	"testing"
)

func TestConfigServerHostPort(t *testing.T) {
	server := &ConfigServer{
		Host: "hostname",
		Port: "0000",
	}

	result := server.HostPort()
	if result != "hostname:0000" {
		t.Logf("Expected hostname:0000, got %s", result)
		t.Fail()
	}
}
