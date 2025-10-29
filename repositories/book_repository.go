package repositories

import (
	"backend_perpustakaan_online/config"
	"backend_perpustakaan_online/models"
	"math"

	"gorm.io/gorm"
)

type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		DB: config.DB,
	}
}

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type BookFilter struct {
	Search   string
	Status   string
	Category string
	Page     int
	Limit    int
}

func (r *BookRepository) GetAll(filter BookFilter) ([]models.Book, *Pagination, error) {
	var books []models.Book
	var total int64

	query := r.DB.Model(&models.Book{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ?", search, search, search)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit < 1 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit
	totalPage := int(math.Ceil(float64(total) / float64(filter.Limit)))

	pagination := &Pagination{
		Page:      filter.Page,
		Limit:     filter.Limit,
		Total:     total,
		TotalPage: totalPage,
	}

	err := query.Offset(offset).Limit(filter.Limit).Order("created_at desc").Find(&books).Error
	if err != nil {
		return nil, nil, err
	}
	return books, pagination, nil
}

func (r *BookRepository) GetBookByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.DB.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) GetByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	err := r.DB.Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) CreateDataBook(book *models.Book) error {
	return r.DB.Create(book).Error
}

func (r *BookRepository) UpdateDataBook(book *models.Book) error {
	return r.DB.Save(book).Error
}

func (r *BookRepository) DeleteDataBook(id uint) error {
	return r.DB.Delete(&models.Book{}, id).Error
}

func (r *BookRepository) UpdateStatus(id uint, status string) error {
	return r.DB.Model(&models.Book{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BookRepository) CheckISBNExists(isbn string, excludeID uint) (bool, error) {
	var count int64
	query := r.DB.Model(&models.Book{}).Where("isbn = ?", isbn)

	if excludeID > 0 {
		query = query.Where("id > ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
