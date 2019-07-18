package resmodels

// LoginPost the response data for login
type LoginPost struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Auth     TokenPost `json:"auth"`
}

// TokenPost the response data for Token Info
type TokenPost struct {
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	ExpireAt  string `json:"expire_at"`
	DeleteAt  string `json:"delete_at"`
}
