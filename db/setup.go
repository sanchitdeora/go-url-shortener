package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func Init(opts *Opts) *sql.DB {
	// establish db connection
	db, err := dbConnection(opts)
    if err != nil {
        log.Error().Msg(fmt.Sprintf("Error %s while getting db connection", err))
        return nil
    }

	// create url shortener table	
	err = createUrlShortenerTable(db, opts.TableName)
    if err != nil {
        log.Printf("Create url shortener table failed with error %s", err)
        return nil
    }

	return db
}

func dbConnection(opts *Opts) (*sql.DB, error) {  
    
	// open sql connection
    // Capture connection properties.
    cfg := mysql.Config{
        User:   "root",
        Passwd: "Password@1",
        Net:    "tcp",
        Addr:   opts.Hostname,
		AllowNativePasswords: true,
    }

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Info().Msg(fmt.Sprintf("Error %s when opening DB\n", err))
        return nil, err
    }

    ctx, cancelfunc := context.WithTimeout(context.Background(), opts.TimeoutSeconds)
    defer cancelfunc()

	// create db if not exists
	_, err = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + opts.DBName)
    if err != nil {
        log.Info().Msg(fmt.Sprintf("Error %s when creating DB\n", err))
        return nil, err
    }

    db.SetMaxOpenConns(opts.MaxOpenConns)
    db.SetMaxIdleConns(opts.MaxIdleConns)
    db.SetConnMaxLifetime(opts.MaxConnsLifetime)

    db.Close()

	// connect to database
	cfg.DBName = opts.DBName
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

    log.Printf("Connected to DB %s successfully\n", opts.DBName)
    return db, nil
}

func createUrlShortenerTable(db *sql.DB, tablename string) error {  
	// create table if not exist
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
			id int primary key auto_increment,
			completeurl text,
			shortenedurl text)`, tablename)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)  
	defer cancelfunc()

	_, err := db.ExecContext(ctx, query)  
	if err != nil {  
		log.Printf("Error %s when creating url shortener table", err)
		return err
	}

	return nil
}