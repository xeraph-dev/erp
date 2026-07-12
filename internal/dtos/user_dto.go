package dtos

type UserLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRegister struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
