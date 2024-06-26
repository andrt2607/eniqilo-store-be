package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		CreatedAt: insertedRow.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
	}, nil
}

func (p *productRepo) UpdateByID(ctx context.Context, productId string, product dto.ReqCreateProduct) (interface{}, error) {
	q := `UPDATE product SET name = $1, sku = $2, category = $3, imageUrl = $4, notes = $5, price = $6, stock = $7, location = $8, is_available = $9 WHERE id = $10 RETURNING id`

	var rowAffected string

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
		productId,
	).Scan(
		&rowAffected,
	)

	if err != nil {
		ierr.LogErrorWithLocation(err)
		if rowAffected == "" {
			return nil, ierr.ErrNotFound
		}
	}

	return struct{ ID string }{ID: rowAffected}, nil
}

func (p *productRepo) Get(ctx context.Context, param dto.ReqParamProductGet) ([]dto.ResGetProduct, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, sku, category, imageurl, notes, price, stock, location, is_available, created_at FROM product WHERE 1=1 ")

	if param.ID != "" {
		query.WriteString(fmt.Sprintf("AND id = '+%s' ", param.ID))
	}
	if param.Sku != "" {
		query.WriteString(fmt.Sprintf("AND sku = '%s' ", param.Sku))
	}
	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND name LIKE '%%%s%%' ", strings.ToLower(param.Name)))
	}
	if param.IsAvailable == "true" {
		query.WriteString(fmt.Sprintf("AND is_available = '%s' ", param.IsAvailable))
	} else if param.IsAvailable == "false" {
		query.WriteString(fmt.Sprintf("AND is_available = '%s' ", param.IsAvailable))
	}
	if param.Category != "" {
		fmt.Println("masuk", param.Category)
		switch dto.Category(param.Category) {
		case dto.Clothing:
			fmt.Println("masuk", dto.Clothing)
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		case dto.Accessories:
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		case dto.Footwear:
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		case dto.Beverages:
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		default:
		}
	}
	if param.InStock == "true" {
		query.WriteString("AND stock > 0 ")
	} else if param.InStock == "false" {
		query.WriteString("AND stock = 0 ")
	}

	fmt.Println("cek sini", param.Price)

	var orderByPrice bool
	if param.Price != "" {
		orderByPrice = true
		if param.Price == "asc" {
			query.WriteString("ORDER BY price ASC ")
		} else if param.Price == "desc" {
			query.WriteString("ORDER BY price DESC ")
		}
	}
	fmt.Println("cek sini", param.CreatedAt)
	if param.CreatedAt != "" {
		if orderByPrice {
			query.WriteString(", ")
		} else {
			query.WriteString("ORDER BY ")
		}
		if param.CreatedAt == "asc" {
			query.WriteString("created_at ASC ")
		}
	} else {
		if orderByPrice {
			query.WriteString(", created_at DESC ")
		} else {
			query.WriteString("ORDER BY created_at DESC ")
		}
	}

	// var orderByCreatedAt bool
	// if param.CreatedAt != "" {
	// 	orderByCreatedAt = true
	// 	if param.CreatedAt == "asc" {
	// 		query.WriteString("ORDER BY created_at ASC ")
	// 	} else {
	// 		query.WriteString("ORDER BY created_at DESC ")
	// 	}
	// }

	// if param.Price != "" {
	// 	if orderByCreatedAt {
	// 		query.WriteString(", ")
	// 	} else {
	// 		query.WriteString("ORDER BY ")
	// 	}
	// 	if param.Price == "asc" {
	// 		query.WriteString("price ASC ")
	// 	} else if param.Price == "asc" {
	// 		query.WriteString("price DESC ")
	// 	}
	// }

	// if param.CreatedAt == "asc" {
	// 	query.WriteString("ORDER BY created_at ASC ")
	// } else {
	// 	query.WriteString("ORDER BY created_at DESC ")
	// }

	// fmt.Println(param.Price)

	// if param.Price == "asc" {
	// 	query.WriteString(", price ASC ")
	// } else if param.Price == "desc" {
	// 	query.WriteString(", price DESC ")
	// }

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	fmt.Println("query: ", query.String())
	rows, err := p.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := []dto.ResGetProduct{}
	for rows.Next() {
		var product dto.ResGetProduct
		var createdAt time.Time
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Category,
			&product.ImageURL,
			&product.Notes,
			&product.Price,
			&product.Stock,
			&product.Location,
			&product.IsAvailable,
			&createdAt,
		)
		product.CreatedAt = createdAt.Format("2006-01-02T15:04:05.000Z")
		if err != nil {
			return nil, err
		}
		results = append(results, product)
	}
	return results, nil
}

