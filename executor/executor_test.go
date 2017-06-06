package executor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFalseExecutor(t *testing.T) {
	value := FalseExecutor(context.Background())

	if value {
		t.Error("expected false")
	}
}

func TestTrueExecutor(t *testing.T) {
	value := TrueExecutor(context.Background())

	if !value {
		t.Error("expected true")
	}
}

func TestCommandExecutor(t *testing.T) {
	value := CommandExecutor("echo")(context.Background())

	if !value {
		t.Error("expected true")
	}

	value = CommandExecutor("unknown-command")(context.Background())

	if value {
		t.Error("expected false")
	}
}

func TestHTTPExecutor(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	}))

	value := HTTPExecutor("GET", ts.URL, []int{200}, time.Second)(context.Background())

	if !value {
		t.Error("expected true")
	}

	value = HTTPExecutor("GET", ts.URL, []int{201}, time.Second)(context.Background())

	if value {
		t.Error("expected false")
	}

	value = HTTPExecutor("GET", "http://localhost:0", []int{200}, time.Second)(context.Background())

	if value {
		t.Error("expected false")
	}

	value = HTTPExecutor("🖕", ts.URL, []int{200}, time.Second)(context.Background())

	if value {
		t.Error("expected false")
	}
}
