package web

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/jneurock/todo-go/internal/store"
)

func newTestServer(t *testing.T, isStoreAvailable bool) *Server {
	store := store.NewTodoInMemoryStore(isStoreAvailable)
	server, err := NewServer(store, "views")

	if err != nil {
		t.Fatal(err)
	}

	return server
}

func testHandler(t *testing.T, r *http.Request, handler func(http.ResponseWriter, *http.Request)) {
	w := httptest.NewRecorder()

	handler(w, r)

	response := w.Result()

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}

	snaps.MatchSnapshot(t, string(b))
}

func Test500(t *testing.T) {
	server := newTestServer(t, false)
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	testHandler(t, r, server.createHandler(server.handleIndex))
}

func TestNotFound(t *testing.T) {
	server := newTestServer(t, true)
	r := httptest.NewRequest(http.MethodGet, "/doesnotexist", nil)

	testHandler(t, r, server.createHandler(server.handleIndex))
}

func TestIndex(t *testing.T) {
	server := newTestServer(t, true)

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	testHandler(t, r, server.createHandler(server.handleIndex))
}

func TestIndexSeeded(t *testing.T) {
	server := newTestServer(t, true)

	server.store.Create("Do chores")
	server.store.Create("Buy groceries")

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	testHandler(t, r, server.createHandler(server.handleIndex))
}

func TestNewTodo(t *testing.T) {
	server := newTestServer(t, true)
	body := strings.NewReader("description=Do%20chores")

	r := httptest.NewRequest(http.MethodPost, "/todo", body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testHandler(t, r, server.createHandler(server.handleNewTodo))
}

func TestDeleteTodo(t *testing.T) {
	server := newTestServer(t, true)
	err := server.store.Create("Delete me")

	if err != nil {
		t.Fatal(err)
	}

	todos, err := server.store.FindAll()

	if err != nil {
		t.Fatal(err)
	}

	id := strconv.FormatInt(todos[0].ID, 10)
	url := fmt.Sprintf("/todo/%s", id)

	r := httptest.NewRequest(http.MethodDelete, url, nil)
	r.SetPathValue("id", id)

	testHandler(t, r, server.createHandler(server.handleDeleteTodo))
}

func TestUpdateTodo(t *testing.T) {
	server := newTestServer(t, true)
	err := server.store.Create("Update me")

	if err != nil {
		t.Fatal(err)
	}

	todos, err := server.store.FindAll()

	if err != nil {
		t.Fatal(err)
	}

	id := strconv.FormatInt(todos[0].ID, 10)
	url := fmt.Sprintf("/todo/%s", id)
	body := strings.NewReader("complete=on&description=Update%20me%20again")

	r := httptest.NewRequest(http.MethodPut, url, body)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.SetPathValue("id", id)

	testHandler(t, r, server.createHandler(server.handleUpdateTodo))
}
