package repo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"

	"eniqilo-store-be/internal/entity"
	"eniqilo-store-be/internal/ierr"

	"github.com/jackc/pgx/v5/pgxpool"
)

type staffRepo struct {
	conn *pgxpool.Pool
}

func newStaffRepo(conn *pgxpool.Pool) *staffRepo {
	return &staffRepo{conn}
}

func (u *staffRepo) Insert(ctx context.Context, staff entity.Staff) (string, error) {
	credVal := staff.PhoneNumber
	q := `INSERT INTO staff (name, phone_number, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, now(), now()) RETURNING id`

	var staffID string
	err := u.conn.QueryRow(ctx, q,
		staff.Name, credVal, staff.Password).Scan(&staffID)

	if err != nil {
		ierr.LogErrorWithLocation(err)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", ierr.ErrDuplicate
			}
		}
		return "", err
	}

	return staffID, nil
}

func (u *staffRepo) GetByPhoneNumber(ctx context.Context, cred string) (entity.Staff, error) {
	staff := entity.Staff{}
	q := `SELECT id, name, phone_number, password_hash FROM staff
	WHERE phone_number = $1`

	var phoneNumber sql.NullString

	err := u.conn.QueryRow(ctx,
		q, cred).Scan(&staff.ID, &staff.Name, &phoneNumber, &staff.Password)

	staff.PhoneNumber = phoneNumber.String

	if err != nil {
		return staff, err
	}

	return staff, nil
}

func (u *staffRepo) GetByID(ctx context.Context, id string) (entity.Staff, error) {
	staff := entity.Staff{}
	q := `SELECT phone_number, name, password FROM staffs
	WHERE id = $1`

	var phoneNumber sql.NullString

	err := u.conn.QueryRow(ctx,
		q, id).Scan(&phoneNumber, &staff.Name, &staff.Password)

	staff.PhoneNumber = phoneNumber.String

	if err != nil {
		if err.Error() == "no rows in result set" {
			return staff, ierr.ErrNotFound
		}
		return staff, err
	}

	return staff, nil
}

func (u *staffRepo) LookUp(ctx context.Context, id string) error {
	q := `SELECT 1 FROM staffs WHERE id = $1`

	v := 0
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

func (u *staffRepo) UpdateAccount(ctx context.Context, id, name, url string) error {
	q := `UPDATE staffs SET image_url = $1, name = $2 WHERE id = $3`
	_, err := u.conn.Exec(ctx, q,
		url, name, id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (u *staffRepo) GetNameBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT name FROM staffs WHERE id = $1`

	v := ""
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		return "", err
	}

	return v, nil
}

func (u *staffRepo) GetEmailBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT email FROM staffs WHERE id = $1`

	v := ""
	err := u.conn.QueryRow(ctx,
		q, id).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		return "", err
	}

	return v, nil
}

// func (u *staffRepo) GetNameByID(ctx context.Context, id string) (string, error) {
// 	name := ""
// 	err := u.conn.QueryRow(ctx,
// 		`SELECT name FROM staffs
// 		WHERE id = $1`,
// 		id).Scan(&name)
// 	if err != nil {
// 		if err.Error() == "no rows in result set" {
// 			return "", ierr.ErrNotFound
// 		}
// 		if pgErr, ok := err.(*pgconn.PgError); ok {
// 			if pgErr.Code == "22P02" {
// 				return "", ierr.ErrNotFound
// 			}
// 		}
// 		return "", err
// 	}

// 	return name, nil
// }
