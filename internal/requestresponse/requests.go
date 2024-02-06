package requestresponse

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoomCreateRequest struct {
	Name           string `json:"name"`
	RecordMessages bool   `json:"record_messages"`
	AllowAnon      bool   `json:"allow_anon"`
}
