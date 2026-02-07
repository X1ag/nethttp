package integration

import (
	"context"
	"nethttppractice/internal/domain"
	postgres "nethttppractice/internal/repository"
	"testing"
)

func TestRepository_Create(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewPgRepo(pool)

	item := &domain.Item{
		Name: "Test",
		Description: "Desc",
	}

	newItem, err := repo.Create(context.Background(), item)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if newItem.ID == 0 {
		t.Errorf("Expected id > 0, got %d", newItem.ID)
	}
	if newItem.Name != "Test" {
		t.Errorf("Expected name %s, got %s", "Test", newItem.Name)
	}
}

func TestRepository_GetAllItems(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewPgRepo(pool)

	item1 := &domain.Item{
		Name: "Test1",
		Description: "Desc",
	}
	item2 := &domain.Item{
		Name: "Test2",
		Description: "Desc",
	}

	_, err := repo.Create(context.Background(), item1)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	_, err = repo.Create(context.Background(), item2)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	items, err := repo.GetAllItems(context.Background())
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if items[0].Name != "Test1" {
		t.Errorf("Expected name %s, got %s", "Test1", items[0].Name)
	}
	if items[1].Name != "Test2" {
		t.Errorf("Expected name %s, got %s", "Test2", items[1].Name)
}}

func TestRepository_DeleteItem(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewPgRepo(pool)
	item1 := &domain.Item{
		Name: "Test1",
		Description: "Desc",
	}

	item, err := repo.Create(t.Context(), item1)
	if err != nil {
		t.Fatalf("cant create a item with fields %v with err=%v", item1, err)
	}

	items, err := repo.GetAllItems(t.Context())
	if items[0].Name != item.Name {
		t.Errorf("item was not found %v", item)
	}

	err = repo.Delete(t.Context(), item.ID)
	if err != nil {
		t.Errorf("cant delete that item %v", err)
	}

	items, err = repo.GetAllItems(t.Context())
	if len(items) != 0 {
		t.Errorf("delete was not executed %v", items)
	}
}