package repository

import (
	"context"
	"fmt"
	models "github.com/Vitaly-Baidin/auth-api/internal/model/db"
	"github.com/Vitaly-Baidin/auth-api/pkg/postgres"
)

const (
	insertUserQuery     = `INSERT INTO users (login, email, phone, password) VALUES ($1, $2, $3, $4) RETURNING id`
	getUserByID         = `SELECT id, login, email, phone, password FROM users WHERE id=$1`
	getUserByEmail      = `SELECT id, login, email, phone, password FROM users WHERE email=$1`
	ifExistsUserByEmail = `SELECT exists(SELECT 1 FROM users WHERE email=$1)`
	ifExistsUserByLogin = `SELECT exists(SELECT 1 FROM users WHERE login=$1)`
)

type UserRepoPG struct {
	*postgres.Postgres
}

func NewUserRepoPG(pg *postgres.Postgres) *UserRepoPG {
	return &UserRepoPG{Postgres: pg}
}

func (r *UserRepoPG) Store(ctx context.Context, u *models.User) (uint, error) {
	var userID uint

	row := r.Pool.QueryRow(ctx, insertUserQuery, u.Login, u.Email, u.Phone, u.Password)
	err := row.Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed store user: %v", err)
	}

	return userID, nil
}

func (r *UserRepoPG) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User

	row := r.Pool.QueryRow(ctx, getUserByID, id)
	err := row.Scan(&u.ID, &u.Login, &u.Email, &u.Phone, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("failed get user by id: %v", err)
	}

	return &u, nil
}

func (r *UserRepoPG) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User

	row := r.Pool.QueryRow(ctx, getUserByEmail, email)
	err := row.Scan(&u.ID, &u.Login, &u.Email, &u.Phone, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("failed get user by email: %v", err)
	}

	return &u, nil
}

func (r *UserRepoPG) IfExistsByEmail(ctx context.Context, email string) error {
	var ifExists bool

	row := r.Pool.QueryRow(ctx, ifExistsUserByEmail, email)
	err := row.Scan(&ifExists)
	if err != nil {
		return fmt.Errorf("failed check exists user by email: %v", err)
	}

	if ifExists {
		return ErrEmailExists
	}

	return nil
}

func (r *UserRepoPG) IfExistsByLogin(ctx context.Context, login string) error {
	var ifExists bool

	row := r.Pool.QueryRow(ctx, ifExistsUserByLogin, login)
	err := row.Scan(&ifExists)
	if err != nil {
		return fmt.Errorf("failed check exists user by login: %v", err)
	}

	if ifExists {
		return ErrUserExists
	}

	return nil
}
