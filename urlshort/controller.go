package urlshort

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiService struct {
	UrlShortService Service
}

func (service *ApiService) HandleUrlShortener(c *gin.Context) {

	completeUrl := c.Request.FormValue("url")
	if completeUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "URL parameter is missing", 
		})
		return
	}

	shortenedUrl, err := service.UrlShortService.UrlShortener(c, completeUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err, 
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"shortenedUrl": shortenedUrl,
	})
}