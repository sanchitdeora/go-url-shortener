package main

import (
	// "net/http"

	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sanchitdeora/go-url-shortener/db"
	"github.com/sanchitdeora/go-url-shortener/urlshort"
)

const (
	PORT_NUMBER = ":8080"
)

func main() {
	// database initialization
	dbopts := &db.Opts{
		Hostname: "127.0.0.1:3306",
		DBName: "urlshortenerdb",
		TableName: "urlshortener",
		TimeoutSeconds: 5 * time.Second,
		MaxOpenConns: 20,
		MaxIdleConns: 20,
		MaxConnsLifetime: 5 * time.Minute,
	}
	dbConn := db.Init(dbopts)
	dbopts.DBConnection = dbConn
	defer dbConn.Close()

	db := db.NewDatabase(dbopts)

	// service initialization
	urlShortService := urlshort.NewUrlShortener(&urlshort.Opts{
		KeyLength: 10,
		ShortKeyDomainPrefix: "http://localhost" + PORT_NUMBER,
		DB: db,
	})

	// routers initialization
	startRouter(&urlshort.ApiService{UrlShortService: urlShortService})

}

func startRouter(service *urlshort.ApiService) {
	r := gin.Default()

	r.LoadHTMLGlob("static/templates/*")

	r.Use(static.Serve("/", static.LocalFile("./static", true)))
	r.GET("/ping", func(c *gin.Context) {
	  c.String(200, "test")
	})

	r.POST("/shorten", service.HandleUrlShortener)
	r.GET("/short/:id", service.HandleUrlRedirect)
	
	log.Info().Str("URL Shortener is listening on", PORT_NUMBER).Send()
	r.Run(PORT_NUMBER)
}