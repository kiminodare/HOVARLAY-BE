package dtoAuth

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
