package handler

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type CgiBinHandlerOptions struct {
	Root    *os.Root
	URI     string
	Timeout time.Duration
}

func CgiBinHandler(opts *CgiBinHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()
		logger = logger.With("path", req.URL.Path)

		switch req.Method {
		case "GET", "POST":
			// pass
		default:
			http.Error(rsp, "Invalid request", http.StatusBadRequest)
			return
		}

		script_path := filepath.Base(req.URL.Path)

		script_info, err := opts.Root.Stat(script_path)

		if err != nil {
			http.Error(rsp, "Script not found or not executable", http.StatusNotFound)
			return
		}

		if !script_info.Mode().IsRegular() {
			http.Error(rsp, "Script not found or not executable", http.StatusNotFound)
			return
		}

		ctx, cancel := context.WithTimeout(req.Context(), opts.Timeout)
		defer cancel()

		env := os.Environ()
		env = append(env,
			"REQUEST_METHOD="+req.Method,
			"SCRIPT_FILENAME="+script_path,
			"QUERY_STRING="+req.URL.RawQuery,
		)

		if contentLength := req.ContentLength; contentLength >= 0 {
			env = append(env, "CONTENT_LENGTH="+fmt.Sprintf("%d", contentLength))
		}

		cmd := exec.CommandContext(ctx, script_path)
		cmd.Env = env

		var stdin io.WriteCloser

		if req.Method == "POST" {

			stdin, err = cmd.StdinPipe()

			if err != nil {
				http.Error(rsp, "Failed to create stdin pipe", http.StatusInternalServerError)
				return
			}
			go func() {
				defer stdin.Close()
				_, _ = io.Copy(stdin, req.Body)
			}()
		}

		stdout, err := cmd.StdoutPipe()

		if err != nil {
			http.Error(rsp, "Failed to create stdout pipe", http.StatusInternalServerError)
			return
		}
		stderr, err := cmd.StderrPipe()

		if err != nil {
			http.Error(rsp, "Failed to create stderr pipe", http.StatusInternalServerError)
			return
		}

		err = cmd.Start()

		if err != nil {
			http.Error(rsp, "Failed to start CGI script: "+err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {

			defer rsp.(http.Flusher).Flush()
			rsp.WriteHeader(http.StatusOK)

			io.Copy(rsp, stdout)

			if err := cmd.Wait(); err != nil && ctx.Err() == nil { // Check if error is not context cancellation
				logger.Error("Script failed", "error", err)
			}
		}()

		b, _ := io.ReadAll(stderr) // Read stderr to capture any errors from the script

		if ctx.Err() != nil {
			logger.Error("Script errors", "error", string(b))
			http.Error(rsp, "Script execution timed out", http.StatusGatewayTimeout)
			return
		}
	}

	return http.HandlerFunc(fn), nil
}
