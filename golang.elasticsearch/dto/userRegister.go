package dto

type UserRegister struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
