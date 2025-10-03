package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/q1ngy/go-tasks/domain"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{
		db: db,
	}
}

func (h *CommentHandler) RegisterRoute(server *gin.Engine) {
	cg := server.Group("comment")
	cg.POST("/create", h.create)
	cg.GET("/by-post/:post-id", h.getByPostId)
}

func (h *CommentHandler) create(ctx *gin.Context) {
	var req = struct {
		PostId  uint   `gorm:"not null"`
		Content string `gorm:"not null"`
	}{}
	if err := ctx.Bind(&req); err != nil {
	}

	userId, _ := ctx.Get("userId")

	h.db.Create(&domain.Comment{
		Content: req.Content,
		UserID:  userId.(uint),
		PostID:  req.PostId,
	})
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
	})
}

func (h *CommentHandler) getByPostId(ctx *gin.Context) {
	postId := ctx.Param("post-id")
	var comments []domain.Comment
	h.db.Model(&domain.Comment{}).Where("post_id = ?", postId).Find(&comments)
	ctx.JSON(http.StatusOK, Resp{
		Success: true,
		Data:    comments,
	})
}
