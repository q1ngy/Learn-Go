package main

import "github.com/jmoiron/sqlx"
import _ "github.com/mattn/go-sqlite3"

var booksSchema = `
CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY,
    title TEXT,
    author TEXT,
    price INTEGER
);`

type Books struct {
	Id     int64
	Title  string
	Author string
	Price  int64
}

func main() {
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	db.MustExec(booksSchema)

	db.Exec("INSERT INTO books (id, title, author, price) VALUES ($1, $2, $3, $4)", 1, "On Java 8", "Slim", 260)
	db.Exec("INSERT INTO books (id, title, author, price) VALUES ($1, $2, $3, $4)", 2, "Go", "Slim", 30)

	var books []Books
	db.Select(&books, "SELECT * FROM books WHERE price > $1", 50)

	db.Exec("DELETE FROM books")
}
