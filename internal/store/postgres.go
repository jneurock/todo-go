package todo

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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
		complete boolean,
		description text,
		created_at timestamp,
		updated_at timestamp
	)`

	_, err := r.db.Exec(query)

	return err
}

func (r *TodoPostgresStore) Create(attrs TodoAttrs) (*domain.Todo, error) {
	query := `INSERT INTO todo
		(complete, description, created_at, updated_at)
		values ($1, $2, $3, $4)
	`

	rows, err := r.db.Query(
		query,
		attrs.Complete,
		attrs.Description,
		time.Now(),
		time.Now(),
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
	query := `UPDATE todo
		SET complete = $1, description = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(query, todo.Complete, todo.Description, time.Now(), todo.ID)

	return err
}

func scanIntoTodo(rows *sql.Rows) (*domain.Todo, error) {
	todo := new(domain.Todo)

	err := rows.Scan(
		&todo.ID,
		&todo.Complete,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	return todo, err
}
