package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	LoadConfig()
	for _, model := range models {
		t.Logf("model: %s", model.Model)
	}
}
