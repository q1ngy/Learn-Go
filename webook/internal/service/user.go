package service

import (
	"context"
	"errors"

	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	EmailDuplicateErr        = repository.EmailDuplicateErr
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) SignUp(ctx context.Context, user domain.User) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encryptedPassword)
	return s.repo.Create(ctx, user)
}

func (s *UserServiceImpl) Login(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return s.repo.UpdateNonZeroFields(ctx, user)
}

func (s *UserServiceImpl) FindById(ctx context.Context, uid int64) (domain.User, error) {
	return s.repo.FindById(ctx, uid)
}

func (s *UserServiceImpl) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := s.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	err = s.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && !errors.Is(err, repository.EmailDuplicateErr) {
		return domain.User{}, err
	}
	// 可能存在主从延迟，理论上来说应该强制走主库
	return s.repo.FindByPhone(ctx, phone)
}
