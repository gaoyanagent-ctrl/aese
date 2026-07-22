package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/httpapi"
)

const usage = `Usage:
  aese-server [flags]

Options:
  --listen <addr>         HTTP listen address (default :8090)
  --pack-dir <path>       scenario pack directory (default scenario-packs/hctm)
  --request-timeout <dur>  request timeout, e.g. 30s (default 30s)
  --body-limit <bytes>    max request body bytes (default 1048576)
`

func main() {
	if code := run(os.Args[1:]); code != 0 {
		os.Exit(code)
	}
}

func run(args []string) int {
	fs := flag.NewFlagSet("aese-server", flag.ContinueOnError)
	listen := fs.String("listen", ":8090", "http listen address")
	packDir := fs.String("pack-dir", "scenario-packs/hctm", "scenario pack directory")
	timeout := fs.Duration("request-timeout", 30*time.Second, "request timeout")
	bodyLimit := fs.Int64("body-limit", 1<<20, "request body byte limit")
	showHelp := fs.Bool("help", false, "show usage")

	fs.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *showHelp {
		fmt.Fprint(os.Stdout, usage)
		return 0
	}

	if *listen == "" {
		fmt.Fprintln(os.Stderr, "--listen is required")
		return 2
	}
	if !strings.Contains(*listen, ":") {
		fmt.Fprintln(os.Stderr, "--listen must include port, e.g. :8090")
		return 2
	}
	if *bodyLimit <= 0 {
		fmt.Fprintln(os.Stderr, "--body-limit must be greater than 0")
		return 2
	}
	if *timeout <= 0 {
		fmt.Fprintln(os.Stderr, "--request-timeout must be greater than 0")
		return 2
	}

	logger := log.New(os.Stdout, "[aese-server] ", log.LstdFlags|log.Lshortfile)
	server := httpapi.New(httpapi.Config{
		PackDir:        *packDir,
		RequestTimeout: *timeout,
		BodyLimit:      *bodyLimit,
		Logger:         logger,
	})

	httpServer := &http.Server{
		Addr:              *listen,
		Handler:           server,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      *timeout + 10*time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Printf("starting aese-server on %s (pack-dir=%s)", *listen, *packDir)
		errCh <- httpServer.ListenAndServe()
	}()

	<-shutdownCtx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Printf("shutdown failed: %v", err)
		return 1
	}
	if err := <-errCh; err != nil && err != http.ErrServerClosed {
		logger.Printf("server exit: %v", err)
		return 1
	}
	return 0
}
