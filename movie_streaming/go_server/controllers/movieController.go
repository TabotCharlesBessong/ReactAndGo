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
	"github.com/go-playground/validator/v10"

)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

var validate = validator.New()

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will go here
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		
		defer cancel()

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

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will go here
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
 
		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie"})
			}
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

func CreateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will go here
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie

		if err := c.BindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		_, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Validation failed","details": err.Error()})

			return
		}

		// insert data

		result, err := movieCollection.InsertOne(ctx,movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
