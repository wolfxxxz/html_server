package responses

type CreateUserResponse struct {
	UserId string `json:"user_id"`
}

type Result struct {
	Answer string `json:"result"`
}

type GetTranslResponse struct {
	English string `json:"english"`
	Russian string `json:"russian"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type WordResp struct {
	ID            string `json:"id"`
	English       string `json:"english"`
	Russian       string `json:"russian"`
	PartsOfSpeech string `json:"part_of_speech"`
}
