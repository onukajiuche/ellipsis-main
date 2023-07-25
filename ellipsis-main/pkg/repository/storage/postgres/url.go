package postgres

import (
	"brief/internal/model"
	"context"

	"gorm.io/gorm/clause"
)

// CreateURL stores 'url' in the database
func (p *Postgres) CreateURL(ctx context.Context, url *model.URL) error {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()
	return db.Model(&model.User{ID: url.UserID}).Association("Urls").Append(url)
}

// GetURL fetches a url entry from the database using its 'hash'
func (p *Postgres) GetURL(ctx context.Context, hash string) (*model.URL, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	var url model.URL
	err := db.First(&url, "hash = ?", hash).Error
	return &url, err
}

// GetURL fetches a url entry from the database using its 'id'
func (p *Postgres) GetURLById(ctx context.Context, id string) (*model.URL, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	var url model.URL
	err := db.First(&url, "id = ?", id).Error
	return &url, err
}

// GetUrls fetches all url's made by a user with 'userID'
func (p *Postgres) GetUrls(ctx context.Context, userID string) ([]model.URL, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	user := model.User{ID: userID}
	var urls []model.URL

	err := db.Model(&user).Association("Urls").Find(&urls)
	return urls, err
}

// GetAll fetches all url's in the database
func (p *Postgres) GetAll(ctx context.Context) ([]model.URL, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	var urls []model.URL
	err := db.Find(&urls).Error
	return urls, err
}

// DeleteUrl deletes a random url by its 'id'
func (p *Postgres) DeleteUrl(ctx context.Context, id string) (*model.URL, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	url := model.URL{ID: id}
	err := db.Model(&url).Clauses(clause.Returning{}).Delete(&url).Error
	return &url, err
}
