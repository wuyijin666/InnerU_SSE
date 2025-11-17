# build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY *.go ./
COPY web ./web
RUN go build -o /todo-sse

# final
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /todo-sse /todo-sse
COPY --from=builder /app/web /web
EXPOSE 8080
ENTRYPOINT [""/todo-sse""]
