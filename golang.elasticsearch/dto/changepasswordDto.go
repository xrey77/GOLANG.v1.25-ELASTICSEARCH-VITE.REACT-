package dto

type ChangePassword struct {
	Password string `json:"password" binding:"required"`
}
