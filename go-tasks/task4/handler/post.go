package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/q1ngy/go-tasks/domain"
	"gorm.io/gorm"
)

type PostHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *PostHandler {
	return &PostHandler{
		db: db,
	}
}

func (h *PostHandler) RegisterRoute(server *gin.Engine) {
	pg := server.Group("post")
	pg.POST("/create", h.create)
	pg.GET("/:id", h.getById)
	pg.GET("/all", h.getAll)
	pg.PUT("/:id", h.updateById)
	pg.DELETE("/:id", h.deleteById)
}

func (h *PostHandler) create(ctx *gin.Context) {
	var req = struct {
		Title   string `gorm:"not null"`
		Content string `gorm:"not null"`
	}{}
	if err := ctx.Bind(&req); err != nil {
	}

	userId, _ := ctx.Get("userId")

	h.db.Create(&domain.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId.(uint),
	})
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
	})

}

func (h *PostHandler) getById(ctx *gin.Context) {
	id := ctx.Param("id")
	var post domain.Post
	h.db.First(&post, id)
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
		Data:    post,
	})
}

func (h *PostHandler) getAll(ctx *gin.Context) {
	var posts []domain.Post
	h.db.Find(&posts)
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
		Data:    posts,
	})
}

func (h *PostHandler) updateById(ctx *gin.Context) {
	var req = struct {
		Title   string `gorm:"not null"`
		Content string `gorm:"not null"`
	}{}
	if err := ctx.Bind(&req); err != nil {
	}
	id := ctx.Param("id")
	var post domain.Post
	h.db.First(&post, id)
	userId, _ := ctx.Get("userId")
	if post.UserID != userId.(uint) {
		ctx.JSON(http.StatusOK, Resp{
			Success: true,
			Message: "没有权限",
		})
		return
	}
	h.db.Model(&post).Updates(domain.Post{Title: req.Title, Content: req.Content})
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
	})
}

func (h *PostHandler) deleteById(ctx *gin.Context) {
	id := ctx.Param("id")
	var post domain.Post
	h.db.First(&post, id)
	userId, _ := ctx.Get("userId")
	if post.UserID != userId.(uint) {
		ctx.JSON(http.StatusOK, Resp{
			Success: true,
			Message: "没有权限",
		})
		return
	}
	h.db.Delete(&post)
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
	})
}
