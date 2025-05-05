package request

type Register struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Login struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
