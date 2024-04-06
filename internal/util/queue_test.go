package util

import "testing"

type TestItem struct {}

func TestAdd(t *testing.T) {
	queue := &Queue[TestItem]{}
	want := &TestItem{}

	queue.Add(want)

	if queue.Head == nil {
		t.Fatalf("want %v, got nil", want)
	}

	if queue.Head.Value != want {
		t.Fatalf("want %v, got %v", want, queue.Head.Value)
	}
}

func TestAddMore(t *testing.T) {
	queue := &Queue[TestItem]{}
	want := &TestItem{}

	queue.Add(&TestItem{})
	queue.Add(want)

if queue.Head.Value != want {
		t.Fatalf("want %v, got %v", want, queue.Head.Value)
	}
}

func TestRemove(t *testing.T) {
	queue := &Queue[TestItem]{}
	want := &TestItem{}
	itemToRemove := &TestItem{}

	queue.Add(want)
	queue.Add(itemToRemove)

	if queue.Head.Value != itemToRemove {
		t.Fatal("head of queue is the wrong item")
	}

	queue.Remove(func (ti *TestItem) bool {
		return ti == itemToRemove
	})

	if queue.Head.Value != want {
		t.Fatalf("want %v, got %v", want, queue.Head.Value)
	}
}

func TestRemoveNoMatch(t *testing.T) {
	queue := &Queue[TestItem]{}
	wantHead := &TestItem{}
	wantTail := &TestItem{}
	
	queue.Add(wantTail)
	queue.Add(wantHead)

	queue.Remove(func (ti *TestItem) bool {
		return ti == nil
	})

	if queue.Head.Value != wantHead {
		t.Fatalf("want head %v, got %v", wantHead, queue.Head.Value)
	}

	if queue.Head.Next.Value != wantTail {
		t.Fatalf("want tail %v, got %v", wantTail, queue.Head.Next.Value)
	}
}
