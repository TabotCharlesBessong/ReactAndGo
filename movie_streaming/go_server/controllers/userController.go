package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/database"
	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

func HashPassword(password string)(string, error){
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	if err != nil{
		return "", err
	}

	return string(HashPassword), nil
}

func RegisterUser () gin.HandlerFunc {
	return func(c *gin.Context) {	
		var user models.User

		if err:= c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details":err.Error()})
			return 
		}

		hashedPassword, err := HashPassword(user.Password)

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"UNable to hash password"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx,bson.M{"email":user.Email})

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to check existing users"})
			return
		}

		if count > 0{
			c.JSON(http.StatusConflict,gin.H{"error":"User already exist"})
			return
		}

		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(ctx,user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated,result)

	}
}