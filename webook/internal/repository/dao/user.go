package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrEmailDuplicate = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CTime = now
	user.UTime = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	if sqlError, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if sqlError.Number == duplicateErr {
			return ErrEmailDuplicate
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx *gin.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *UserDao) UpdateById(ctx *gin.Context, entity User) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		Updates(map[string]any{
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

func (dao *UserDao) FindById(ctx *gin.Context, uid int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&user).Error
	return user, err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string         `gorm:"type=varchar(4096)"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	//  时区：UTC 0 毫秒
	CTime int64
	UTime int64
}
