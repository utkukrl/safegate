package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/utkukrl/safegate/internal/core"
)

func main() {
	configPath := "../../configs/config.yaml"
	cfg, err := core.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	parsedURL, err := url.Parse(cfg.Proxy.Target + ":" + cfg.Proxy.Port)
	if err != nil {
		log.Fatalf("failed to parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	ruleSet := core.NewRuleSet()

	if cfg.REPL.Enabled {
		go func() {
			core.StartREPL(ruleSet)
		}()
	}

	certFile := "../../scripts/certs/server.crt"
	keyFile := "../../scripts/certs/server.key"

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Printf("Certificate file not found: %s", certFile)
		log.Printf("Please run scripts/generate_certs.sh to generate certificates")
		os.Exit(1)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Printf("Key file not found: %s", keyFile)
		log.Printf("Please run scripts/generate_certs.sh to generate certificates")
		os.Exit(1)
	}

	startServer(":"+cfg.Dashboard.Port, handler, certFile, keyFile)
}

func startServer(addr string, handler http.Handler, certFile, keyFile string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting HTTPS server at %s", addr)
		serverErrors <- srv.ListenAndServeTLS(certFile, keyFile)
	}()

	select {
	case err := <-serverErrors:
		log.Printf("Server error: %v", err)
	case <-stop:
		log.Printf("Shutting down server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
		os.Exit(1)
	}

	log.Printf("Server gracefully stopped")
}
