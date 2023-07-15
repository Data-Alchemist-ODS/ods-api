package request

type UserCreateRequest struct {
	Name 		string `json:"name"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type UserLoginRequest struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}