package domain

import "time"

type User struct {
	Id       int64
	Nickname string
	Birthday time.Time
	Email    string
	Phone    string
	Password string
	AboutMe  string
	Ctime    time.Time
}
