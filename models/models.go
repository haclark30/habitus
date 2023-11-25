package models

type Habit struct {
	Id        int
	Name      string
	UpCount   int
	DownCount int
	HasDown   bool
}

type Daily struct {
	Id   int
	Name string
	Due  bool
	Done bool
}
