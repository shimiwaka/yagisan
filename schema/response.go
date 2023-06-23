package schema

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `jsoin:"token"`
}

type ShowBoxResponse struct {
	Success   bool       `json:"success"`
	Questions []Question `json:"questions"`
}

type SendQuestionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `jsoin:"token"`
}

type SendAnswerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}
