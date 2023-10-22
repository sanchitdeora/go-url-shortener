package urlshort

import (
	"context"
	"math/rand"

	"github.com/rs/zerolog/log"
	"github.com/sanchitdeora/db"
)

type Service interface {
	UrlShortener(c context.Context, url string) (string, error)
}

type serviceImpl struct {
	*Opts
}

type Opts struct {
	KeyLength 		int
	ShortKeyPrefix  string
	DB 			    db.Database
}

func NewUrlShortener(opts *Opts) Service {
	return &serviceImpl{Opts: opts}
}

func (s *serviceImpl) UrlShortener(c context.Context, completeUrl string) (string, error) {
	
	// Generate a unique shortened key for the original URL
	shortenedKey := createKey(completeUrl, s.KeyLength)

	log.Printf(shortenedKey)

	shortenedUrl := s.ShortKeyPrefix + shortenedKey
	err := s.DB.SaveUrl(c, &db.UrlShortener{CompleteURL: completeUrl, ShortenedURL: shortenedUrl})
	if err != nil {
		log.Error().AnErr("Error while saving URL", err)
		return "", err
	}

	return shortenedUrl, nil
}

func createKey(url string, length int) (string) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	key := make([]byte, length)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}
	return string(key)
}