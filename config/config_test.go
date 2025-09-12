package config

import (
	"os"
	"sparrow-cli/env"
	"sparrow-cli/global"
	"testing"
)

func init() {
	homePath := os.Getenv("SparrowCliHome")
	if homePath == "" {
		homePath = os.Getenv("HOME") + "/.sparrow-cli"
	}
	env.SparrowCliHome = homePath
}

func TestLoadConfig(t *testing.T) {
	LoadConfig()
	for _, model := range Models {
		t.Logf("model: %s", model.Model)
	}

	t.Logf("logger: %+v", Logger)

	t.Logf("current model: %+v", global.CurrentModel)
}
