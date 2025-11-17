package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Todo struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Category    string     `json:"category,omitempty"`
	Priority    int        `json:"priority,omitempty"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=1", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // sqlite best practice
	s := &Store{db: db}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) init() error {
	schema := `
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    category TEXT,
    priority INTEGER DEFAULT 0,
    due_at DATETIME,
    completed INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) CreateTodo(t *Todo) (int64, error) {
	now := time.Now().UTC()
	res, err := s.db.Exec(`INSERT INTO todos (title,description,category,priority,due_at,completed,created_at,updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		t.Title, t.Description, t.Category, t.Priority, t.DueAt, boolToInt(t.Completed), now, now)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func (s *Store) GetTodos() ([]Todo, error) {
	rows, err := s.db.Query(`SELECT id,title,description,category,priority,due_at,completed,created_at,updated_at FROM todos ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Todo
	for rows.Next() {
		var t Todo
		var due sql.NullTime
		var completed int
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Category, &t.Priority, &due, &completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if due.Valid {
			d := due.Time
			t.DueAt = &d
		}
		t.Completed = completed != 0
		out = append(out, t)
	}
	return out, nil
}

func (s *Store) GetTodoByID(id int64) (*Todo, error) {
	row := s.db.QueryRow(`SELECT id,title,description,category,priority,due_at,completed,created_at,updated_at FROM todos WHERE id = ?`, id)
	var t Todo
	var due sql.NullTime
	var completed int
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Category, &t.Priority, &due, &completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	if due.Valid {
		d := due.Time
		t.DueAt = &d
	}
	t.Completed = completed != 0
	return &t, nil
}

func (s *Store) UpdateTodo(t *Todo) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(`UPDATE todos SET title=?,description=?,category=?,priority=?,due_at=?,completed=?,updated_at=? WHERE id=?`,
		t.Title, t.Description, t.Category, t.Priority, t.DueAt, boolToInt(t.Completed), now, t.ID)
	return err
}

func (s *Store) DeleteTodo(id int64) error {
	_, err := s.db.Exec(`DELETE FROM todos WHERE id = ?`, id)
	return err
}

func (s *Store) SetCompleted(id int64, completed bool) error {
	now := time.Now().UTC()
	_, err := s.db.Exec(`UPDATE todos SET completed=?,updated_at=? WHERE id=?`, boolToInt(completed), now, id)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
