package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	controllers "github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/controllers"
)

func main() {
	fmt.Println("Hello, World!")

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/movies", controllers.getMovies())
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to run server:", err)
	}
}

// This is a simple Go server application that prints "Hello, World!" to the console.