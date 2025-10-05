
package routes

import (
	controller "github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/controllers"
	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/middleware"
	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.CreateMovie())
	// router.GET("/recommendedmovies", controller.GetRecommendedMovies())
	// router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate())
}