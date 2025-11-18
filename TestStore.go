package main

import (
	"os"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	tempPath := "./test_store.db"
	defer os.Remove(tempPath)

	store, err := NewStore(tempPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	if store == nil {
		t.Fatal("Store should not be nil")
	}

	if store.db == nil {
		t.Error("Database connection should be established")
	}
}

func TestTodoCRUD(t *testing.T) {
	tempPath := "./test_crud.db"
	defer os.Remove(tempPath)
	
	store, err := NewStore(tempPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// 测试 CreateTodo
	dueTime := time.Now().UTC().Add(24 * time.Hour)
	todo := &Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		Category:    "test",
		Priority:    1,
		DueAt:       &dueTime,
		Completed:   false,
	}

	id, err := store.CreateTodo(todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}
	
	if id <= 0 {
		t.Error("Created todo should have positive ID")
	}

	// 测试 GetTodoByID
	retrievedTodo, err := store.GetTodoByID(id)
	if err != nil {
		t.Fatalf("Failed to get todo by ID: %v", err)
	}

	// 验证所有字段
	if retrievedTodo.Title != todo.Title {
		t.Errorf("Expected title %s, got %s", todo.Title, retrievedTodo.Title)
	}
	
	if retrievedTodo.Description != todo.Description {
		t.Errorf("Expected description %s, got %s", todo.Description, retrievedTodo.Description)
	}
	
	if retrievedTodo.Category != todo.Category {
		t.Errorf("Expected category %s, got %s", todo.Category, retrievedTodo.Category)
	}
	
	if retrievedTodo.Priority != todo.Priority {
		t.Errorf("Expected priority %d, got %d", todo.Priority, retrievedTodo.Priority)
	}

	// 测试 GetTodos
	todos, err := store.GetTodos()
	if err != nil {
		t.Fatalf("Failed to get todos: %v", err)
	}
	
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}

	// 测试 UpdateTodo
	retrievedTodo.Title = "Updated Title"
	err = store.UpdateTodo(retrievedTodo)
	if err != nil {
		t.Fatalf("Failed to update todo: %v", err)
	}

	// 测试 DeleteTodo
	err = store.DeleteTodo(id)
	if err != nil {
		t.Fatalf("Failed to delete todo: %v", err)
	}

	_, err = store.GetTodoByID(id)
	if err == nil {
		t.Error("Expected error when getting deleted todo")
	}
}

func TestBoolToInt(t *testing.T) {
	if boolToInt(true) != 1 {
		t.Error("boolToInt(true) should return 1")
	}
	
	if boolToInt(false) != 0 {
		t.Error("boolToInt(false) should return 0")
	}
}