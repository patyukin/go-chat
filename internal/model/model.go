package model

type ErrorResponse struct {
	Error *string `json:"error,omitempty"`
}

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UUID string `json:"user_id"`
}
