package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseIDFromPath(p string) (int64, bool) {
	// expects path like /api/todos/123 or /api/todos/123/complete
	parts := strings.Split(strings.Trim(p, "/"), "/")
	if len(parts) < 3 {
		return 0, false
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func todosHandler(store *Store, hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todos, err := store.GetTodos()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, todos)
			return
		case http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			var t Todo
			if err := json.Unmarshal(body, &t); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			if strings.TrimSpace(t.Title) == "" {
				http.Error(w, "title required", http.StatusBadRequest)
				return
			}
			id, err := store.CreateTodo(&t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// fetch the inserted row to get DB-populated fields (created_at/updated_at, etc)
			created, err := store.GetTodoByID(id)
			if err != nil {
				// if fetch fails, still return created id but indicate server error
				http.Error(w, "created but failed to fetch", http.StatusInternalServerError)
				return
			}

			// send SSE notification with the full created record
			evt := map[string]interface{}{"type": "todo.created", "todo": created}
			b, _ := json.Marshal(evt)
			hub.Broadcast(string(b))

			writeJSON(w, http.StatusCreated, created)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func todoDetailHandler(store *Store, hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// path examples:
		// /api/todos/123
		// /api/todos/123/complete (we support PATCH via this path)
		id, ok := parseIDFromPath(r.URL.Path)
		if !ok {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			t, err := store.GetTodoByID(id)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, t)
			return
		case http.MethodPut:
			body, _ := io.ReadAll(r.Body)
			var t Todo
			if err := json.Unmarshal(body, &t); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			t.ID = id
			if err := store.UpdateTodo(&t); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			updated, _ := store.GetTodoByID(id)
			evt := map[string]interface{}{"type": "todo.updated", "todo": updated}
			b, _ := json.Marshal(evt)
			hub.Broadcast(string(b))
			writeJSON(w, http.StatusOK, updated)
			return
		case http.MethodPatch:
			// support /api/todos/{id}/complete
			if strings.HasSuffix(r.URL.Path, "/complete") {
				// read body for { "completed": true }
				var payload map[string]bool
				body, _ := io.ReadAll(r.Body)
				if err := json.Unmarshal(body, &payload); err != nil {
					http.Error(w, "invalid json", http.StatusBadRequest)
					return
				}
				completed := payload["completed"]
				if err := store.SetCompleted(id, completed); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				t, _ := store.GetTodoByID(id)
				evt := map[string]interface{}{"type": "todo.completed", "todo": t}
				b, _ := json.Marshal(evt)
				hub.Broadcast(string(b))
				writeJSON(w, http.StatusOK, t)
				return
			}
			http.Error(w, "unsupported patch", http.StatusBadRequest)
			return
		case http.MethodDelete:
			if err := store.DeleteTodo(id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			evt := map[string]interface{}{"type": "todo.deleted", "id": id}
			b, _ := json.Marshal(evt)
			hub.Broadcast(string(b))
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}