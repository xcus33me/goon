package request

type Register struct {
	Login    string `json:"login" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=6"`
}

type Login struct {
	Login    string `json:"login" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=6"`
}
