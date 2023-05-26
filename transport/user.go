package transport

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jasonbronson/kwikportal-api/config"
	"github.com/jasonbronson/kwikportal-api/library"
	"github.com/jasonbronson/kwikportal-api/models"
	"github.com/jasonbronson/kwikportal-api/repositories"
	"golang.org/x/crypto/bcrypt"
)

func handleLogin(g *gin.Context) {

	var user models.User
	if err := g.ShouldBindJSON(&user); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	userDB, err := repositories.GetUser(user.Email)
	if err != nil {
		responseError(g, fmt.Errorf("Failed to find a user account %v", err))
		return
	}
	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		responseError(g, fmt.Errorf("Invalid Credentials %v", err))
		return
	}

	log.Println(user)
	// Assuming the login is successful, generate a bearer token
	token, err := generateBearerToken(userDB)
	if err != nil {
		responseError(g, fmt.Errorf("Failed to generate bearer token %v", err))
		return
	}

	responseSuccess(g, "token", token)
}

func handleSignup(g *gin.Context) {
	var user models.User
	if err := g.ShouldBindJSON(&user); err != nil {
		responseError(g, err)
		return
	}

	// Check if the user already exists
	existingUser, err := repositories.GetUser(user.Email)
	if err != nil {
		responseError(g, errors.New(fmt.Sprintf("User lookup failed %v", err)))
		return
	}
	if existingUser.Email != "" {
		responseError(g, errors.New("User already exists"))
		return
	}
	password := library.GeneratePassword(user.Password)

	// Create the new user record
	newUser := models.User{
		Email:    user.Email,
		Password: string(password),
	}

	// Insert the new user record into the "users" table
	err = repositories.SaveUser(newUser)
	if err != nil {
		responseError(g, err)
		return
	}

	responseSuccess(g, "message", "User created successfully")
}

func generateBearerToken(user models.User) (string, error) {

	jwtConfig := config.Cfg.JwtConfig
	uuid, _ := uuid.NewV4()
	claims := CustomClaims{
		jwt.StandardClaims{
			Audience: jwtConfig.Audience,
			Id:       uuid.String(),
			IssuedAt: time.Now().Unix(),
			Issuer:   jwtConfig.Issuer,
		},
		"user",
		user.Email,
		user.ID,
		nil,
		"free",
		"none",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", err
	}
	return signedString, nil
}
