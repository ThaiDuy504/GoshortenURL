package main

import (
	"net/http"
	"os"

	"Go_shortenURL/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	for _, f := range []string{".env", "configs/config.env"} {
		if err := godotenv.Load(f); err == nil {
			break
		}
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handler.IndexPage)
	router.POST("/shorten", handler.ShortenURL)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/:shortCode", handler.RedirectURL)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
