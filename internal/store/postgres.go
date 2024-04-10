package store

import (
	"database/sql"
	"fmt"

	"github.com/jneurock/todo-go/internal/domain"

	_ "github.com/lib/pq"
)

type TodoPostgresStore struct {
	isAvailable bool
	db          *sql.DB
}

func NewTodoPostsgresStore(db *sql.DB, isAvailable bool) *TodoPostgresStore {
	store := &TodoPostgresStore{db: db, isAvailable: true}

	if err := db.Ping(); err != nil {
		fmt.Println("Could not ping database")
		fmt.Println(err)
		store.isAvailable = false
	}

	if err := store.init(); err != nil {
		fmt.Println("Could not initialize database")
		fmt.Println(err)
		store.isAvailable = false
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
	return s.isAvailable
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
	query := "SELECT * FROM todo ORDER BY \"complete\" ASC, \"id\" DESC"
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
