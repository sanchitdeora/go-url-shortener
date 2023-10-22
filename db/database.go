package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	// username = os.Getenv("DBUSER")
	// password = os.Getenv("DBPASS")
)

type UrlShortener struct {
    CompleteURL  string
    ShortenedURL string 
}

type Database interface {
	// Init() *sql.DB
	SaveUrl(c context.Context, urlShortener *UrlShortener) error
}

type Opts struct {
    Hostname string
    DBName string
    TableName string
    TimeoutSeconds time.Duration
    MaxOpenConns int
    MaxIdleConns int
    MaxConnsLifetime time.Duration
    DBConnection *sql.DB
}

type dbImpl struct {
    *Opts
}

func NewDatabase(opts *Opts) Database {
	return &dbImpl{Opts: opts}
}

func (db *dbImpl) SaveUrl(c context.Context, urlShortener *UrlShortener) error {
    query := `INSERT INTO urlshortener(completeurl, shortenedurl) VALUES (?, ?)`
    ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancelfunc()

    res, err := db.DBConnection.ExecContext(ctx, query, urlShortener.CompleteURL, urlShortener.ShortenedURL )
    if err != nil {
        log.Printf("Error %s when inserting row into products table", err)
        return err
    }

    rows, err := res.RowsAffected()
    if err != nil {
        log.Printf("Error %s when finding rows affected", err)
        return err
    }
    log.Printf("%d URLs created ", rows)
	return nil
}

