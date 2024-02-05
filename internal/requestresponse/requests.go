package requestresponse

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
