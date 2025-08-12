package entity

type CreateUser struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type UserInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Debt        float32 `json:"debt"`
	CreatedAt   string  `json:"created_at"`
}

type UpdateUser struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateUserBody struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type UserList struct {
	Users []UserInfo `json:"users"`
	Count int        `json:"count"`
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRes struct {
	UserInfo UserInfo `json:"user_info"`
	Token    string   `json:"token"`
}
