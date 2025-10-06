package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/cache"
	"github.com/q1ngy/Learn-Go/webook/internal/repository/dao"
)

var (
	EmailDuplicateErr = dao.ErrEmailDuplicate
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, user domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CachedUserRepository) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func (repo *CachedUserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: time.UnixMilli(user.Birthday),
		AboutMe:  user.AboutMe,
		CTime:    time.UnixMilli(user.CTime),
	}
}

func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
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

func (repo *CachedUserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
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
		log.Println("[RedisUserCache]", err)
	}
	return du, nil
}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
