package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE IF NOT EXISTS employees (
    id INTEGER PRIMARY KEY,
    name TEXT,
    department TEXT,
    salary INTEGER
);`

type Employees struct {
	Id         int64
	Name       string
	Department string
	Salary     int64
}

func main() {
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	db.MustExec(schema)

	db.Exec("INSERT INTO employees (id, name, department, salary) VALUES ($1, $2, $3, $4)", 1, "Jason", "技术部", 100)
	db.Exec("INSERT INTO employees (id, name, department, salary) VALUES ($1, $2, $3, $4)", 2, "Max", "技术部", 200)

	var employees []Employees
	db.Select(&employees, "SELECT * FROM employees WHERE department = $1", "技术部")

	employee := Employees{}
	db.Get(&employee, "SELECT * FROM employees ORDER BY salary DESC limit 1")

	db.Exec("DELETE FROM employees")
}
