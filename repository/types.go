package repository

type InsertUser struct {
	PhoneNumber string
	Name        string
	Password    string
}

type GetUser struct {
	PhoneNumber string
	Password    string
}

type User struct {
	ID          int
	PhoneNumber string
	Name        string
	Password    string
}

type InsertUserLoginHistory struct {
	ID          int
	PhoneNumber string
	Name        string
	Password    string
}

type UpdateUser struct {
	ID          int
	PhoneNumber string
	Name        string
}
