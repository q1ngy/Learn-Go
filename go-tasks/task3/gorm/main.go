package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	defer clearTable()

	// 关联查询
	u := User{}
	db.Preload("Posts.Comments").First(&u, 1)
	fmt.Println(u)

	p := Post{}
	db.Model(&Post{}).
		Joins("LEFT JOIN comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("COUNT(comments.id) DESC").
		Limit(1).
		Scan(&p)
	fmt.Println(p)

	// Hook
	p3 := Post{
		Id:            3,
		Title:         "Post3",
		CommentCount:  0,
		CommentStatus: false,
		UserId:        1,
	}
	fmt.Printf("user posts count before create a post: %d\n", u.PostCount)
	db.Create(&p3)
	db.First(&u, p3.UserId)
	fmt.Printf("user posts count after create a post: %d\n", u.PostCount)

	// 删除有两条评论的 Post
	p1 := Post{}
	db.First(&p1, 1)
	fmt.Printf("post comment count before create a post: %d, commentStatus: %t\n", p1.CommentCount, p1.CommentStatus)
	c1 := Comment{}
	db.First(&c1, 1)
	db.Delete(&c1)
	db.First(&p1, 1)
	fmt.Printf("post comment count after create a post: %d, commentStatus: %t\n", p1.CommentCount, p1.CommentStatus)

	// 删除只有一条评论的 Post
	p2 := Post{}
	db.First(&p2, 2)
	fmt.Printf("post comment count before create a post: %d, commentStatus: %t\n", p2.CommentCount, p2.CommentStatus)
	c3 := Comment{}
	db.First(&c3, 3)
	db.Delete(&c3)
	db.First(&p2, 2)
	fmt.Printf("post comment count after create a post: %d, commentStatus: %t\n", p2.CommentCount, p2.CommentStatus)
}

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	u := User{Id: 1, Name: "slim", PostCount: 0}
	p1 := Post{Id: 1, Title: "Post1", CommentCount: 2, CommentStatus: true, UserId: 1}
	p2 := Post{Id: 2, Title: "Post2", CommentCount: 1, CommentStatus: true, UserId: 1}
	c1 := Comment{Id: 1, Content: "Post1 Comment1", PostId: 1}
	c2 := Comment{Id: 2, Content: "Post1 Comment2", PostId: 1}
	c3 := Comment{Id: 3, Content: "Post2 Comment3", PostId: 2}

	db.Create(&u)
	db.Create(&p1)
	db.Create(&p2)
	db.Create(&c1)
	db.Create(&c2)
	db.Create(&c3)
}

func clearTable() {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM comments")
}
