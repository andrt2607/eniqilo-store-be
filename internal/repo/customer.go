package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"eniqilo-store-be/internal/dto"
	"eniqilo-store-be/internal/entity"
	"eniqilo-store-be/internal/ierr"

	"github.com/jackc/pgx/v5/pgxpool"
)

type customerRepo struct {
	conn *pgxpool.Pool
}

func newCustomerRepo(conn *pgxpool.Pool) *customerRepo {
	return &customerRepo{conn}
}

func (u *customerRepo) Insert(ctx context.Context, customer entity.Customer) (string, error) {

	q := `INSERT INTO customer (name, phone_number, created_at)
	VALUES ($1, $2, now()) RETURNING id, phone_number, name `

	var customerID, customerPhoneNumber, customerName string
	err := u.conn.QueryRow(ctx, q,
		customer.Name, customer.PhoneNumber).Scan(&customerID, &customerPhoneNumber, &customerName)

	if err != nil {
		ierr.LogErrorWithLocation(err)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", ierr.ErrDuplicate
			}
		}
		return "", err
	}

	return customerID, nil
}

func (cr *customerRepo) GetCustomer(ctx context.Context, param dto.ReqParamCustomerGet) ([]dto.ResCustomerGet, error) {
	var query strings.Builder

	query.WriteString("SELECT id, phone_number, name FROM customer WHERE 1=1 ")

	if param.PhoneNumber != "" {
		query.WriteString(fmt.Sprintf("AND phone_number = '+%s' ", param.PhoneNumber))
	}

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	query.WriteString("ORDER BY created_at DESC")

	rows, err := cr.conn.Query(ctx, query.String()) // Replace $1 with sub
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "02000" {
				return []dto.ResCustomerGet{}, nil
			}
		}
		return nil, err
	}
	defer rows.Close()

	results := []dto.ResCustomerGet{}
	for rows.Next() {

		result := dto.ResCustomerGet{}
		err := rows.Scan(
			&result.CustomerID,
			&result.PhoneNumber,
			&result.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
