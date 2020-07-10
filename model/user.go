package model

type User struct {
	ID        int64  `json:"id" gorm:"column:id"`
	Email     string `json:"email" gorm:"column:email1"`
	Username  string `json:"username" gorm:"column:user_name"`
	Password  string `json:"-" gorm:"column:user_password"`
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
}
