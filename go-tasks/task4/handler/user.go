package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mattn/go-sqlite3"
	"github.com/q1ngy/go-tasks/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) RegisterRoute(server *gin.Engine) {
	ug := server.Group("user")
	{
		ug.POST("/register", h.register)
		ug.POST("/login", h.Login)
	}
}

func (h *UserHandler) register(ctx *gin.Context) {
	var req = struct {
		Username        string `json:"username" binding:"required,min=3,max=30"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}{}
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, Resp{
			Success: false,
			Message: "Passwords do not match",
		})
		return
	}
	encrypted, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
	}
	result := h.db.Model(&domain.User{}).Create(&domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(encrypted),
	})
	err = result.Error
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) {
			if errors.Is(sqlErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				if strings.Contains(sqlErr.Error(), "users.username") {
					ctx.JSON(http.StatusBadRequest, Resp{
						Success: false,
						Message: "用户名已存在",
					})
					return
				}
				if strings.Contains(sqlErr.Error(), "users.email") {
					ctx.JSON(http.StatusBadRequest, Resp{
						Success: false,
						Message: "邮箱已存在",
					})
					return
				}
			}
		}
	}
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
	})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req = struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}{}
	if err := ctx.Bind(&req); err != nil {
		return
	}
	var user domain.User
	result := h.db.Model(&domain.User{}).Where("username = ?", req.Username).Find(&user)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, Resp{
			Success: false,
			Message: "用户名或密码不正确",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		ctx.JSON(http.StatusOK, Resp{
			Success: false,
			Message: "用户名或密码不正确",
		})
		return
	}
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("123456"))
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
		Data:    tokenStr,
	})
}
