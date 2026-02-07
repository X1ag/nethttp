package domain

import (
	"context"
	"time"
)

type Item struct {
	ID   int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	GetByID(ctx context.Context, id int) (*Item, error)
	GetAllItems(ctx context.Context) ([]Item, error)
	Create(ctx context.Context, u *Item) (*Item, error)
	Update(ctx context.Context, u *Item) error
	Delete(ctx context.Context, id int) error
}
