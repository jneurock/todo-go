package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewTestPostgresStore(t *testing.T) (TodoStore, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}

	return &TodoPostgresStore{db: db}, db, mock
}

func TestPostgresCreate(t *testing.T) {
	store, db, mock := NewTestPostgresStore(t)

	defer db.Close()

	q := mock.ExpectQuery("INSERT INTO todo \\(complete, description\\) values \\(\\$1, \\$2\\)")
	q.WithArgs(false, "Do chores")
	q.WillReturnRows(sqlmock.NewRows([]string{"id", "complete", "description"}))

	if err := store.Create("Do chores"); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Create failed with error: %s", err.Error())
	}
}

func TestPostgresDelete(t *testing.T) {
	store, db, mock := NewTestPostgresStore(t)

	defer db.Close()

	q := mock.ExpectExec("DELETE FROM todo WHERE id = \\$1")
	q.WithArgs("0")
	q.WillReturnResult(sqlmock.NewResult(1, 1))

	if err := store.Delete("0"); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Delete failed with error: %s", err.Error())
	}
}

func TestPostgresFind(t *testing.T) {
	store, db, mock := NewTestPostgresStore(t)

	defer db.Close()

	q := mock.ExpectQuery("SELECT \\* FROM todo WHERE id = \\$1")
	q.WithArgs("0")
	q.WillReturnRows(sqlmock.NewRows([]string{"id", "complete", "description"}))

	if _, err := store.Find("0"); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Find failed with error: %s", err.Error())
	}
}

func TestPostgresFindAll(t *testing.T) {
	store, db, mock := NewTestPostgresStore(t)

	defer db.Close()

	q := mock.ExpectQuery("SELECT \\* FROM todo")
	q.WillReturnRows(sqlmock.NewRows([]string{"id", "complete", "description"}))

	if _, err := store.FindAll(); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Find all failed with error: %s", err.Error())
	}
}

func TestPostgresUpdate(t *testing.T) {
	store, db, mock := NewTestPostgresStore(t)

	defer db.Close()

	q := mock.ExpectExec("UPDATE todo SET complete = \\$1, description = \\$2 WHERE id = \\$3")
	q.WithArgs(true, "Do more chores", "0")
	q.WillReturnResult(sqlmock.NewResult(1, 1))

	if err := store.Update("0", "Do more chores", true); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Update failed with error: %s", err.Error())
	}
}
