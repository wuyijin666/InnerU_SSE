package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHub(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// 测试消息广播功能
	client1 := &Client{
		id:   "client1",
		ch:   make(chan string, 1),
		done: make(chan struct{}),
	}
	
	client2 := &Client{
		id:   "client2",
		ch:   make(chan string, 1),
		done: make(chan struct{}),
	}

	// 注册客户端
	hub.register <- client1
	hub.register <- client2
	
	// 等待注册完成
	time.Sleep(10 * time.Millisecond)
	
	// 广播消息
	hub.Broadcast("test message")
	
	// 检查两个客户端是否都收到了消息
	select {
	case msg := <-client1.ch:
		if msg != "test message" {
			t.Errorf("Expected 'test message', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Client1 did not receive message in time")
	}
	
	select {
	case msg := <-client2.ch:
		if msg != "test message" {
			t.Errorf("Expected 'test message', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Client2 did not receive message in time")
	}
	
	// 取消注册一个客户端
	hub.unregister <- client1
	time.Sleep(10 * time.Millisecond)
	
	// 再次广播
	hub.Broadcast("second message")
	
	// 现在只有 client2 应该收到消息
	select {
	case msg := <-client2.ch:
		if msg != "second message" {
			t.Errorf("Expected 'second message', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Error("Client2 did not receive second message in time")
	}
	
	// Client1 不应该收到任何消息
	select {
	case msg := <-client1.ch:
		t.Errorf("Client1 should not receive message but got: %s", msg)
	default:
		// 符合预期 - 没有收到消息
	}
}

func TestSSEHandler(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	
	handler := sseHandler(hub)
	
	// 测试缺少 token 的情况
	req := httptest.NewRequest("GET", "/sse", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
	
	// 测试有效请求
	req = httptest.NewRequest("GET", "/sse?token=test", nil)
	w = httptest.NewRecorder()
	
	// 使用短超时避免测试挂起
	go func() {
		time.Sleep(100 * time.Millisecond)
		handler(w, req)
	}()
	
	// 检查响应
	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", result.StatusCode)
	}
	
	contentType := result.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		t.Errorf("Expected text/event-stream content type, got %s", contentType)
	}
}