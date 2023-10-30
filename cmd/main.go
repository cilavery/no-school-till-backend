package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	internal "github.com/cilavery/no-school-till-backend/cmd/main.go/internal/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// TODO: load different vars based off of environment i.e. dev, prod
func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load .env file")
		os.Exit(1)
	}
}

func main() {
	loadEnv()

	controller := internal.NewController()
	r := mux.NewRouter()
	r.HandleFunc("/courses", controller.GetCourseInfo())
	fmt.Println("Server is listening on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
