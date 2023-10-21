package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

const (  
    hostname = "127.0.0.1:3306"
    dbname   = "urlshortenerdb"
)

var (
	// username = os.Getenv("DBUSER")
	// password = os.Getenv("DBPASS")
)

type Database interface {
	Init()
	InsertUrl() error
}

type dbImpl struct {}

func NewDatabase() Database {
	return &dbImpl{}
}

func (db *dbImpl) InsertUrl() error {
	return nil
}

func (d *dbImpl) Init() {
	// establish db connection
	db, err := dbConnection()
    if err != nil {
        log.Error().Msg(fmt.Sprintf("Error %s while getting db connection", err))
        return
    }
    defer db.Close()

	// create url shortener table	
	err = createUrlShortenerTable(db)
    if err != nil {
        log.Printf("Create url shortener table failed with error %s", err)
        return
    }
}

func dbConnection() (*sql.DB, error) {  
    
	// open sql connection
    // Capture connection properties.
    cfg := mysql.Config{
        User:   "root",
        Passwd: "Password@1",
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
		AllowNativePasswords: true,
    }

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Info().Msg(fmt.Sprintf("Error %s when opening DB\n", err))
        return nil, err
    }

    ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelfunc()

	// create db if not exists
	_, err = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + dbname)
    if err != nil {
        log.Info().Msg(fmt.Sprintf("Error %s when creating DB\n", err))
        return nil, err
    }

    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(20)
    db.SetConnMaxLifetime(time.Minute * 5)

    db.Close()

	// connect to database
	cfg.DBName = dbname
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Printf("Error %s when opening DB", err)
        return nil, err
    }

    ctx, cancelfunc = context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancelfunc()

	err = db.PingContext(ctx)
    if err != nil {
        log.Info().Msg(fmt.Sprintf("Errors %s pinging DB", err))
        return nil, err
    }

    log.Printf("Connected to DB %s successfully\n", dbname)
    return db, nil
}

func createUrlShortenerTable(db *sql.DB) error {  
	// create table if not exist
	query := `CREATE TABLE IF NOT EXISTS urlshortener(
			id int primary key auto_increment,
			completeurl text,
			shortenedurl text
		)`
	
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)  
	defer cancelfunc()

	_, err := db.ExecContext(ctx, query)  
	if err != nil {  
		log.Printf("Error %s when creating url shortener table", err)
		return err
	}

	return nil
}

