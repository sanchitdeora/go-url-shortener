package urlshort

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/rs/zerolog/log"
	"github.com/sanchitdeora/go-url-shortener/db"
)

//go:generate mockgen -destination=./mocks/mock_urlshort.go -package=mock_urlshort github.com/sanchitdeora/go-url-shortener/urlshort Service
type Service interface {
	UrlShortener(c context.Context, url string) (string, error)
	GetCompleteUrl(c context.Context, shortenedUrl string) (string, error)
	GetShortenedUrl(c context.Context, completeUrl string) (string, error)
}

type serviceImpl struct {
	*Opts
}

type Opts struct {
	KeyLength 				int
	ShortKeyDomainPrefix    string
	DB 			    		db.Database
}

func NewUrlShortener(opts *Opts) Service {
	return &serviceImpl{Opts: opts}
}

func (s *serviceImpl) GetCompleteUrl(c context.Context, shortenedKey string) (string, error) {
	completeUrl, err := s.DB.GetCompleteUrl(c, shortenedKey)
	if err != nil {
		log.Error().AnErr("Error while fetching completeUrl URL", err)
		return "", err
	}
	return completeUrl, nil
}

func (s *serviceImpl) GetShortenedUrl(c context.Context, completeUrl string) (string, error) {
	shortenedKey, err := s.DB.GetShortenedKey(c, completeUrl)
	if err != nil {
		log.Error().AnErr("Error while fetching shortened key", err)
		return "", err
	}
	return shortenedKey, nil
}

func (s *serviceImpl) UrlShortener(c context.Context, completeUrl string) (string, error) {

	shortenedKey, err := s.GetShortenedUrl(c, completeUrl)
	if err != nil {
		return "", err
	}
	if shortenedKey != "" {
		log.Info().Str("shortened kry already exists", shortenedKey).Send()
		return s.buildShortenedUrl(shortenedKey), nil
	}

	// Generate a unique shortened key for the original URL
	shortenedKey = createKey(completeUrl, s.KeyLength)

	err = s.DB.SaveUrl(c, completeUrl, shortenedKey)
	if err != nil {
		log.Error().AnErr("Error while saving URL", err)
		return "", err
	}

	return s.buildShortenedUrl(shortenedKey), nil
}

func createKey(url string, length int) (string) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	key := make([]byte, length)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}
	return string(key)
}

func (s *serviceImpl) buildShortenedUrl(key string) string {
	return fmt.Sprintf("%s/short/%s", s.ShortKeyDomainPrefix, key)
}