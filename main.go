package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
)

func main() {
    hub := NewHub()
    go hub.Run()

    http.Handle("/", http.FileServer(http.Dir("./web")))
    http.HandleFunc("/sse", sseHandler(hub))
    http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
        // 简单 POST 或 GET 接口用于测试广播
        if r.Method == http.MethodPost {
            var payload map[string]interface{}
            if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
                http.Error(w, "invalid json", http.StatusBadRequest)
                return
            }
            bytes, _ := json.Marshal(payload)
            hub.Broadcast(string(bytes))
            w.WriteHeader(http.StatusNoContent)
            return
        }
        // GET: query msg
        msg := r.URL.Query().Get("msg")
        if msg == "" {
            msg = "hello"
        }
        hub.Broadcast(msg)
        w.WriteHeader(http.StatusNoContent)
    })

    addr := ":8080"
    if v := os.Getenv("PORT"); v != "" {
        addr = ":" + v
    }
    srv := &http.Server{
        Addr:         addr,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 0, // allow long-lived SSE
        IdleTimeout:  120 * time.Second,
    }
    log.Printf("listening on %s\n", addr)
    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
