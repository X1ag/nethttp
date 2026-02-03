package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"nethttppractice/internal/domain"
	"strings"
	"testing"
	"time"
)

type mockRepo struct {
	items []domain.Item
	err error
}

func (m *mockRepo) GetAllItems(ctx context.Context) ([]domain.Item, error) {
	return m.items, m.err
}

func(m *mockRepo) Create(ctx context.Context, u *domain.Item) (*domain.Item, error) {
	if m.err != nil {
		return nil, m.err
	}
	u.ID = len(m.items)
	u.CreatedAt = time.Now()
	m.items = append(m.items, *u) 
	return u, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id int) (*domain.Item, error) {
	for _, item := range m.items {
		if item.ID == id	{
			return &item, nil
		} 
	}
	m.err = errors.New("not found")
	return nil, errors.New("not found") 
}

func (m *mockRepo) Update(ctx context.Context, u *domain.Item) error {
	for i, item := range m.items {
		if item.ID == u.ID {
			m.items[i] = *u
			return nil
		}
	}	
	m.err = errors.New("not found")
	return errors.New("not found")
}

func (m *mockRepo) Delete(ctx context.Context, id int) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	m.err = errors.New("not found")
	return errors.New("not found")
}

func TestGetItem(t *testing.T) {
	mock := &mockRepo{
		items: []domain.Item{
			{ID: 1, Name: "Item 1", Description: "Description 1", CreatedAt: time.Now()},
		},
	}

	req := httptest.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()

	h := NewItemHandler(mock)
	h.GetItems(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if w.Header().Get("content-type") != "application/json" {
		t.Errorf("Expected content-type %s, got %s", "application/json", w.Header().Get("content-type"))
	}

	item := []domain.Item{}
	body := w.Body.Bytes()

	if err := json.Unmarshal(body, &item); err != nil {
		t.Errorf("Expected body %s, got %s", "[]", body)
	}

	if item[0].Name != mock.items[0].Name {
		t.Errorf("Expected body %v, got %v", mock.items[0], item[0])
	}
}

func TestGetItem_DBError(t *testing.T) {
	mock := &mockRepo{
		err: errors.New("db error"),
	}

	req := httptest.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()

	h := NewItemHandler(mock)
	h.GetItems(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "db error") {
		t.Errorf("Expected error message, got %s", body)
	}
}

func TestGetItem_EmptyBody(t *testing.T) {
	mock := &mockRepo{
		items: []domain.Item{},
	}

	req := httptest.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()

	h := NewItemHandler(mock)
	h.GetItems(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()

	if !strings.Contains(body, "[]") {
		t.Errorf("Expected body %s, got %s", "[]", body)
	}
}

func TestInsertItems(t *testing.T) {
	mock := &mockRepo{}
	handler := NewItemHandler(mock)

	item := domain.Item{Name: "Test", Description: "Desc"}
	jsonItem, err := json.Marshal(item)
	if err != nil {
		t.Errorf("Error while parsing %v", item)
	}

	req := httptest.NewRequest("POST", "/items", bytes.NewReader(jsonItem))
	w := httptest.NewRecorder()

	handler.InsertItem(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, god %d", w.Code)
	}

	var response domain.Item 
	json.NewDecoder(w.Body).Decode(&response)

	if response.Name != item.Name {
		t.Errorf("expected name %s, got %s", item.Name, response.Name)
	}
}	

func TestInsert_InvalidJSON(t *testing.T) {
	mock := &mockRepo{}
	handler := NewItemHandler(mock)


	req := httptest.NewRequest("POST", "/items", strings.NewReader("invalid json"))
	w := httptest.NewRecorder()

	handler.InsertItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "body error") {
		t.Errorf("expected error message, got %s", body)
	}
}

func TestInsert_DBError(t *testing.T) {
	mock := &mockRepo{
		err: errors.New("db error"),
	}
	handler := NewItemHandler(mock)

	item := domain.Item{Name: "Test", Description: "Desc"}
	jsonItem, err := json.Marshal(item)
	if err != nil {
		t.Errorf("Error while parsing %v", item)
	}

	req := httptest.NewRequest("POST", "/items", bytes.NewReader(jsonItem))
	w := httptest.NewRecorder()

	handler.InsertItem(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "db error") {
		t.Errorf("expected error message, got %s", body)
	}
}
func TestDeleteItem(t *testing.T) {
	mock := &mockRepo{
		items: []domain.Item{
			{ID: 1, Name: "Item 1", Description: "Description 1", CreatedAt: time.Now()},
		},
	}

	req := httptest.NewRequest("DELETE", "/items", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h := NewItemHandler(mock)
	h.DeleteItem(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if w.Header().Get("content-type") != "application/json" {
		t.Errorf("Expected content-type %s, got %s", "application/json", w.Header().Get("content-type"))
	}
}	

func TestDeleteItem_InvalidID(t *testing.T) {
	mock := &mockRepo{
		items: []domain.Item{
			{ID: 1, Name: "Item 1", Description: "Description 1", CreatedAt: time.Now()},
		},
	}

	req := httptest.NewRequest("DELETE", "/items", nil)
	w := httptest.NewRecorder()

	h := NewItemHandler(mock)
	h.DeleteItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "invalid id") {
		t.Errorf("Expected error message, got %s", body)
	}
}

func TestDeleteItem_DBError(t *testing.T) {
	mock := &mockRepo{
		items: []domain.Item{
			{ID: 1, Name: "Item 1", Description: "Description 1", CreatedAt: time.Now()},
		},
		err: errors.New("db error"),
	}

	req := httptest.NewRequest("DELETE", "/items", nil)
	req.SetPathValue("id", "3")
	w := httptest.NewRecorder()

	h := NewItemHandler(mock)
	h.DeleteItem(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "db error") {
		t.Errorf("Expected error message, got %s", body)
	}
}