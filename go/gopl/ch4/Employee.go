package ch4

import "time"

type Employee struct {
	ID        int
	Name      string
	Address   string
	DoB       time.Time
	Position  string
	Salary    int
	ManagerID int
}

var dilbert Employee

func main() {
	var employeeOfTheMonth *Employee = &dilbert
	employeeOfTheMonth.Position += " (proactive team player)" // 自动解引用
	(*employeeOfTheMonth).Position += " (proactive team player)"

}
