package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Students struct {
	Id    int64 `gorm:"primaryKey;autoIncrement"`
	Name  string
	Age   int32
	Grade string
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Students{})

	// insert
	_ = db.Exec("INSERT INTO students(name, age, grade) VALUES ('张三', 20, '三年级')")

	// find
	var students []Students
	result := db.Raw("SELECT * FROM students WHERE age > 18").Scan(&students)
	if result.Error != nil {
	}

	// update
	_ = db.Exec("UPDATE students SET grade = '四年级' WHERE name = '张三'")

	// delete
	_ = db.Exec("DELETE FROM students WHERE age < 14")
}
