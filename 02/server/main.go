package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	const root = "/srv/site"
	files := http.FileServer(http.Dir(root))
	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				w.Header().Set("Allow", "GET, HEAD")
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			clean := path.Clean("/" + r.URL.Path)
			info, err := os.Stat(filepath.Join(root, clean))
			if err == nil && !info.IsDir() {
				files.ServeHTTP(w, r)
				return
			}
			if clean == "/" || (path.Ext(clean) == "" && strings.Contains(r.Header.Get("Accept"), "text/html")) {
				http.ServeFile(w, r, filepath.Join(root, "index.html"))
				return
			}
			http.NotFound(w, r)
		}),
	}
	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("static server: %v", err)
			os.Exit(1)
		}
	}()

	traefik := exec.Command("/usr/local/bin/traefik", "--configFile=/etc/traefik/traefik.yml")
	traefik.Stdout, traefik.Stderr = os.Stdout, os.Stderr
	if err := traefik.Start(); err != nil {
		log.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- traefik.Wait() }()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-done:
		log.Fatalf("Traefik stopped unexpectedly: %v", err)
	case sig := <-signals:
		_ = traefik.Process.Signal(sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
		select {
		case <-done:
		case <-ctx.Done():
			_ = traefik.Process.Kill()
			<-done
		}
	}
}
