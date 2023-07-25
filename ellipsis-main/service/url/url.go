package url

import (
	"brief/internal/config"
	"brief/internal/constant"
	"brief/internal/model"
	"brief/pkg/repository/storage"
	"brief/utility"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	urlPkg "net/url"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UrlService interface {
	Redirect(hash string) (*model.URL, error)
	Shorten(url *model.URL, ctxInfo *model.ContextInfo, r *http.Request) error
	Delete(ctxInfo *model.ContextInfo, urlId string) (*model.URL, error)
	GetURLs(userID string) ([]model.URL, error)
	GetAll() ([]model.URL, error)
}

type urlService struct {
	dbRepo storage.StorageRepository
}

func NewUrlService(dbRepo storage.StorageRepository) UrlService {
	return &urlService{dbRepo: dbRepo}
}

// Redirect contains business logic to redirect a shortened url to the original url
func (u *urlService) Redirect(hash string) (*model.URL, error) {

	url, err := u.dbRepo.GetURL(context.TODO(), hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("url not found")
		}
		return nil, fmt.Errorf("could not fetch url, got error %w", err)
	}

	return url, nil
}

// ADMIN & USER

// Link contains business logic to shorten and store a URL
func (u *urlService) Shorten(url *model.URL, ctxInfo *model.ContextInfo, r *http.Request) error {

	{
		// Check that URL is valid
		_, err := urlPkg.Parse(url.LongURL)
		if err != nil {
			return fmt.Errorf("invalid url specified: '%v'", url.LongURL)
		}

		if err := ping(url.LongURL); err != nil {
			return fmt.Errorf("invalid url specified: '%v', got error: '%v'", url.LongURL, err)
		}
	}

	// URL shortening logic
	url.ID = uuid.NewString()
	if ctxInfo != nil && ctxInfo.ID != "" {
		url.UserID = ctxInfo.ID
	} else {
		url.UserID = config.GetConfig().AdminID
	}

	if url.Hash == "" {
		// Run indefinite loop to prevent possible collision
		for {
			hash, err := utility.GetURLHash(url.ID, url.LongURL)
			if err != nil {
				return fmt.Errorf("could not generate hash, got error %w", err)
			}
			url.Hash = hash

			if err := u.dbRepo.CreateURL(context.TODO(), url); err != nil {
				if !errors.Is(err, gorm.ErrDuplicatedKey) {
					return fmt.Errorf("could not store url, got error %w", err)
				}
			} else {
				break
			}
		}
	} else {
		if err := u.dbRepo.CreateURL(context.TODO(), url); err != nil {
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return fmt.Errorf("could not store url, got error %w", err)
			}
			return fmt.Errorf("oops, '%s' already exists", url.Hash)
		}
	}

	hashUrl := urlPkg.URL{
		Host:   r.Host,
		Scheme: r.URL.Scheme,
		Path:   url.Hash,
	}
	if hashUrl.Scheme == "" {
		hashUrl.Scheme = "https"
	}
	url.Hash = hashUrl.String()
	return nil
}

// Delete contains business logic to delete a user's saved URL or a random url by its 'id'
func (u *urlService) Delete(ctxInfo *model.ContextInfo, urlId string) (*model.URL, error) {

	if ctxInfo.Role != constant.Roles[constant.Admin] {
		url, err := u.dbRepo.GetURLById(context.TODO(), urlId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("url not found")
			}
			return nil, fmt.Errorf("could not fetch url, got error %w", err)
		}

		if url.UserID != ctxInfo.ID {
			return nil, fmt.Errorf("unauthorized to perform this action")
		}
	}

	url, err := u.dbRepo.DeleteUrl(context.TODO(), urlId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("url not found")
		}
		return nil, fmt.Errorf("could not delete url, got error %w", err)
	}

	return url, nil
}

// GetURLs contains business logic to fetch all URL's created by a user with 'userID'
func (u *urlService) GetURLs(userID string) ([]model.URL, error) {

	urls, err := u.dbRepo.GetUrls(context.TODO(), userID)
	if err != nil {
		return nil, fmt.Errorf("could not get urls, got error : %w", err)
	}

	return urls, nil
}

// ADMIN

// GetAll contains business logic to fetch all URL's
func (u *urlService) GetAll() ([]model.URL, error) {

	urls, err := u.dbRepo.GetAll(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("could not get urls, got error : %w", err)
	}

	return urls, nil
}

func ping(url string) error {
	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{Timeout: 2 * time.Second}).DialContext,
		},
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	regx, _ := regexp.Compile("^20")
	ok := regx.Match([]byte(fmt.Sprint(resp.StatusCode)))
	if !ok {
		return fmt.Errorf("invalid")
	}
	return nil
}
