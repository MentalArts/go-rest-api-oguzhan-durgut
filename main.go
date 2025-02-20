package main

import (
	"log"
	"mentalartsapi/handlers"
	"mentalartsapi/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dsn := "host=localhost user=postgres password=123abcd dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect database: %v", err)
	}

	db.AutoMigrate(&models.Author{})

	handlers.InitDB(db)

	router := gin.Default()

	router.GET("/ping", handlers.HandlePing)
	router.GET("/hello", handlers.HandleHello)
	router.GET("/helloWithPayload", handlers.HandleHelloWithPayload)

	router.POST("/author", handlers.CreateAuthor)
	router.GET("/author", handlers.GetAllAuthors)
	router.GET("/author/:id", handlers.GetAuthor)
	router.PUT("/author/:id", handlers.UpdateAuthor)
	router.DELETE("/author/:id", handlers.DeleteAuthor)

	router.Run(":8000")
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
