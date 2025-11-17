# InnerU_SSE

A small TODO server with SSE (Server-Sent Events) notifications and a minimal front-end demo.

## Overview

This project provides:
- A Go HTTP server that implements a simple TODO API (CRUD) backed by SQLite.
- SSE endpoint to broadcast changes (todo.created, todo.updated, todo.completed, todo.deleted) to connected clients.
- A minimal front-end demo at `web/index.html` to connect and view real-time notifications.

## Requirements

- Go 1.24+ (development & CI)
- (Optional) curl or PowerShell for testing

## Build

From the project root:

# build executable
go build -o ./bin/todo-sse.exe

# or run directly
go run .

## Run

Run in the foreground to see logs:

./bin/todo-sse.exe

The server listens on port 8080 by default. Open http://localhost:8080 in your browser.

## Front-end demo

Open in your browser:

http://localhost:8080/index.html

- Enter a token (for demo use `demo-token`) and click Connect.
- The page will show connection status and a log area where SSE messages appear.

## API examples (PowerShell)

Use `Invoke-RestMethod` in PowerShell (recommended on Windows). Replace IDs as needed.

# Create
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/todos" -ContentType "application/json" -Body '{"title":"Buy milk","description":"2L"}'

# List
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/todos"

# Update (PUT)
Invoke-RestMethod -Method Put -Uri "http://localhost:8080/api/todos/1" -ContentType "application/json" -Body '{"title":"Buy milk (2L)","description":"low-fat"}'

# Mark completed (PATCH to /complete)
Invoke-RestMethod -Method Patch -Uri "http://localhost:8080/api/todos/1/complete" -ContentType "application/json" -Body '{"completed":true}'

# Delete
Invoke-RestMethod -Method Delete -Uri "http://localhost:8080/api/todos/1"

## Testing SSE from terminal (curl)

# Listen for SSE stream (use curl.exe on Windows)
curl.exe -N "http://localhost:8080/sse?token=demo-token"

# Then in another terminal, create a todo to see the SSE notification
curl.exe -X POST -H "Content-Type: application/json" -d '{"title":"Test SSE"}' "http://localhost:8080/api/todos"


## Notes and Best Practices

- Do NOT commit built binaries or local database files. Add them to `.gitignore` (see `.gitignore` in repo). If you accidentally committed them, remove from git index with `git rm --cached`.
- For production use, replace the demo token mechanism with proper authentication and consider a more robust DB (Postgres) or a shared pub/sub broker for multi-instance SSE.
- The server writes to a local SQLite file; SQLite allows only one concurrent writer — the app config should set `db.SetMaxOpenConns(1)`.

## Development

- Run `gofmt -w .`, `go vet ./...`, and (recommended) `staticcheck ./...`.
- Run unit tests: `go test ./... -v`.

## Demo script

You can create a small PowerShell demo script to start the server and create a sample todo; see `demo.ps1` (not included). 

## License

MIT