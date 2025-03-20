package dto

type Response struct {
	Msg string `json:"message"`
}

type User struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
