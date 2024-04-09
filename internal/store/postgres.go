package store

import (
	"database/sql"

	"github.com/jneurock/todo-go/internal/domain"

	_ "github.com/lib/pq"
)

type TodoPostgresStore struct {
	available bool
	db        *sql.DB
}

// TODO: Let this accept a database
// TODO: Don't depend on sql.DB directly
func NewTodoPostsgresStore(connStr string) *TodoPostgresStore {
	db, err := sql.Open("postgres", connStr)

	store := &TodoPostgresStore{db: db, available: err != nil}

	if err != nil {
		return store
	}

	if err := db.Ping(); err != nil {
		store.available = false
	}

	if err := store.init(); err != nil {
		store.available = false
	}

	return store
}

func (s *TodoPostgresStore) init() error {
	query := `CREATE TABLE IF NOT EXISTS todo (
		id serial primary key,
		complete boolean,
		description text
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *TodoPostgresStore) IsAvailable() bool {
	return s.available
}

func (s *TodoPostgresStore) Create(description string) error {
	todo := domain.NewTodo(description)
	query := "INSERT INTO todo (complete, description) values ($1, $2)"

	_, err := s.db.Query(
		query,
		todo.Complete,
		todo.Description,
	)

	return err
}

func (s *TodoPostgresStore) Delete(id string) error {
	query := "DELETE FROM todo WHERE id = $1"
	_, err := s.db.Exec(query, id)

	return err
}

func (s *TodoPostgresStore) Find(id string) (*domain.Todo, error) {
	query := "SELECT * FROM todo WHERE id = $1"
	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	todo := new(domain.Todo)

	for rows.Next() {
		err := rows.Scan(
			&todo.ID,
			&todo.Complete,
			&todo.Description,
		)

		if err != nil {
			return nil, err
		}
	}

	return todo, nil
}

func (s *TodoPostgresStore) FindAll() ([]*domain.Todo, error) {
	query := "SELECT * FROM todo"
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	todos := []*domain.Todo{}

	for rows.Next() {
		todo := new(domain.Todo)

		err := rows.Scan(
			&todo.ID,
			&todo.Complete,
			&todo.Description,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *TodoPostgresStore) Update(id, description string, complete bool) error {
	query := "UPDATE todo SET complete = $1, description = $2 WHERE id = $3"
	_, err := s.db.Exec(query, complete, description, id)

	return err
}
