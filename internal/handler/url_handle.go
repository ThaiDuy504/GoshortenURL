package handler

import (
	"fmt"
	"net/http"
	"os"

	"Go_shortenURL/pkg/shortener"

	"github.com/gin-gonic/gin"
)

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func ShortenURL(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Please enter a URL.",
		})
		return
	}

	shortCode := shortener.Encode(url)
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"ShortURL":    fmt.Sprintf("%s/%s", baseURL, shortCode),
		"OriginalURL": url,
	})
}

func RedirectURL(c *gin.Context) {
	// shortCode := c.Param("shortCode")
	url := "https://www.google.com"
	if url == "" {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"Error": "URL not found.",
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url)
}
