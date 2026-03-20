package handler

import (
	"log"
	"net/http"

	"Go_shortenURL/internal/service"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	URLService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{URLService: urlService}
}

func (h *URLHandler) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Please enter a URL.",
		})
		return
	}

	shortCode, err := h.URLService.ShortenURL(url)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"Error": "Failed to shorten URL.",
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"ShortURL": shortCode,
		"OriginalURL": url,
	})
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	originalURL, err := h.URLService.GetURL(shortCode)
	if err != nil {
		c.HTML(http.StatusNotFound, "index.html", gin.H{
			"Error": "URL not found.",
		})
		log.Fatal(err)
		return
	}
	log.Println(originalURL)
	// c.Redirect(http.StatusMovedPermanently, originalURL)
	//testing
	c.Redirect(http.StatusFound, "https://www.google.com")
}
