package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/creative"
	"github.com/industrial-ai/iaos-aese/internal/genesisworkspace"
	"github.com/industrial-ai/iaos-aese/internal/httpapi"
)

const usage = `Usage:
  aese-server [flags]

Options:
  --listen <addr>         HTTP listen address (default :8090)
  --pack-dir <path>       scenario pack directory (default scenario-packs/hctm)
  --iaos-base-url <url>   live IAOS API; absent uses deterministic offline projection
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
	iaosBaseURL := fs.String("iaos-base-url", "", "live IAOS API base URL")
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
	authMode := strings.ToLower(strings.TrimSpace(envOrDefault("AESE_AUTH_MODE", "iaos")))
	if authMode != "iaos" && authMode != "local_dev" {
		fmt.Fprintln(os.Stderr, "AESE_AUTH_MODE must be iaos or local_dev")
		return 2
	}
	if authMode == "local_dev" && !isLoopbackListen(*listen) {
		fmt.Fprintln(os.Stderr, "AESE_AUTH_MODE=local_dev may only listen on 127.0.0.1 or localhost")
		return 2
	}

	logger := log.New(os.Stdout, "[aese-server] ", log.LstdFlags|log.Lshortfile)
	var creativeProvider creative.Provider = creative.DeterministicProvider{}
	if key := strings.TrimSpace(os.Getenv("MINMAX_API_KEY")); key != "" {
		provider, err := creative.NewMiniMaxProvider(creative.MiniMaxConfig{
			BaseURL: os.Getenv("MINMAX_API_BASE"),
			APIKey:  key,
			Model:   os.Getenv("MINMAX_MODEL"),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid MiniMax configuration: %v\n", err)
			return 2
		}
		creativeProvider = provider
		logger.Printf("creative provider enabled (provider=MiniMax model=%s)", os.Getenv("MINMAX_MODEL"))
	} else {
		logger.Printf("creative provider fallback enabled (provider=deterministic)")
	}
	server := httpapi.New(httpapi.Config{
		PackDir:               *packDir,
		IAOSBaseURL:           *iaosBaseURL,
		RequestTimeout:        *timeout,
		BodyLimit:             *bodyLimit,
		Logger:                logger,
		CreativeProvider:      creativeProvider,
		CreativeJobStore:      creative.NewJobStore(envOrDefault("GENESIS_CREATIVE_JOB_STORE", ".aese-data/genesis-creative-jobs.json")),
		GenesisPlayerAuth:     &genesisworkspace.PlayerAuthClient{BaseURL: *iaosBaseURL},
		AllowLocalGenesisAuth: authMode == "local_dev",
		GenesisWorkspaceService: &genesisworkspace.Service{
			Store:        genesisworkspace.NewStore(envOrDefault("GENESIS_WORKSPACE_STORE", ".aese-data/genesis-workspaces.json")),
			ControlPlane: &genesisworkspace.ControlPlaneClient{BaseURL: *iaosBaseURL},
			Provisioner: genesisworkspace.IAOSClient{
				BaseURL:       *iaosBaseURL,
				PlatformToken: os.Getenv("GENESIS_PLATFORM_TOKEN"),
			},
		},
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

func isLoopbackListen(address string) bool {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return false
	}
	host = strings.Trim(host, "[]")
	return strings.EqualFold(host, "localhost") || net.ParseIP(host).IsLoopback()
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
