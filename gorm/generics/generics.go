package main

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// User represents a user in the system with company association
type User struct {
	ID        uint     `gorm:"primaryKey"`
	Name      string   `gorm:"not null"`
	Email     string   `gorm:"unique;not null"`
	Age       int      `gorm:"check:age > 0"`
	CompanyID *uint    `gorm:"index"`
	Company   *Company `gorm:"foreignKey:CompanyID"`
	Posts     []Post   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Company represents a company that users can belong to
type Company struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null;unique"`
	Industry    string `gorm:"not null"`
	FoundedYear int    `gorm:"not null"`
	Users       []User `gorm:"foreignKey:CompanyID"`
	CreatedAt   time.Time
}

// Post represents a blog post by a user
type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID"`
	ViewCount int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConnectDB establishes a connection to the SQLite database and auto-migrates models
func ConnectDB() (*gorm.DB, error) {
	// Connect to SQLite database and auto-migrate all models
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
	db.AutoMigrate(&User{}, &Company{}, &Post{})
	return db, nil
}

// CreateUser creates a new user using GORM's generics API
func CreateUser(ctx context.Context, db *gorm.DB, user *User) error {
	// Use gorm.G[User](db).Create() with context
	err := gorm.G[User](db).Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by ID using generics API
func GetUserByID(ctx context.Context, db *gorm.DB, id uint) (*User, error) {
	// Use gorm.G[User](db).Where().First() with context
	user, err := gorm.G[User](db).Where("id = ?", id).First(ctx)
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

// UpdateUserAge updates a user's age using generics API
func UpdateUserAge(ctx context.Context, db *gorm.DB, userID uint, age int) error {
	// Use gorm.G[User](db).Where().Update() with context
	_, err := gorm.G[User](db).Where("id = ?", userID).Update(ctx, "age", age)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID using generics API
func DeleteUser(ctx context.Context, db *gorm.DB, userID uint) error {
	// Use gorm.G[User](db).Where().Delete() with context
	_, err := gorm.G[User](db).Where("id = ?", userID).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// CreateUsersInBatches creates multiple users in batches using generics API
func CreateUsersInBatches(ctx context.Context, db *gorm.DB, users []User, batchSize int) error {
	// Use gorm.G[User](db).CreateInBatches() with context
	err := gorm.G[User](db).CreateInBatches(ctx, &users, batchSize)
	if err != nil {
		return err
	}
	return nil
}

// FindUsersByAgeRange finds users within an age range using generics API
func FindUsersByAgeRange(ctx context.Context, db *gorm.DB, minAge, maxAge int) ([]User, error) {
	// Use gorm.G[User](db).Where() with range conditions and Find()
	users, err := gorm.G[User](db).Where("age >= ? and age <= ?", minAge, maxAge).Find(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpsertUser creates or updates a user handling conflicts using OnConflict
func UpsertUser(ctx context.Context, db *gorm.DB, user *User) error {
	// Use gorm.G[User](db, clause.OnConflict{...}).Create() with conflict handling
	err := gorm.G[User](db, clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "age", "company_id"}),
	}).Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// CreateUserWithResult creates a user and returns result metadata
func CreateUserWithResult(ctx context.Context, db *gorm.DB, user *User) (int64, error) {
	// Use gorm.WithResult() to capture metadata and return rows affected
	result := gorm.WithResult()
	err := gorm.G[User](db, result).Create(ctx, user)
	return result.RowsAffected, err
}

// GetUsersWithCompany retrieves users with their company information using enhanced joins
// https://gorm.io/zh_CN/docs/the_generics_way.html
func GetUsersWithCompany(ctx context.Context, db *gorm.DB) ([]User, error) {
	// Use gorm.G[User](db).Joins() with enhanced join syntax
	users, err := gorm.G[User](db).Joins(clause.Has("Company"), nil).Find(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersWithPosts retrieves users with their posts using enhanced preloading
func GetUsersWithPosts(ctx context.Context, db *gorm.DB, limit int) ([]User, error) {
	// Use gorm.G[User](db).Preload() with LimitPerRecord
	users, err := gorm.G[User](db).Preload("Posts", func(db gorm.PreloadBuilder) error {
		db.Limit(limit)
		return nil
	}).Find(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserWithPostsAndCompany retrieves a user with both posts and company preloaded
func GetUserWithPostsAndCompany(ctx context.Context, db *gorm.DB, userID uint) (*User, error) {
	// Use multiple Preload() calls with generics API
	user, err := gorm.G[User](db).Preload("Posts", nil).Preload("Company", nil).Where("id = ?", userID).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SearchUsersInCompany finds users working in a specific company using join with filters
func SearchUsersInCompany(ctx context.Context, db *gorm.DB, companyName string) ([]User, error) {
	// Use enhanced joins with custom filter functions
	users, err := gorm.G[User](db).Joins(clause.Has("Company"), nil).Where("Company.name = ?", companyName).Find(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetTopActiveUsers retrieves users with the most posts using complex joins and grouping
func GetTopActiveUsers(ctx context.Context, db *gorm.DB, limit int) ([]User, error) {
	// Use joins, grouping, and ordering to find most active users
	// select u.*
	// from users u
	// join posts p on u.id = p.user_id
	// group by u.id
	// order by count(p.id) desc
	// limit 10
	var users []User
	db.Model(&User{}).
		Joins("join posts on users.id = posts.user_id").
		Group("users.id").
		Order("count(posts.id) desc").
		Limit(limit).
		Find(&users)
	return users, nil
}
