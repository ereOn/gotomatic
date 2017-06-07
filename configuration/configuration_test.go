package configuration

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestLoadFailReader(t *testing.T) {
	f, _ := os.Open("fixture/configuration.yaml")
	f.Close()
	_, err := Load(f)

	if err == nil {
		t.Error("expected an error")
	}
}

func TestLoadFailYAML(t *testing.T) {
	f, _ := os.Open("fixture/invalid")
	defer f.Close()
	_, err := Load(f)

	if err == nil {
		t.Error("expected an error")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	f, _ := os.Open("fixture/invalid.yaml")
	defer f.Close()
	_, err := Load(f)

	if err == nil {
		t.Error("expected an error")
	}
}

func TestLoad(t *testing.T) {
	f, _ := os.Open("fixture/configuration.yaml")
	defer f.Close()
	_, err := Load(f)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}
