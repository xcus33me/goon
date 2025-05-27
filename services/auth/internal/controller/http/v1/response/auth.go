package response

type Register struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type Login struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}
