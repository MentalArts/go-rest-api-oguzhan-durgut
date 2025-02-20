package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg string `json:"message"`
}

func main() {
	router := gin.Default()
	router.GET("/ping", handlePing)
	router.GET("/hello", handleHello)
	router.GET("/helloWithPayload", handleHelloWithPayload)
	router.Run(":8000")
}

func handlePing(c *gin.Context) {
	res := Response{Msg: "pong"}
	c.JSON(http.StatusOK, res)
}

func handleHello(c *gin.Context) {
	name := c.Query("name")

	var msg string
	if name != "" {
		msg = fmt.Sprintf("Welcome, %s", name)
	} else {
		msg = "Welcome, user"
	}

	c.String(http.StatusOK, msg)
}

type DTO struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func handleHelloWithPayload(c *gin.Context) {
	// Binding (Get Payload from request)
	var dto DTO
	err := c.BindJSON(&dto)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

	// Validaton (Validate payload)
	if dto.Name == "" {
		c.String(http.StatusBadRequest, "empty name is not accepted")
		return
	}

	msg := fmt.Sprintf("Welcome, %s %s", dto.Name, dto.Surname)
	c.String(http.StatusOK, msg)

}

// Vanilla implementation
// func main() {
// 	http.HandleFunc("GET /ping", handlePing)
// 	log.Println("Server listening...")
// 	log.Fatal(http.ListenAndServe(":8000", nil))
// }

// func handlePing(w http.ResponseWriter, r *http.Request) {
// 	res := Response{Msg: "pong"}
// 	json.NewEncoder(w).Encode(res)
// 	w.WriteHeader(http.StatusOK)
// 	log.Println("Request recieved")
// }
