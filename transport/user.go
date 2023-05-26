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

// handleLogin handles the login request.
//
// This function receives a JSON payload containing the user's email and password,
// validates the request, compares the provided password with the stored hashed password,
// and generates a bearer token if the login is successful.
//
// If the request is invalid or the login fails, appropriate error responses are sent.
//
// Parameters:
// - g: The Gin context.
func handleLogin(g *gin.Context) {
	// Parse the JSON payload
	var user models.User
	if err := g.ShouldBindJSON(&user); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Retrieve the user from the database
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

// handleSignup handles the signup request.
//
// It extracts user information from the request body and performs the following steps:
// - Checks if the user already exists.
// - Generates a hashed password for the user.
// - Creates a new user record with the provided email and hashed password.
// - Inserts the new user record into the "users" table.
//
// If any error occurs during these steps, an appropriate error response is sent.
// If the signup process is successful, a success message is returned.
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

// generateBearerToken generates a bearer token for the given user.
//
// It takes a user model and creates a JWT token with custom claims based on the user information.
// The token is signed using the JWT secret from the configuration.
//
// The generated bearer token is returned as a string.
// If an error occurs during token generation, an error is returned.
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
