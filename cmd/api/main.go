package main

import (
	"Go_shortenURL/configs"
	"Go_shortenURL/internal/handler"
	"Go_shortenURL/internal/repository"
	"Go_shortenURL/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	config := configs.NewConfig()

	urlRepository := repository.NewURLRepository(config.Database, config.Redis)

	urlService := service.NewURLService(urlRepository)

	h := handler.NewURLHandler(urlService)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", h.IndexPage)
	router.POST("/shorten", h.ShortenURL)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/:shortCode", h.RedirectURL)

	router.Run(":" + config.Server.Port)
}
