package schema

type RegisterResponse struct {
	Success string `json:"sucess"`
	Message string `json:"message"`
	Token   string `jsoin:"token"`
}
