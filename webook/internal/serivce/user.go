package serivce

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	EmailDuplicateErr        = repository.EmailDuplicateErr
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) SignUp(ctx *gin.Context, user domain.User) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encryptedPassword)
	return s.repo.Create(ctx, user)
}

func (s *UserService) Login(ctx *gin.Context, email string, password string) (domain.User, error) {
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
