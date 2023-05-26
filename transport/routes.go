package transport

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	requestid "github.com/sumit-tembe/gin-requestid"
)

// Router func
func Router(newRelicApp *newrelic.Application) http.Handler {

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}
	router := gin.Default()
	router.RedirectTrailingSlash = true

	router.Use(nrgin.Middleware(newRelicApp))
	router.Use(cors.New(corsConfig))
	router.Use(requestid.RequestID(nil))

	// if os.Getenv("ENVIRONMENT") == "load" {
	// 	pprof.Register(router, "debug")
	// }

	router.GET("/", HealthCheck)
	router.GET("/healthz", HealthCheck)

	//Performance verify key on load forge
	loaderVerification := os.Getenv("LOAD_FORGE")
	if loaderVerification != "" {
		router.GET(fmt.Sprintf("/%v", loaderVerification), func(g *gin.Context) {
			g.Writer.WriteHeader(http.StatusOK)
			g.Writer.Write([]byte(loaderVerification))
		})
	}

	api := router.Group("api/v1")
	{
		api.GET("", HealthCheck)
		api.POST("/login", handleLogin)
		api.POST("/signup", handleSignup)

		members := api.Group("/members")
		{
			members.Use(AuthMiddleware()).POST("/bookmarks", uploadBookmarks)
			members.Use(AuthMiddleware()).POST("/settings")
			members.Use(AuthMiddleware()).GET("/settings")
			members.Use(AuthMiddleware()).GET("/bookmarks", getBookmarks)

		}

	}

	return router
}
