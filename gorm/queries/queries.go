package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User represents a user in the social media system
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Age       int    `gorm:"not null"`
	Country   string `gorm:"not null"`
	CreatedAt time.Time
	Posts     []Post `gorm:"foreignKey:UserID"`
	Likes     []Like `gorm:"foreignKey:UserID"`
}

// Post represents a social media post
type Post struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Content     string `gorm:"type:text"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
	Category    string `gorm:"not null"`
	ViewCount   int    `gorm:"default:0"`
	IsPublished bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Likes       []Like `gorm:"foreignKey:PostID"`
}

// Like represents a user's like on a post
type Like struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	PostID    uint `gorm:"not null"`
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
	CreatedAt time.Time
}

// ConnectDB establishes a connection to the SQLite database with auto-migration
func ConnectDB() (*gorm.DB, error) {
	// Implement database connection with auto-migration
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Post{}, &Like{})
	return db, nil
}

// GetTopUsersByPostCount retrieves users with the most posts
func GetTopUsersByPostCount(db *gorm.DB, limit int) ([]User, error) {
	// Implement top users by post count aggregation
	// select u.id, count(p.id)
	// from users u
	// left join posts p on p.user_id = u.id
	// group by u.id
	// order by count(p.id) desc
	// limit 10
	var users []User
	db.Model(&User{}).
		Select("users.*").
		InnerJoins("left join posts on posts.user_id = users.id").
		Group("users.id").
		Order("count(posts.id) desc").
		Limit(limit).
		Find(&users)
	return users, nil
}

// GetPostsByCategoryWithUserInfo retrieves posts by category with pagination and user info
func GetPostsByCategoryWithUserInfo(db *gorm.DB, category string, page, pageSize int) ([]Post, int64, error) {
	// Implement paginated posts retrieval with user info
	if page < 1 {
		return nil, 0, nil
	}
	offset := (page - 1) * pageSize
	var posts []Post
	var total int64
	db.Model(&Post{}).Where("category = ?", category).Count(&total)
	db.Preload("User").Model(&Post{}).Where("category = ?", category).Limit(pageSize).Offset(offset).Find(&posts)
	return posts, total, nil
}

// GetUserEngagementStats calculates engagement statistics for a user
func GetUserEngagementStats(db *gorm.DB, userID uint) (map[string]interface{}, error) {
	// Implement user engagement statistics
	var user User
	db.Preload("Posts.Likes").Preload("Likes").First(&user, userID)
	stat := make(map[string]interface{}, 4)
	stat["total_posts"] = len(user.Posts)
	stat["total_likes_given"] = len(user.Likes)
	totalLikes := 0
	for _, post := range user.Posts {
		totalLikes += len(post.Likes)
	}
	stat["total_likes_received"] = totalLikes
	totalViews := 0
	for _, post := range user.Posts {
		totalViews += post.ViewCount
	}
	averageViews := totalViews / len(user.Posts)
	stat["average_views_per_post"] = averageViews
	return stat, nil
}

// GetPopularPostsByLikes retrieves popular posts by likes in a time period
func GetPopularPostsByLikes(db *gorm.DB, days int, limit int) ([]Post, error) {
	// Implement popular posts by likes
	end := time.Now()
	start := end.AddDate(0, 0, -days)
	var posts []Post
	db.Model(&Post{}).
		Preload("Likes").
		Select("posts.*").
		InnerJoins("left join likes on likes.post_id = posts.id").
		Where("posts.created_at >= ? and posts.created_at <= ?", start, end).
		Group("posts.id").
		Order("count(likes.id) desc").
		Limit(limit).
		Find(&posts)
	return posts, nil
}

// GetCountryUserStats retrieves user statistics grouped by country
func GetCountryUserStats(db *gorm.DB) ([]map[string]interface{}, error) {
	// Implement country-based user statistics
	// select country, count(*)
	// from users
	// group by country
	var stat []struct {
		Country   string
		UserCount int
	}
	db.Model(&User{}).Select("country, count(*) as user_count").Group("country").Scan(&stat)
	var result []map[string]interface{}
	for _, s := range stat {
		result = append(result, map[string]interface{}{
			"country":   s.Country,
			"userCount": s.UserCount,
		})
	}
	return result, nil
}

// SearchPostsByContent searches posts by content using full-text search
func SearchPostsByContent(db *gorm.DB, query string, limit int) ([]Post, error) {
	// TODO: Implement full-text search
	return nil, nil
}

// GetUserRecommendations retrieves user recommendations based on similar interests
func GetUserRecommendations(db *gorm.DB, userID uint, limit int) ([]User, error) {
	// TODO: Implement user recommendations algorithm
	return nil, nil
}
