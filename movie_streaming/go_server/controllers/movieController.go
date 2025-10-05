package controllers

import (
	"net/http"
	"time"
	"context"

	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/database"
	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

)


func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will go here
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		
		defer cancel()
		var movieCollection *mongo.Collection = database.OpenCollection("movies")

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx,bson.M{})

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to fetch movies."})
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx,&movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to fetch movies."})
		}

		c.JSON(http.StatusOK,movies)
	}
}
