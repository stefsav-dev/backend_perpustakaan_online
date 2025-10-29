package models

import "time"

type Book struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Author      string    `json:"author" gorm:"type:varchar(255);not null"`
	ISBN        string    `json:"isbn" gorm:"type:varchar(20);uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category" gorm:"type:varchar(100);"`
	TotalPages  int       `json:"total_pages" gorm:"type:int;"`
	Publisher   string    `json:"publisher" gorm:"type:varchar(255);"`
	PublisherAt time.Time `json:"publisher_at" gorm:"type:date"`
	Status      string    `json:"status" gorm:"type:enum('available','borrowed','maintenance');default:'available';'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type BookRequest struct {
	Title       string    `json:"title" validate:"required"`
	Author      string    `json:"author" validate:"required"`
	ISBN        string    `json:"isbn" validate:"required"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	TotalPages  int       `json:"total_pages"`
	Publisher   string    `json:"publisher"`
	PublisherAt time.Time `json:"publisher_at"`
	Status      string    `json:"status"`
}
type BookResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	TotalPages  int       `json:"total_pages"`
	Publisher   string    `json:"publisher"`
	PublisherAt time.Time `json:"publisher_at"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}
