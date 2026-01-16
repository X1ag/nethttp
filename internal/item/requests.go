package itemrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepository struct {
	pg *pgxpool.Pool
}

func NewPgRepo(pg *pgxpool.Pool) *pgRepository {
	return &pgRepository{
		pg: pg,
	}
}

func (r *pgRepository) GetByID(ctx context.Context, id int) (*Item, error) {
	row := r.pg.QueryRow(ctx, "SELECT id, name FROM items WHERE id = $1", id)
	i := new(Item)
	if err := row.Scan(&i.ID, &i.Name); err != nil {
		return nil, err
	}
	return i, nil
}

func (r *pgRepository) GetAllItems(ctx context.Context) ([]Item, error) {
	rows, err := r.pg.Query(ctx, "SELECT id, name FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0, 10)	
	for rows.Next() {
		i := new(Item)
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, *i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *pgRepository) Create(ctx context.Context, u *Item) error {
	cnt, err := r.pg.Exec(ctx, "INSERT INTO items (id, name) VALUES ($1, $2)", u.ID, u.Name)
	if err != nil {
		return err
	}
	if cnt.RowsAffected()== 0 {
		return errors.New("no rows inserted")
	}

	return nil
}

func (r *pgRepository) Update(ctx context.Context, u *Item) error {
	return nil
}

func (r *pgRepository) Delete(ctx context.Context, id int) error {
	return nil
}