func (cr *productRepo) GetProductSKU(ctx context.Context, param dto.ReqParamProductSKUGet) ([]dto.ResProductSKUGet, error) {
	var query strings.Builder
	// categoryProduct := []string{"Clothing", "Accessories", "Footwear", "Beverages"}

	query.WriteString(`SELECT id, name, sku, category, imageUrl, stock, price, location, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z') FROM product WHERE is_available = true `)

	if param.SKU != "" {
		query.WriteString(fmt.Sprintf("AND sku = '%s' ", param.SKU))
	}

	if param.Stock == "true" {
		query.WriteString("AND stock > 0 ")
	} else if param.Stock == "false" {
		query.WriteString("AND stock = 0 ")
	}

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	// if validator.IsInArray(param.Category, categoryProduct) {
	// 	query.WriteString(fmt.Sprintf("AND category = '%s') ", param.Category))
	// }
	fmt.Println(param.Category)
	if param.Category != "" {
		fmt.Println("masuk", param.Category)
		switch dto.Category(param.Category) {
		case dto.Clothing:
			fmt.Println("masuk", dto.Clothing)
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		case dto.Accessories:
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		case dto.Footwear:
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		case dto.Beverages:
			query.WriteString(fmt.Sprintf("AND category = '%s' ", param.Category))
		default:
		}
	}

	var orderByPrice bool
	fmt.Println("cek sini", param.Price)
	if param.Price != "" {
		orderByPrice = true
		if param.Price == "asc" {
			query.WriteString("ORDER BY price ASC ")
		} else if param.Price == "desc" {
			query.WriteString("ORDER BY price DESC ")
		}
	}

	fmt.Println("cek sini", param.CreatedAt)
	if param.CreatedAt != "" {
		if orderByPrice {
			query.WriteString(", ")
		} else {
			query.WriteString("ORDER BY ")
		}
		if param.CreatedAt == "asc" {
			query.WriteString("created_at ASC ")
		}
	} else {
		if orderByPrice {
			query.WriteString(", created_at DESC ")
		} else {
			query.WriteString("ORDER BY created_at DESC ")
		}
	}

	// var orderByCreatedAt bool
	// if param.CreatedAt != "" {
	// 	orderByCreatedAt = true
	// 	if param.CreatedAt == "asc" {
	// 		query.WriteString("ORDER BY created_at ASC ")
	// 	} else {
	// 		query.WriteString("ORDER BY created_at DESC ")
	// 	}
	// }

	// if param.Price != "" {
	// 	if orderByCreatedAt {
	// 		query.WriteString(", ")
	// 	} else {
	// 		query.WriteString("ORDER BY ")
	// 	}
	// 	if param.Price == "asc" {
	// 		query.WriteString("price ASC ")
	// 	} else if param.Price == "asc" {
	// 		query.WriteString("price DESC ")
	// 	}
	// }

	// if param.CreatedAt == "asc" {
	// 	query.WriteString("ORDER BY created_at ASC ")
	// } else {
	// 	query.WriteString("ORDER BY created_at DESC ")
	// }

	// if param.Price == "asc" {
	// 	query.WriteString(", price ASC ")
	// } else if param.Price == "desc" {
	// 	query.WriteString(", price DESC ")
	// }

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	fmt.Println("ini query", query.String())
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
			&result.Location,
			&result.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (p *productRepo) DeleteProduct(ctx context.Context, id string) error {
	// q := `UPDATE product SET deleted_at = now() WHERE id = $1 RETURNING id`

	q := `DELETE FROM product WHERE id = $1 RETURNING id`

	res, err := p.conn.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}
