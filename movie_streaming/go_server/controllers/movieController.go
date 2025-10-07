package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/database"
	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/models"
	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")

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

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will go here

		role,err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}

		if role != "ADMIN"{
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can update reviews"})
			return
		}

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		sentiment, rankVal, err := GetReviewEanking(req.AdminReview)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze review"})
			return
		}

		filter := bson.M{"imdb_id": movieId}
		update := bson.M{
			"$set": bson.M{
				"admin_review":  req.AdminReview,
				"ranking": bson.M{
					"ranking_name":  sentiment,
					"ranking_value": rankVal,
				},
			},
		}

		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)
		
		defer cancel()


		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie review"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)

	}
}

func GetReviewEanking(admin_review string) (string, int, error) {
	rankings, err := GetRankings()

	if err != nil{
		return "", 0 , err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings{
		if ranking.RankingValue != 999{
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}

	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	OpenAiApiKey := os.Getenv("OPEN_AI_API_KEY")

	if OpenAiApiKey == "" {
		return "", 0, errors.New("OPEN_AI_API_KEY is not set in the environment variables")
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	if base_prompt_template == "" {
		return "", 0, errors.New("BASE_PROMPT_TEMPLATE is not set in the environment variables")
	}

	base_prompt := strings.ReplaceAll(base_prompt_template, "{rankings}", sentimentDelimited)

	response, err := llm.Call(context.Background(), base_prompt + " " + admin_review)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0

	for _, ranking := range rankings{
		if strings.EqualFold(ranking.RankingName, response){
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}

func GetRankings() ([]models.Ranking, error) {

	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx,bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx,&rankings); err != nil{
		return nil, err
	}

	return rankings, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will go here
		// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()

		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}

		favourite_genres, err := GetUserFavouriteGenres(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user's favourite genres"})
			return
		}

		err = godotenv.Load(".env")

		if err != nil {
			log.Println("Error loading .env file:", err)
		}

		var recommendedMovieLimitVal int64 = 5

		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIES_COUNT")

		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		filter := bson.M{"genre.genre_name": bson.M{"$in": favourite_genres}}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := movieCollection.Find(ctx, filter, findOptions.SetLimit(recommendedMovieLimitVal))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommended movies."})
			return
		}

		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie

		if err = cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommended movies."})
			return
		}
	  c.JSON(http.StatusOK, recommendedMovies)
	}
}


func GetUserFavouriteGenres(userId string) ([]string, error) {
	// Implementation will go here

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{"favourite_genres.genre_name": 1, "_id": 0}

	filterOptions := options.FindOne().SetProjection(projection)
	var result bson.M

	err := userCollection.FindOne(ctx, filter, filterOptions).Decode(&result)
	if err != nil {
		return nil, err
	}

	favouriteGenresRaw, ok := result["favourite_genres"].([]interface{})
	if !ok {
		return nil, errors.New("failed to cast favourite_genres")
	}

	var genreNames []string

	for _, genre := range favouriteGenresRaw {
		if genreMap, ok := genre.(map[string]interface{}); ok {
			if name, ok := genreMap["genre_name"].(string); ok {
				genreNames = append(genreNames, name)
			}
		}
	}
	return genreNames, nil
}
