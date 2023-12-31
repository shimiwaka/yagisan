package schema

import (
	"time"
)

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `jsoin:"token"`
}

type UpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ShowBoxResponse struct {
	Success   bool       `json:"success"`
	Username  string     `json:"username"`
	Questions []Question `json:"questions"`
}

type ProfileResponse struct {
	Success     bool   `json:"success"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Description string `json:"description"`
	SecureMode  bool   `jsoin:"secure_mode"`
	Message     string `json:"message"`
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

type GetQuestionReponse struct {
	Success    bool      `json:"success"`
	Message    string    `jsoin:"message"`
	Email      string    `json:"email"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"useragent"`
	Body       string    `json:"body"`
	Token      string    `json:"token"`
	AnswerBody string    `json:"answer_body"`
	CreatedAt  time.Time `json:"created_at"`
	QuestionID uint      `json:"question_id"`
}
