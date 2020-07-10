package model

type LoginRequest struct {
	HotelCode string `json:"hotel_code"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserInfo struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ResponseLogin struct {
	Token string   `json:"token"`
	Info  UserInfo `json:"info"`
}
