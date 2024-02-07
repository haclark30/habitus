package models

type Habit struct {
	Name      string
	Id        int
	UserId    int
	UpCount   int
	DownCount int
	HasUp     bool
	HasDown   bool
}

type Daily struct {
	Name   string
	Id     int
	UserId int
	Due    bool
	Done   bool
}

type User struct {
	UserName     string
	PasswordHash string
	Id           int
}
