package main

import (
	"github.com/gin-gonic/gin"
	"github.com/q1ngy/go-tasks/domain"
	"github.com/q1ngy/go-tasks/handler"
	"github.com/q1ngy/go-tasks/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// gorm
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.User{}, &domain.Post{}, &domain.Comment{})

	// gin
	server := gin.New()
	server.Use(middleware.ErrorHandlerMiddleware())
	server.Use(middleware.AuthMiddleware())

	userHandler := handler.NewUserHandler(db)
	userHandler.RegisterRoute(server)

	postHandler := handler.NewPostHandler(db)
	postHandler.RegisterRoute(server)

	commentHandler := handler.NewCommentHandler(db)
	commentHandler.RegisterRoute(server)

	server.Run(":8080")
}
