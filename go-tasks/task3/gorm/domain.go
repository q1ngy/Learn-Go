package main

import "gorm.io/gorm"

type User struct {
	Id        int64 `gorm:"primaryKey;autoIncrement"`
	Name      string
	PostCount int32

	Posts []Post `gorm:"foreignKey:UserId"`
}

type Post struct {
	Id            int64 `gorm:"primaryKey;autoIncrement"`
	Title         string
	CommentCount  int32
	CommentStatus bool

	UserId int64
	User   User `gorm:"foreignKey:UserId"`

	Comments []Comment `gorm:"foreignKey:PostId"`
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(&User{}).
		Where("id = ?", p.UserId).
		Update("post_count", gorm.Expr("post_count + ?", 1))
	return
}

type Comment struct {
	Id      int64 `gorm:"primaryKey;autoIncrement"`
	Content string

	PostId int64
	Post   Post `gorm:"foreignKey:PostId"`
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	p := Post{}
	tx.First(&p, c.PostId)
	if p.CommentCount > 1 {
		tx.Model(&p).Update("comment_count", gorm.Expr("comment_count - ?", 1))
	} else if p.CommentCount == 1 {
		tx.Model(&p).Select("comment_count", "comment_status").Updates(Post{CommentCount: 0, CommentStatus: false})
	}
	return
}
