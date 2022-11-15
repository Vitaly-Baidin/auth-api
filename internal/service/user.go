package service

import (
	"context"
	"fmt"
	db "github.com/Vitaly-Baidin/auth-api/internal/model/db"
	"github.com/Vitaly-Baidin/auth-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *db.User) (*db.User, error)
	FindUserById(ctx context.Context, userID uint) (*db.User, error)
	FindUserByEmail(ctx context.Context, email string) (*db.User, error)
	CheckUserMail(ctx context.Context, email string) error
	CheckUserLogin(ctx context.Context, login string) error
}

type UserServ struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserServ {
	return &UserServ{repo: r}
}

func (s *UserServ) CreateUser(ctx context.Context, u *db.User) (*db.User, error) {
	password, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed generate hash pass: %v", err)
	}

	u.Password = password

	userID, err := s.repo.Store(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed store to db: %v", err)
	}

	u.ID = userID

	return u, nil
}

func (s *UserServ) FindUserById(ctx context.Context, userID uint) (*db.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed find user by id: %v", err)
	}

	return user, nil
}

func (s *UserServ) FindUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed find user by email: %v", err)
	}

	return user, nil
}

func (s *UserServ) CheckUserMail(ctx context.Context, email string) error {
	err := s.repo.IfExistsByEmail(ctx, email)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServ) CheckUserLogin(ctx context.Context, login string) error {
	err := s.repo.IfExistsByLogin(ctx, login)
	if err != nil {
		return err
	}

	return nil
}
