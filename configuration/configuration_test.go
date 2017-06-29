package configuration

import (
	"context"
	"os"
	"testing"

	"github.com/intelux/gotomatic/conditional"
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

	conf, err := Load(f)
	defer conf.Close()

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}

func TestWatch(t *testing.T) {
	f, _ := os.Open("fixture/configuration.yaml")
	defer f.Close()

	conf, _ := Load(f)
	defer conf.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	err := conf.Watch(ctx)

	if err != nil {
		t.Errorf("expected no error but got: %s", err)
	}
}

func TestWatchTriggerFailure(t *testing.T) {
	f, _ := os.Open("fixture/configuration.yaml")
	defer f.Close()

	conf, _ := Load(f)
	defer conf.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	condition := conf.GetCondition("a")
	go condition.(conditional.Settable).Set(true)

	err := conf.Watch(ctx)

	if err == nil {
		t.Error("expected an error")
	}
}
