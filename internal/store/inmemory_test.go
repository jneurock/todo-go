package store

import (
	"testing"
)

func TestInMemoryCreateAndFind(t *testing.T) {
	var id int64 = 0
	want := "Do chores"
	store := NewTodoInMemoryStore()

	err := store.Create(want)

	if err != nil {
		t.Fatal(err)
	}

	got, err := store.Find("0")

	if err != nil {
		t.Fatal(err)
	}

	if got.Complete {
		t.Fatal("expected new todo to be incomplete")
	}

	if got.Description != want {
		t.Fatalf("want %s, got %s", want, got.Description)
	}

	if got.ID != id {
		t.Fatalf("expected new todo with id: %d", id)
	}
}

func TestInMemoryDelete(t *testing.T) {
	store := NewTodoInMemoryStore()
	err := store.Create("Do chores")

	if err != nil {
		t.Fatal(err)
	}

	_ = store.Delete("0")
	_, err = store.Find("0")

	if err.Error() != "todo not found" {
		t.Fatal("expected delete todo to not be found")
	}
}

func TestInMemoryFindAll(t *testing.T) {
	store := NewTodoInMemoryStore()
	want := []string{"Do chores", "Go to the store"}

	_ = store.Create("Go to the store")
	_ = store.Create("Do chores")

	todos, err := store.FindAll()

	if err != nil {
		t.Fatal(err)
	}

	for i, todo := range todos {
		got := todo.Description

		if got != want[i] {
			t.Fatalf("want %s, got %s", want[i], got)
		}
	}
}

func TestInMemoryUpdate(t *testing.T) {
	store := NewTodoInMemoryStore()
	want := "Do more chores"

	err := store.Create("Do chores")

	if err != nil {
		t.Fatal(err)
	}

	err = store.Update("0", "Do more chores", true)

	if err != nil {
		t.Fatal(err)
	}

	got, err := store.Find("0")

	if err != nil {
		t.Fatal(err)
	}

	if got.Description != want {
		t.Fatalf("want %s, got %s", want, got.Description)
	}

	if !got.Complete {
		t.Fatal("expected todo to be complete")
	}
}
