package main

import (
	"IUMS/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	r := gin.Default()
	iumsRepo := controllers.New()

	/*
	* line 25 will create a table and insert data into it,
	* please execute it only once
	 */
	//iumsRepo.InsertDataIntoDb()

	r.GET("/analytics", iumsRepo.GetUsageDetails)
	r.GET("/user/search", iumsRepo.GetUserDetails)

	r.Run("localhost:8000")
	fmt.Println("Server is running")
}
