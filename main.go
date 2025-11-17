package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var dbPath string
	var port string
	flag.StringVar(&dbPath, "db", "todos.db", "sqlite db file")
	flag.StringVar(&port, "port", "8080", "server port")
	flag.Parse()

	// store
	store, err := NewStore(dbPath)
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}

	// hub (SSE)
	hub := NewHub()
	go hub.Run()

	// static + sse
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/sse", sseHandler(hub))

	// api routes
	http.HandleFunc("/api/todos", todosHandler(store, hub))
	http.HandleFunc("/api/todos/", todoDetailHandler(store, hub)) // handles /api/todos/{id} and /api/todos/{id}/complete

	// keep simple notify for quick tests (optional)
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		// legacy quick test API
		if r.Method == http.MethodPost {
			// already handled in code elsewhere
			w.WriteHeader(http.StatusNoContent)
			return
		}
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			msg = "hello"
		}
		hub.Broadcast(msg)
		w.WriteHeader(http.StatusNoContent)
	})

	srv := &http.Server{
		Addr: ":" + port,
		// ReadTimeout left small; WriteTimeout=0 allowed for SSE in handlers where we'll flush periodically
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	// graceful shutdown
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("server stopped")
}
