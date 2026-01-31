package domain

import (
	"context"
)

type Item struct {
	ID   int
	Name string
	Description string
}

type Repository interface {
	GetByID(ctx context.Context, id int) (*Item, error)
	GetAllItems(ctx context.Context) ([]Item, error)
	Create(ctx context.Context, u *Item) error
	Update(ctx context.Context, u *Item) error
	Delete(ctx context.Context, id int) error
}
