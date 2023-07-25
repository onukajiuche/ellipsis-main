package user_test

import (
	"brief/pkg/repository/storage"
	"brief/service/mock"
	"brief/service/url"
)

var mockStorage storage.StorageRepository = &mock.Repo{}
var storageService url.UrlService = url.NewUrlService(mockStorage)
