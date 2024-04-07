package domain

import (
	"testing"
)

func TestNewDescription(t *testing.T) {
	want := "Do chores"
	got, err := NewDescription(want)

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestNewDescriptionBad(t *testing.T) {
	want := "description cannot be empty"
	_, err := NewDescription("")

	if err == nil {
		t.Fatal("expected NewDescription to return error")
	}

	got := err.Error()

	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}
