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

	// show query
	fmt.Println(query.String())

	rows, err := cr.conn.Query(ctx, query.String()) // Replace $1 with sub
	if err != nil {
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

// func (cr *customerRepo) GetCatByID(ctx context.Context, id, sub string) (dto.ResCustomerGet, error) {
// 	q := `SELECT id,
// 		name,
// 		race,
// 		sex,
// 		age_in_month,
// 		description,
// 		image_urls,
// 		EXISTS (
// 			SELECT 1 FROM match_customers m WHERE m.user_customer_id = c.id AND m.user_id = $1
// 		) AS has_matched,
// 		created_at
// 	FROM cats c WHERE id = $2`

// 	var imageUrl sql.NullString
// 	var createdAt int64
// 	var description string

// 	result := dto.ResCustomerGet{}
// 	err := cr.conn.QueryRow(ctx, q, sub, id).Scan(&result.ID, &result.Name, &result.Race, &result.Sex, &result.AgeInMonth, &description, &imageUrl, &result.HasMatched, &createdAt)
// 	if err != nil {
// 		return dto.ResCustomerGet{}, err
// 	}

// 	result.ImageUrls = strings.Split(imageUrl.String, ",")
// 	result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
// 	result.Description = description

// 	return result, nil
// }
