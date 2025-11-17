package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

// Client represents a single SSE connection.
type Client struct {
    id   string
    ch   chan string
    done chan struct{}
}

// Hub manages clients and broadcasts messages.
type Hub struct {
    mu         sync.RWMutex
    clients    map[string]*Client
    register   chan *Client
    unregister chan *Client
    broadcastC chan string
}

// NewHub creates a Hub.
func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]*Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcastC: make(chan string, 256),
    }
}

// Run loop handles register/unregister/broadcast.
func (h *Hub) Run() {
    for {
        select {
        case c := <-h.register:
            h.mu.Lock()
            h.clients[c.id] = c
            h.mu.Unlock()
        case c := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[c.id]; ok {
                close(c.ch)
                delete(h.clients, c.id)
            }
            h.mu.Unlock()
        case msg := <-h.broadcastC:
            h.mu.RLock()
            for _, c := range h.clients {
                select {
                case c.ch <- msg:
                default:
                    // client channel full, drop (避免阻塞)
                }
            }
            h.mu.RUnlock()
        }
    }
}

// Broadcast sends a message to all clients.
func (h *Hub) Broadcast(msg string) {
    select {
    case h.broadcastC <- msg:
    default:
        // drop if hub overloaded
    }
}

// sseHandler returns an SSE HTTP handler.
func sseHandler(h *Hub) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 认证：从 query token 获取（示例）
        token := r.URL.Query().Get("token")
        if token == "" {
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }

        flusher, ok := w.(http.Flusher)
        if !ok {
            http.Error(w, "streaming unsupported", http.StatusInternalServerError)
            return
        }

        // headers
        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")
        // CORS if needed
        w.Header().Set("Access-Control-Allow-Origin", "*")

        lastEventID := r.Header.Get("Last-Event-ID")
        _ = lastEventID // 在需要时用来补发

        clientID := fmt.Sprintf("%s-%d", token, time.Now().UnixNano())
        client := &Client{
            id:   clientID,
            ch:   make(chan string, 16),
            done: make(chan struct{}),
        }
        h.register <- client
        defer func() { h.unregister <- client }()

        // send initial event
        fmt.Fprintf(w, "event: connected\ndata: %s\n\n", "ok")
        flusher.Flush()

        // heartbeat ticker to keep connection alive (some proxies close idle)
        heart := time.NewTicker(25 * time.Second)
        defer heart.Stop()

        ctx := r.Context()

        for {
            select {
            case <-ctx.Done():
                return
            case msg, ok := <-client.ch:
                if !ok {
                    return
                }
                // send message (as JSON/text)
                fmt.Fprintf(w, "data: %s\n\n", msg)
                flusher.Flush()
            case <-heart.C:
                // comment as heartbeat
                fmt.Fprintf(w, ": heartbeat\n\n")
                flusher.Flush()
            }
        }
    }
}
