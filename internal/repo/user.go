package repo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"

	"eniqilo-store-be/internal/entity"
	"eniqilo-store-be/internal/ierr"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	conn *pgxpool.Pool
}

func newUserRepo(conn *pgxpool.Pool) *userRepo {
	return &userRepo{conn}
}

func (u *userRepo) Insert(ctx context.Context, user entity.User) (string, error) {
	credVal := user.PhoneNumber
	q := `INSERT INTO staff (name, phone_number, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, now(), now()) RETURNING id`

	var userID string
	err := u.conn.QueryRow(ctx, q,
		user.Name, credVal, user.Password).Scan(&userID)

	if err != nil {
		ierr.LogErrorWithLocation(err)
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", ierr.ErrDuplicate
			}
		}
		return "", err
	}

	return userID, nil
}

// func (u *userRepo) GetByEmail(ctx context.Context, cred string) (entity.User, error) {
// 	user := entity.User{}
// 	q := `SELECT id, name, email, password FROM users
// 	WHERE email = $1`

// 	var email sql.NullString

// 	err := u.conn.QueryRow(ctx,
// 		q, cred).Scan(&user.ID, &user.Name, &email, &user.Password)

// 	user.Email = email.String

// 	if err != nil {
// 		if err.Error() == "no rows in result set" {
// 			return user, ierr.ErrNotFound
// 		}
// 		return user, err
// 	}

// 	return user, nil
// }

func (u *userRepo) GetByPhoneNumber(ctx context.Context, cred string) (entity.User, error) {
	user := entity.User{}
	q := `SELECT id, name, phone_number, password_hash FROM staff
	WHERE phone_number = $1`

	var phoneNumber sql.NullString

	err := u.conn.QueryRow(ctx,
		q, cred).Scan(&user.ID, &user.Name, &phoneNumber, &user.Password)

	user.PhoneNumber = phoneNumber.String

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	user := entity.User{}
	q := `SELECT phone_number, name, password FROM users
	WHERE id = $1`

	var phoneNumber sql.NullString

	err := u.conn.QueryRow(ctx,
		q, id).Scan(&phoneNumber, &user.Name, &user.Password)

	user.PhoneNumber = phoneNumber.String

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) LookUp(ctx context.Context, id string) error {
	q := `SELECT 1 FROM users WHERE id = $1`

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

func (u *userRepo) UpdateAccount(ctx context.Context, id, name, url string) error {
	q := `UPDATE users SET image_url = $1, name = $2 WHERE id = $3`
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

func (u *userRepo) GetNameBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT name FROM users WHERE id = $1`

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

func (u *userRepo) GetEmailBySub(ctx context.Context, id string) (string, error) {
	q := `SELECT email FROM users WHERE id = $1`

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

// func (u *userRepo) GetNameByID(ctx context.Context, id string) (string, error) {
// 	name := ""
// 	err := u.conn.QueryRow(ctx,
// 		`SELECT name FROM users
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
