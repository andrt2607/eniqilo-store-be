package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/ierr"

	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	conn *pgxpool.Pool
}

type insertResult struct {
	ID        string
	CreatedAt time.Time
}

func newProductRepo(conn *pgxpool.Pool) *productRepo {
	return &productRepo{conn}
}

func (p *productRepo) Insert(ctx context.Context, product dto.ReqCreateProduct) (interface{}, error) {
	q := `INSERT INTO product (name, sku, category, imageUrl, notes, price, stock, location, is_available, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now()) RETURNING id, created_at`

	var insertedRow insertResult

	err := p.conn.QueryRow(ctx, q,
		product.Name,
		product.SKU,
		product.Category,
		product.ImageURL,
		product.Notes,
		product.Price,
		product.Stock,
		product.Location,
		product.IsAvailable,
	).Scan(
		&insertedRow.ID,
		&insertedRow.CreatedAt,
	)

	if err != nil {
		ierr.LogErrorWithLocation(err)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return nil, ierr.ErrDuplicate
			}
		}
		return nil, err
	}

	return dto.ResCreateProduct{
		ID:        insertedRow.ID,
		CreatedAt: insertedRow.CreatedAt.Format(time.RFC3339),
	}, nil
}
