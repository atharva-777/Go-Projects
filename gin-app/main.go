package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []user{
	{ID: "1", Name: "Alice", Email: "alice@example.com"},
	{ID: "2", Name: "Bob", Email: "bob@example.com"},
}

func getUsers(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{"users": users})
}

func postUsers(ctx *gin.Context) {
	var newUser user
	if err := ctx.BindJSON(&newUser); err != nil {
		return
	}

	users = append(users, newUser)
	ctx.IndentedJSON(http.StatusCreated, newUser)
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(ctx *gin.Context) {
		data := map[string]interface{}{"name": "gin", "version": "1.0"}
		ctx.JSON(http.StatusOK, data)
	})

	router.GET("/users", getUsers)

	router.POST("/users", postUsers)

	router.GET("/ping", func(ctx *gin.Context) {
		fmt.Printf("ClientIP: %s\n", ctx.ClientIP())
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.Run("localhost:8080") // listens on :8080 by default
}
