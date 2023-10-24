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
 
	c.HTML(http.StatusOK, "shorten.tmpl", gin.H{
		"title": "URL Shortener",
		"completeUrl": completeUrl,
		"shortenedUrl": shortenedUrl,
	})
}

func (service *ApiService) HandleUrlRedirect(c *gin.Context) {
	shortenedKey := c.Param("id")

	completeUrl, err := service.UrlShortService.GetCompleteUrl(c, shortenedKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err, 
		})
		return
	}
	if completeUrl == "" {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "no url found", 
		})
		return
	}

	c.Redirect(http.StatusFound, completeUrl)
}