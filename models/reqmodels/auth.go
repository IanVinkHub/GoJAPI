package reqmodels

// LoginPost the post data for login
type LoginPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangePasswordPost the post data for changing password
type ChangePasswordPost struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

// TokenPost the post data for posting token
type TokenPost struct {
	Username  string `json:"username"`
	AuthToken string `json:"auth_token"`
}
