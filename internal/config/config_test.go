package config_test

import (
	"os"
	"scpsave/internal/config"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCode(t *testing.T) {
	err := config.MakeSampleConfig()
	if err != nil {
		t.Fatalf("MakeSampleConfig failed: %v", err)
	}
	_ = os.Rename("./config.sample.yaml", "./config.yaml")
	err = config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	bt, err := yaml.Marshal(config.Value)
	if err != nil {
		t.Fatalf("yaml.Marshal failed: %v", err)
	}
	t.Log(string(bt))
	_ = os.Remove("./config.yaml")
}
