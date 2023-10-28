package urlshort_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock_db "github.com/sanchitdeora/go-url-shortener/db/mocks"
	"github.com/sanchitdeora/go-url-shortener/urlshort"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_DOMAIN = "http://example.com"
	TEST_COMPLETE_URL = "http://test.com/complete-url/it-is-very-long/so-decide-to-shorten-it"
)

var ErrSomeError = errors.New("some error")

type ServiceMocks struct {
	mockDB mock_db.MockDatabase
}

func createUrlShortService(ctrl *gomock.Controller) (urlshort.ApiService, *ServiceMocks) {
	mockDb := mock_db.NewMockDatabase(ctrl)
	return urlshort.ApiService{UrlShortService: urlshort.NewUrlShortener(&urlshort.Opts{
		KeyLength: 10,
		ShortKeyDomainPrefix: TEST_DOMAIN,
		DB: mockDb,
	})}, 
	&ServiceMocks{mockDB: *mockDb} 
}

// HandleUrlShortener

func TestHandleUrlShortener_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mocks := createUrlShortService(ctrl)

	{	// happy path
		expectedUrl := fmt.Sprintf("%s/short/", TEST_DOMAIN)

		mocks.mockDB.EXPECT().
			GetShortenedKey(gomock.Any(), gomock.Any()).
			Return("", nil).
			Times(1)

		mocks.mockDB.EXPECT().
			SaveUrl(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		url, err := service.UrlShortService.UrlShortener(context.Background(), TEST_COMPLETE_URL)

		assert.True(t, strings.Contains(url, expectedUrl))
		assert.Nil(t, err)
	}
	{	// happy path with url found in db
		expectedKey := "abcdefghij"
		expectedUrl := fmt.Sprintf("%s/short/%s", TEST_DOMAIN, expectedKey)

		mocks.mockDB.EXPECT().
			GetShortenedKey(gomock.Any(), gomock.Any()).
			Return(expectedKey, nil).
			Times(1)

		url, err := service.UrlShortService.UrlShortener(context.Background(), TEST_COMPLETE_URL)

		assert.Equal(t, expectedUrl, url)
		assert.Nil(t, err)
	}
}

func TestHandleUrlShortener_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mocks := createUrlShortService(ctrl)

	{	// error while saving db
		mocks.mockDB.EXPECT().
			GetShortenedKey(gomock.Any(), gomock.Any()).
			Return("", nil)

		mocks.mockDB.EXPECT().
			SaveUrl(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(ErrSomeError)

		url, err := service.UrlShortService.UrlShortener(context.Background(), TEST_COMPLETE_URL)

		assert.Empty(t, url)
		assert.Equal(t, ErrSomeError, err)
	}
	{	// error while fetching shortened key from db
		mocks.mockDB.EXPECT().
			GetShortenedKey(gomock.Any(), gomock.Any()).
			Return("", ErrSomeError)

		url, err := service.UrlShortService.UrlShortener(context.Background(), TEST_COMPLETE_URL)

		assert.Empty(t, url)
		assert.Equal(t, ErrSomeError, err)
	}
}

func TestHandleGetCompleteUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service, mocks := createUrlShortService(ctrl)

	testShortenedKey := "abcdefghij"
	{	// happy path
		mocks.mockDB.EXPECT().
			GetCompleteUrl(gomock.Any(), gomock.Any()).
			Return(TEST_COMPLETE_URL, nil)

		url, err := service.UrlShortService.GetCompleteUrl(context.Background(), testShortenedKey)

		assert.Equal(t, TEST_COMPLETE_URL, url)
		assert.Nil(t, err)
	}
	{	// error while fetching complete url from db
		mocks.mockDB.EXPECT().
			GetCompleteUrl(gomock.Any(), gomock.Any()).
			Return("", ErrSomeError)

		url, err := service.UrlShortService.GetCompleteUrl(context.Background(), testShortenedKey)

		assert.Empty(t, url)
		assert.Equal(t, ErrSomeError, err)
	}
}