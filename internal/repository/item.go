package postgres

import (
	"context"
	"errors"

	"nethttppractice/internal/domain"

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

func (r *pgRepository) GetByID(ctx context.Context, id int) (*domain.Item, error) {
	query := `SELECT id, name FROM items WHERE id = $1`
	row := r.pg.QueryRow(ctx, query, id)
	item := &domain.Item{}
	if err := row.Scan(&item.ID, &item.Name); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *pgRepository) GetAllItems(ctx context.Context) ([]domain.Item, error) {
	query := `SELECT id, name, description FROM items`
	rows, err := r.pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Item, 0, 10)	

	i := domain.Item{}
	for rows.Next() {
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}	
		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *pgRepository) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	query := `INSERT INTO items (name, description, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := r.pg.QueryRow(
		ctx,
		query,
		item.Name,
		item.Description, 
		item.CreatedAt,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *pgRepository) Update(ctx context.Context, u *domain.Item) error {
	query := `UPDATE items SET name = $1, description = $2 WHERE id = $3`
	row, err := r.pg.Exec(ctx, query, u.Name, u.Description, u.ID)	
	if err != nil {
		return err
	}

	if row.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (r *pgRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM items WHERE id = $1`
	row, err := r.pg.Exec(ctx, query, id)	
	if err != nil {
		return err
	}

	if row.RowsAffected() == 0 {
		return errors.New("no rows deleted")
	}

	return nil
}