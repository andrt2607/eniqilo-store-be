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

type checkoutRepo struct {
	conn *pgxpool.Pool
}

func newCheckoutRepo(conn *pgxpool.Pool) *checkoutRepo {
	return &checkoutRepo{conn}
}

func (u *checkoutRepo) Insert(ctx context.Context, checkout entity.Checkout) (string, error) {

	q := `INSERT INTO checkout (name, phone_number, created_at)
	VALUES ($1, $2, now()) RETURNING id, phone_number, name `

	var checkoutID, checkoutPhoneNumber, checkoutName string
	err := u.conn.QueryRow(ctx, q,
		checkout.Name, checkout.PhoneNumber).Scan(&checkoutID, &checkoutPhoneNumber, &checkoutName)

	if err != nil {
		ierr.LogErrorWithLocation(err)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", ierr.ErrDuplicate
			}
		}
		return "", err
	}

	return checkoutID, nil
}

func (cr *checkoutRepo) GetCheckout(ctx context.Context, param dto.ReqParamCheckoutGet, sub string) ([]dto.ResCheckoutGet, error) {
	var query strings.Builder

	query.WriteString("SELECT id, phone_number, name FROM checkout WHERE 1=1 ")

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

	results := make([]dto.ResCheckoutGet, 0, 10)
	for rows.Next() {

		result := dto.ResCheckoutGet{}
		err := rows.Scan(
			&result.CheckoutID,
			&result.PhoneNumber,
			&result.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// func (cr *checkoutRepo) GetCatByID(ctx context.Context, id, sub string) (dto.ResCheckoutGet, error) {
// 	q := `SELECT id,
// 		name,
// 		race,
// 		sex,
// 		age_in_month,
// 		description,
// 		image_urls,
// 		EXISTS (
// 			SELECT 1 FROM match_checkouts m WHERE m.user_checkout_id = c.id AND m.user_id = $1
// 		) AS has_matched,
// 		created_at
// 	FROM cats c WHERE id = $2`

// 	var imageUrl sql.NullString
// 	var createdAt int64
// 	var description string

// 	result := dto.ResCheckoutGet{}
// 	err := cr.conn.QueryRow(ctx, q, sub, id).Scan(&result.ID, &result.Name, &result.Race, &result.Sex, &result.AgeInMonth, &description, &imageUrl, &result.HasMatched, &createdAt)
// 	if err != nil {
// 		return dto.ResCheckoutGet{}, err
// 	}

// 	result.ImageUrls = strings.Split(imageUrl.String, ",")
// 	result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
// 	result.Description = description

// 	return result, nil
// }
