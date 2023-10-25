package entity

import "database/sql"

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Sex         string `json:"sex"`
}

type UserDAO struct {
	Id          int
	Name        string
	Surname     string
	Patronymic  sql.NullString
	Age         int
	Nationality string
	Sex         string
}
