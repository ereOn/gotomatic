// Package executor implements executors for cut-off conditions.
package executor

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

// An Executor is a callback that returns a boolean status.
type Executor func(ctx context.Context) bool

// FalseExecutor returns always false.
func FalseExecutor(ctx context.Context) bool { return false }

// TrueExecutor returns always true.
func TrueExecutor(ctx context.Context) bool { return true }

// CommandExecutor returns an Executor that runs an external command.
func CommandExecutor(command string, args ...string) Executor {
	return func(ctx context.Context) bool {
		return exec.CommandContext(ctx, command, args...).Run() == nil
	}
}

// HTTPExecutor returns an Executor that runs a HTTP request.
func HTTPExecutor(method string, url string, statusCodes []int, timeout time.Duration) Executor {
	return func(ctx context.Context) bool {
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			return false
		}

		req = req.WithContext(ctx)
		client := &http.Client{Timeout: timeout}

		resp, err := client.Do(req)

		if err != nil {
			return false
		}

		io.Copy(ioutil.Discard, resp.Body)

		for _, statusCode := range statusCodes {
			if statusCode == resp.StatusCode {
				return true
			}
		}

		return false
	}
}
