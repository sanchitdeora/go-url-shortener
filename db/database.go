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

type Database interface {
    GetShortenedKey(c context.Context, completeUrl string) (string, error)
    GetCompleteUrl(c context.Context, shortenedUrl string) (string, error)
    SaveUrl(c context.Context, completeUrl, shortenedUrl string) error
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

func (db *dbImpl) SaveUrl(c context.Context, completeUrl, shortenedUrl string) error {
    query := `INSERT INTO urlshortener(completeurl, shortenedurl) VALUES (?, ?)`
    ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancelfunc()

    res, err := db.DBConnection.ExecContext(ctx, query, completeUrl, shortenedUrl)
    if err != nil {
        log.Printf("Error %s when inserting row into urlshortener table", err)
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

func (db *dbImpl) GetShortenedKey(c context.Context, completeUrl string) (string, error) {
    query := `SELECT shortenedurl FROM urlshortener WHERE completeurl = ?`
    ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancelfunc()

    row := db.DBConnection.QueryRowContext(ctx, query, completeUrl)
    var shortenedURL string
    if err := row.Scan(&shortenedURL); err != nil {
        if err == sql.ErrNoRows {
            return "", nil
        }
        return "", err
    }

    return shortenedURL, nil
}

func (db *dbImpl) GetCompleteUrl(c context.Context, shortenedUrl string) (string, error) {
    query := `SELECT completeurl FROM urlshortener WHERE shortenedurl = ?`
    ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancelfunc()

    row := db.DBConnection.QueryRowContext(ctx, query, shortenedUrl)
    var completeUrl string
    if err := row.Scan(&completeUrl); err != nil {
        if err == sql.ErrNoRows {
            return "", nil
        }
        return "", err
    }

    return completeUrl, nil
}