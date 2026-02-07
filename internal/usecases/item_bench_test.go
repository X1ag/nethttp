package usecases

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"nethttppractice/internal/domain"
	"strings"
	"testing"
	"time"
)

func BenchmarkItemHandler(b *testing.B) {
	mock := &mockRepo{
		items: []domain.Item{
			{ID: 1, Name: "Item 1", Description: "Desc 1", CreatedAt: time.Now()},
			{ID: 2, Name: "Item 2", Description: "Desc 2", CreatedAt: time.Now()},
		},
	}
	handler := NewItemHandler(mock)

	b.Run("GetItems", func(b *testing.B) {
		req := httptest.NewRequest("GET", "/items", nil)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w := httptest.NewRecorder()
			handler.GetItems(w, req)
		}
	})

	b.Run("InsertItem", func(b *testing.B) {
		item := domain.Item{Name: "Test", Description: "Bench"}
		jsonData, _ := json.Marshal(item)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("POST", "/items", bytes.NewReader(jsonData))
			w := httptest.NewRecorder()
			handler.InsertItem(w, req)
		}
	})

	b.Run("DeleteItem", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("DELETE", "/items", nil)
			req.SetPathValue("id", "1")
			w := httptest.NewRecorder()
			handler.DeleteItem(w, req)
		}
	})
}

func BenchmarkGetItemsWithDifferentSizes(b *testing.B) {
	scenarios := []struct {
		name      string
		itemCount int
	}{
		{"1_item", 1},
		{"10_items", 10},
		{"100_items", 100},
		{"1000_items", 1000},
	}

	for _, scenario := range scenarios {
		items := make([]domain.Item, scenario.itemCount)
		for i := 0; i < scenario.itemCount; i++ {
			items[i] = domain.Item{
				ID:          i,
				Name:        "Item",
				Description: "Description",
				CreatedAt:   time.Now(),
			}
		}

		mock := &mockRepo{items: items}
		handler := NewItemHandler(mock)

		b.Run(scenario.name, func(b *testing.B) {
			req := httptest.NewRequest("GET", "/items", nil)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				handler.GetItems(w, req)
			}
		})
	}
}

func BenchmarkInsertItemWithDifferentPayloads(b *testing.B) {
	scenarios := []struct {
		name      string
		payloadSize int
	}{
		{"short_name", 10},
		{"middle_name", 100},
		{"long_name", 1000},
	}

	for i, scenario := range scenarios {
		item := domain.Item{
				ID:          i,
				Name:        strings.Repeat("a", scenario.payloadSize),
				Description: "Description",
				CreatedAt:   time.Now(),
		}
		jsonData, _ := json.Marshal(item)

		mock := &mockRepo{}
		handler := NewItemHandler(mock)

		b.Run(scenario.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				req := httptest.NewRequest("POST", "/items", bytes.NewReader(jsonData))
				w := httptest.NewRecorder()
				handler.InsertItem(w, req)
			}
		})
	}
}
