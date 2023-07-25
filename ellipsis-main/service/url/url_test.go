// build+ unit
package url_test

import (
	"brief/internal/constant"
	"brief/internal/model"
	"brief/pkg/repository/storage"
	"brief/service/mock"
	"brief/service/url"
	"net/http"
	"strings"
	"testing"
	"time"
)

var mockStorage storage.StorageRepository = &mock.Repo{}
var storageService url.UrlService = url.NewUrlService(mockStorage)

func TestRedirect(t *testing.T) {
	hashString := "hashString"
	url, err := storageService.Redirect(hashString)
	if err != nil {
		t.Errorf("Expected 'error' to be nil, got '%v'", err)
	}

	if url == nil {
		t.Errorf("Expected 'url' to be not nil, got '%v'", url)
	}
}

func TestShorten(t *testing.T) {
	testLongUrl := "http://my-url.com"
	req, err := http.NewRequest("POST", testLongUrl, nil)
	if err != nil {
		t.Errorf("Expected 'error' to be nil when creating request, got '%v'", err)
	}

	t.Run("Specified Hash", func(t *testing.T) {
		hash := "specific"
		url := &model.URL{LongURL: "https://google.com", UserID: "test-id", Hash: hash}
		if err = storageService.Shorten(url, &model.ContextInfo{ID: "test-id"}, req); err != nil {
			t.Errorf("Expected 'error' to be nil when creating request, got '%v'", err)
		}

		if url.ID == "" {
			t.Errorf("Expected 'url.ID' to be not empty")
		}

		if !url.CreatedAt.Before(time.Now()) {
			t.Errorf("Expected 'url.CreatedAt' to be before the current time")
		}

		if !strings.HasPrefix(url.Hash, testLongUrl) {
			t.Errorf("Expected 'url.Hash' to have prefix '%v', got '%v'", testLongUrl, url.Hash)
		}

		if expVal := testLongUrl + "/" + hash; url.Hash != expVal {
			t.Errorf("Expected 'url.Hash' to be '%v', got '%v'", expVal, url.Hash)
		}

	})

	t.Run("Random Hash", func(t *testing.T) {
		url := &model.URL{LongURL: "https://google.com", UserID: "test-id"}
		if err = storageService.Shorten(url, &model.ContextInfo{ID: "test-id"}, req); err != nil {
			t.Errorf("Expected 'error' to be nil when creating request, got '%v'", err)
		}

		if url.ID == "" {
			t.Errorf("Expected 'url.ID' to be not empty")
		}

		if !url.CreatedAt.Before(time.Now()) {
			t.Errorf("Expected 'url.CreatedAt' to be before the current time")
		}

		if !strings.HasPrefix(url.Hash, testLongUrl) {
			t.Errorf("Expected 'url.Hash' to have prefix '%v', got '%v'", testLongUrl, url.Hash)
		}
	})

}

func TestDelete(t *testing.T) {
	uniformID := "test-id"
	t.Run("Authorized", func(t *testing.T) {
		_, err := storageService.Delete(&model.ContextInfo{ID: uniformID, Role: constant.Roles[constant.Admin]}, uniformID)
		if err != nil {
			t.Errorf("Expected 'error' to be nil, got '%v'", err)
		}

	})

	t.Run("Unauthorized", func(t *testing.T) {
		_, err := storageService.Delete(&model.ContextInfo{ID: uniformID, Role: constant.Roles[constant.User]}, "test-id-2")
		if err == nil {
			t.Errorf("Expected 'error' to be not nil")
		}
	})
}

func TestGetUrls(t *testing.T) {
	_, err := storageService.GetURLs("test-id")
	if err != nil {
		t.Errorf("Expected 'error' to be nil, got '%v'", err)
	}
}

func TestGetAll(t *testing.T) {
	_, err := storageService.GetAll()
	if err != nil {
		t.Errorf("Expected 'error' to be nil, got '%v'", err)
	}
}
