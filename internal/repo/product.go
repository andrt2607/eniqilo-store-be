package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/ierr"

	validator "eniqilo-store-be/pkg/validator"

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

func (cr *productRepo) GetProductSKU(ctx context.Context, param dto.ReqParamProductSKUGet) ([]dto.ResProductSKUGet, error) {
	var query strings.Builder
	categoryProduct := []string{"Clothing", "Accessories", "Footwear", "Beverages"}

	query.WriteString("SELECT id, name, sku, category, imageUrl, stock, price, location, createdAt FROM product WHERE isAvailable = true ")

	if param.SKU != "" {
		query.WriteString(fmt.Sprintf("AND sku = '%s' ", param.SKU))
	}

	if param.Stock == "true" {
		query.WriteString("AND stock > 0 ")
	}

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	if validator.IsInArray(param.Category, categoryProduct) {
		query.WriteString(fmt.Sprintf("AND category = '%s') ", param.Category))
	}

	if param.CreatedAt == "asc" {
		query.WriteString("ORDER BY created_at ASC ")
	} else {
		query.WriteString("ORDER BY created_at DESC ")
	}

	if param.Price == "asc" {
		query.WriteString(", price ASC ")
	} else if param.Price == "desc" {
		query.WriteString(", price DESC ")
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	rows, err := cr.conn.Query(ctx, query.String()) // Replace $1 with sub
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "02000" {
				return []dto.ResProductSKUGet{}, nil
			}
		}
		return nil, err
	}
	defer rows.Close()

	results := []dto.ResProductSKUGet{}
	for rows.Next() {

		result := dto.ResProductSKUGet{}
		err := rows.Scan(
			&result.Id,
			&result.Name,
			&result.SKU,
			&result.Category,
			&result.ImageURL,
			&result.Stock,
			&result.Price,
			&result.SKU,
			&result.Location,
			&result.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
