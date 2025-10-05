package routes

import (
	controller "github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/controllers"
	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnProtectedRoutes(router *gin.Engine) {

	router.GET("/movies", controller.GetMovies(nil))
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
}