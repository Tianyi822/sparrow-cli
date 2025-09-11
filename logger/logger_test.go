package logger

import (
	"context"
	"os"
	"sparrow-cli/config"
	"sparrow-cli/env"
	"testing"
)

func init() {
	homePath := os.Getenv("SparrowCliHome")
	if homePath == "" {
		homePath = os.Getenv("HOME") + "/.sparrow-cli"
	}
	env.SparrowCliHome = homePath

	config.LoadConfig()
}

func TestInitLogger(t *testing.T) {
	ctx := context.Background()
	err := InitLogger(ctx)
	if err != nil {
		t.Errorf("InitLogger() error = %v", err)
	}
	Info("hello world")
}
