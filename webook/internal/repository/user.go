package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/cache"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/dao"
)

var (
	EmailDuplicateErr = dao.ErrEmailDuplicate
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func (repo *UserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: time.UnixMilli(user.Birthday),
		AboutMe:  user.AboutMe,
	}
}

func (repo *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func (repo *UserRepository) UpdateNonZeroFields(ctx *gin.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *UserRepository) FindById(ctx *gin.Context, uid int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return u, nil
	}

	user, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du := repo.toDomain(user)

	err = repo.cache.Set(ctx, du)
	if err != nil {
		log.Println("[UserCache]", err)
	}
	return du, nil
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
