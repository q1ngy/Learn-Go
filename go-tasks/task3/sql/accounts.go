package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Accounts struct {
	Id      int64 `gorm:"primaryKey"`
	Balance int64
}

type Transactions struct {
	Id            int64 `gorm:"primaryKey;autoIncrement"`
	FromAccountId int64
	ToAccountId   int64
	Amount        int64
}

var db *gorm.DB

func main() {
	defer clearTable()

	// 1 账户 100 1->3 100 提交
	// 2 账户 50 2->3 50 回滚
	// 3 账户 0

	var a Accounts
	_ = db.Raw("SELECT * FROM accounts WHERE id = 1").Scan(&a)
	var b Accounts
	_ = db.Raw("SELECT * FROM accounts WHERE id = 2").Scan(&b)
	var c Accounts
	_ = db.Raw("SELECT * FROM accounts WHERE id = 3").Scan(&c)

	transfer(&a, &c, 100)
	transfer(&b, &c, 100)
}

func transfer(a *Accounts, b *Accounts, amount int64) {
	db.Exec("BEGIN")
	//db.Begin() 这样事务才会在一个连接中
	if a.Balance >= amount {
		a.Balance = a.Balance - amount
		b.Balance = b.Balance + amount
		db.Exec("UPDATE accounts SET balance = ? WHERE id = ?", a.Balance, a.Id)
		db.Exec("UPDATE accounts SET balance = ? WHERE id = ?", b.Balance, b.Id)
		db.Exec("INSERT INTO transactions(from_account_id, to_account_id, amount) VALUES (?, ?, ?)", a.Id, b.Id, amount)
		db.Exec("COMMIT")
	} else {
		db.Exec("ROLLBACK")
	}
}

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Accounts{})
	db.AutoMigrate(&Transactions{})
	db.Exec("INSERT INTO accounts VALUES (1, 100)")
	db.Exec("INSERT INTO accounts VALUES (2, 50)")
	db.Exec("INSERT INTO accounts VALUES (3, 0)")

}

func clearTable() {
	db.Exec("DELETE FROM accounts")
}
