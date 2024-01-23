package models

type Habit struct {
	Id        int
	UserId    int
	Name      string
	UpCount   int
	DownCount int
	HasDown   bool
}

type Daily struct {
	Id     int
	UserId int
	Name   string
	Due    bool
	Done   bool
}

type User struct {
	Id           int
	UserName     string
	PasswordHash string
}
