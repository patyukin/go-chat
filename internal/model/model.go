package model

type MainPageData struct {
	Users []User `json:"users"`
	Rooms []Room `json:"rooms"`
	User  User   `json:"user"`
}

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

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	UUID string `json:"id"`
}
