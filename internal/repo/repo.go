package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	Staff    *staffRepo
	Product  *productRepo
	Customer *customerRepo
	Checkout *checkoutRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.Staff = newStaffRepo(conn)
	repo.Product = newProductRepo(conn)
	// repo.Cat = newCatRepo(conn)
	// repo.Match = newMatchRepo(conn)
	repo.Customer = newCustomerRepo(conn)
	repo.Checkout = newCheckoutRepo(conn)

	return &repo
}
