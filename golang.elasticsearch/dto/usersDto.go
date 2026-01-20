package dto

type Users struct {
	Id          string  `json:"id"`
	Firstname   string  `json:"firstname"`
	Lastname    string  `json:"lastname"`
	Email       string  `json:"email"`
	Mobile      string  `json:"mobile"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Roles       string  `json:"roles"`
	Isactivated bool    `json:"isactivated"`
	Isblocked   bool    `json:"isblocked"`
	Userpicture string  `json:"userpicture"`
	Mailtoken   float64 `json:"mailtoken"`
	Qrcodeurl   *string `json:"qrcodeurl"`
	Secret      *string `json:"secret"`
}
