package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	// "log"
	"net/http"
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

func (cr *checkoutRepo) PostValidateCheckout(ctx context.Context, checkout dto.ReqCheckoutPost) (int, error) {
	q := `SELECT id FROM customer where id = $1`
	_, err := cr.conn.Query(ctx, q,
		checkout.CustomerId)

	if err != nil {
		err_message := fmt.Sprintf("customerId %v is not found", checkout.CustomerId)
		return http.StatusNotFound, errors.New(err_message)
	}

	// var TotalCharge *int
	// *TotalCharge = 0

	TotalCharge := 0

	for _, element := range checkout.ProductDetails {
		// productId, err := strconv.Atoi(element.ProductId)
		data := dto.ResValidateCheckoutPost{}
		// data.Charge
		q := `SELECT price*$2 as charge, stock-$2 as stock FROM product where id = $1 and is_available = true`

		err := cr.conn.QueryRow(ctx, q,
			element.ProductId, element.Quantity).Scan(&data.Charge, &data.Stock)

		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == "22P02" {
					return http.StatusNotFound, ierr.ErrNotFound
				}
			}
			return http.StatusBadRequest, err
		}

		if *data.Stock < 0 {
			err_message := fmt.Sprintf("one of productIds (%v) stock is not enough (stock %v need %v)", element.ProductId, data.Stock, element.Quantity)
			return http.StatusBadRequest, errors.New(err_message)
		}

		TotalCharge = TotalCharge + *data.Charge
		if TotalCharge > *checkout.Paid {
			err_message := fmt.Sprintf("paid %v is not enough based on all bought product %v", checkout.Paid, TotalCharge)
			return http.StatusBadRequest, errors.New(err_message)
		}
	}

	if (*checkout.Paid - TotalCharge) != *checkout.Change {
		err_message := fmt.Sprintf("change %v is not right, based on all bought product %v, and what is paid %v", checkout.Change, TotalCharge, checkout.Paid)
		return http.StatusBadRequest, errors.New(err_message)
	}

	return http.StatusContinue, nil
}

func (cr *checkoutRepo) PostCheckout(ctx context.Context, checkout dto.ReqCheckoutPost) (int, dto.ResCreateProduct, error) {

	qiProduct := `INSERT INTO order_product (customer_id, paid, change, created_at)
	VALUES ($1, $2, $3, now()) RETURNING id, created_at`

	var OrderID string
	var CreatedAt time.Time
	err := cr.conn.QueryRow(ctx, qiProduct,
		checkout.CustomerId, checkout.Paid, checkout.Change).Scan(&OrderID, &CreatedAt)

	if err != nil {
		ierr.LogErrorWithLocation(err)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return http.StatusBadRequest, dto.ResCreateProduct{}, ierr.ErrDuplicate
			}
		}
		return http.StatusBadRequest, dto.ResCreateProduct{}, err
	}

	for _, element := range checkout.ProductDetails {

		qiDetail := `INSERT INTO order_detail (order_id, product_id, product_order_quantity) VALUES($1, $2, $3)`
		_, err := cr.conn.Exec(ctx, qiDetail,
			OrderID, element.ProductId, element.Quantity)

		if err != nil {
			return http.StatusInternalServerError, dto.ResCreateProduct{}, err
		}

		qu := `UPDATE product SET stock = stock - $2, updated_at = now() WHERE id = $1`
		_, err = cr.conn.Exec(ctx, qu,
			element.ProductId, element.Quantity)

		if err != nil {
			return http.StatusInternalServerError, dto.ResCreateProduct{}, err
		}
	}

	return http.StatusOK, dto.ResCreateProduct{
		ID:        OrderID,
		CreatedAt: CreatedAt.Format(time.RFC3339),
	}, nil
}

func (cr *checkoutRepo) GetCheckout(ctx context.Context, param dto.ReqParamCheckoutGet) ([]dto.ResCheckoutGet, error) {
	var query strings.Builder

	query.WriteString("SELECT id, customer_id, paid, change FROM order_product WHERE 1=1 ")

	if param.CustomerId != "" {
		query.WriteString(fmt.Sprintf("AND customer_id = '+%s' ", param.CustomerId))
	}

	if param.CreatedAt == "asc" {
		query.WriteString("ORDER BY created_at ASC ")
	} else {
		query.WriteString("ORDER BY created_at DESC ")
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	rows, err := cr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []dto.ResCheckoutGet{}
	for rows.Next() {

		result := dto.ResCheckoutGet{}
		err := rows.Scan(
			&result.OrderId,
			&result.CustomerId,
			&result.Paid,
			&result.Change)
		if err != nil {
			return nil, err
		}

		orderDetail, err := cr.GetOrderDetail(ctx, result.OrderId)
		if err != nil {
			return nil, err
		}
		result.ProductDetails = orderDetail

		results = append(results, result)
	}

	return results, nil
}

func (cr *checkoutRepo) GetOrderDetail(ctx context.Context, orderId string) ([]entity.CheckoutDetail, error) {
	var query strings.Builder

	query.WriteString("SELECT product_id, product_order_quantity FROM order_detail WHERE order_id = $1")

	rows, err := cr.conn.Query(ctx, query.String(), orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []entity.CheckoutDetail{}
	for rows.Next() {

		result := entity.CheckoutDetail{}
		err := rows.Scan(
			&result.ProductId,
			&result.Quantity)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}
