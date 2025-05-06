package response

type Register struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type Login struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Token string `json:"token"`
}
