package transport

import (
	"log"
	"net/http"
	"strings"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwikportal-api/config"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware is a middleware function for JWT authentication.
// It verifies the bearer token provided in the request header and sets the custom claims in the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		jwtConfig := config.Cfg.JwtConfig

		// 1. Parse token and get token text
		tokenText, err := getTokenFromRequest(g.Request)
		if err != nil {
			g.AbortWithStatusJSON(http.StatusUnauthorized, "cannot get token from request")
			return
		}

		if tokenText == "" {
			log.Println("AuthMiddleware: Bearer token is required")
			g.AbortWithStatusJSON(http.StatusUnauthorized, "Bearer token is required")
			return
		}

		// 3. Get the token
		token, _ := jwt.ParseWithClaims(tokenText, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		})

		if token == nil {
			log.Printf("AuthMiddleware: Token is not parsable %v", tokenText)
			g.AbortWithStatusJSON(http.StatusUnauthorized, "token is not parsable")
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {

			// 4. Verify integrity of token
			err = VerifyClaims(claims, jwtConfig)
			if err != nil {
				log.Printf("AuthMiddleware: VerifyClaims token failed %v", claims)
				g.AbortWithStatusJSON(http.StatusUnauthorized, "invalid token")
				return
			}

			// 9. Set the auth context
			SetContext(ContextCustomClaims, *claims, g)
			return

		} else {
			log.Printf("Parsing token failed %v", token.Claims.Valid())
			g.AbortWithStatusJSON(http.StatusUnauthorized, "parsing token failure")
			return
		}
	}
}

// getTokenFromRequest extracts the bearer token from the request header.
// It returns the token string or an error if the token is not found or incorrect.
func getTokenFromRequest(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	splitToken := strings.Split(tokenString, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("auth token incorrect or wasn't supplied")
	}
	return strings.TrimSpace(splitToken[1]), nil
}

// VerifyClaims verifies the issuer and audience claims of the JWT token.
// It returns an error if the claims are invalid.
func VerifyClaims(claims *CustomClaims, jwtConfig *config.JWTConfig) error {
	if !claims.VerifyIssuer(jwtConfig.Issuer, true) {
		return errors.New("invalid JWT Issuer claim")
	}
	if !claims.VerifyAudience(jwtConfig.Audience, true) {
		return errors.New("invalid JWT Audience claim")
	}
	return nil
}

// SetContext sets the custom claims in the Gin context.
// It returns the modified Gin context.
func SetContext(name ContextKey, claims CustomClaims, g *gin.Context) *gin.Context {
	g.Set(string(name), claims)
	return g
}

// GetClaimsFromContext retrieves the custom claims from the Gin context.
// It returns a pointer to the custom claims.
func GetClaimsFromContext(g *gin.Context) *CustomClaims {
	claimsCtx := g.Value(string(ContextCustomClaims)).(CustomClaims)
	return &claimsCtx
}

// GetTokenFromRequest retrieves the bearer token from the request header.
// It returns the token string or an error if the token is not found or incorrect.
func GetTokenFromRequest(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	splitToken := strings.Split(tokenString, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("auth token incorrect or wasn't supplied")
	}
	return strings.TrimSpace(splitToken[1]), nil
}

// GetClaimsFromRequest retrieves the custom claims from the Gin context using the token in the request.
// It returns the custom claims or nil if the token is not valid or not provided.
func GetClaimsFromRequest(g *gin.Context) *CustomClaims {
	tokenText, err := GetTokenFromRequest(g.Request)
	if err != nil {
		return nil
	}
	token, _ := jwt.ParseWithClaims(tokenText, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JwtConfig.Secret), nil
	})
	claims, _ := token.Claims.(*CustomClaims)
	return claims
}

// GetUserIDFromRequest retrieves the user ID from the Gin context using the token in the request.
// It returns the user ID string or an empty string if the token is not valid or not provided.
func GetUserIDFromRequest(g *gin.Context) string {
	tokenText, err := GetTokenFromRequest(g.Request)
	if err != nil {
		return ""
	}
	token, _ := jwt.ParseWithClaims(tokenText, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JwtConfig.Secret), nil
	})
	claims, _ := token.Claims.(*CustomClaims)
	if claims != nil {
		return claims.UserID
	}
	return ""
}

// IsTokenExpired checks if the bearer token is expired based on the expiration claim.
// It returns true if the token is expired, false otherwise.
func IsTokenExpired(tokenString string, jwtConfig *config.JWTConfig) (isExpired bool) {
	claim := GetCustomClaimFromString(tokenString, jwtConfig)
	if claim != nil && claim.Expiration != nil {
		if claim.Expiration.Before(time.Now()) {
			isExpired = true
		}
	}
	return
}

// GetCustomClaimFromString parses the token string and retrieves the custom claims.
// It returns the custom claims or nil if the token is invalid.
func GetCustomClaimFromString(tokenString string, jwtConfig *config.JWTConfig) *CustomClaims {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})
	if err != nil {
		return nil
	}
	return token.Claims.(*CustomClaims)
}

// CustomClaims represents the custom claims in the JWT token.
type CustomClaims struct {
	jwt.StandardClaims
	Scope             string     `json:"scope"`
	Email             string     `json:"email"`
	UserID            string     `json:"user_id"`
	Expiration        *time.Time `json:"expiration"`
	SubscriberType    string     `json:"subscriber_type"`
	SubscriptionLevel string     `json:"subscription_level"`
}

// ContextKey represents the key for storing values in the Gin context.
type ContextKey string

var (
	ContextKeyAccountType = ContextKey("accountType")
	ContextCustomClaims   = ContextKey("customClaims")
	ContextKeyUUID        = ContextKey("UUID")
)
