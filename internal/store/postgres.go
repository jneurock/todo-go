package store

import (
	"database/sql"
	"fmt"
	"os"

	domain "github.com/jneurock/todo-go/internal/domain"

	_ "github.com/lib/pq"
)

type TodoPostgresStore struct {
	db *sql.DB
}

func NewTodoPostsgresStore(sslmode string) (*TodoPostgresStore, error) {
	user := os.Getenv("DB_USER")
	pw := os.Getenv("DB_PW")
	name := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", user, name, pw, sslmode)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &TodoPostgresStore{db: db}

	if err := store.init(); err != nil {
		return nil, err
	}

	return store, nil
}

func (r *TodoPostgresStore) init() error {
	query := `CREATE TABLE IF NOT EXISTS todo (
		id serial primary key,
		completed boolean,
		title text
	)`

	_, err := r.db.Exec(query)

	return err
}

func (r *TodoPostgresStore) Create(attrs TodoAttrs) (*domain.Todo, error) {
	query := `INSERT INTO todo (completed, title) values ($1, $2)`

	rows, err := r.db.Query(
		query,
		attrs.Complete,
		attrs.Description,
	)

	if err != nil {
		return nil, err
	}

	return scanIntoTodo(rows)
}

func (r *TodoPostgresStore) Delete(id string) error {
	query := `DELETE FROM todo WHERE id = $1`

	_, err := r.db.Exec(query, id)

	return err
}

func (r *TodoPostgresStore) Find(id string) (*domain.Todo, error) {
	query := "SELECT * FROM todo WHERE id = $1"

	rows, err := r.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoTodo(rows)
	}

	return nil, fmt.Errorf("todo not found")
}

func (r *TodoPostgresStore) FindAll() ([]*domain.Todo, error) {
	query := "SELECT * FROM todo"

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	todos := []*domain.Todo{}

	for rows.Next() {
		todo, err := scanIntoTodo(rows)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *TodoPostgresStore) Update(todo *domain.Todo) error {
	query := `UPDATE todo SET completed = $1, title = $2 WHERE id = $3`

	_, err := r.db.Exec(query, todo.Complete, todo.Description, todo.ID)

	return err
}

func scanIntoTodo(rows *sql.Rows) (*domain.Todo, error) {
	todo := new(domain.Todo)

	err := rows.Scan(
		&todo.ID,
		&todo.Complete,
		&todo.Description,
	)

	return todo, err
}
