package database

import (
	"backend_perpustakaan_online/config"
	"backend_perpustakaan_online/models"
	"log"
	"time"
)

func Migrate() {
	err := config.DB.AutoMigrate(
		&models.Book{},
	)

	if err != nil {
		log.Fatal("Gagal migrasi ke database", err)
	}

	log.Println("Database migrasi berhasil")
}

func Seeder() {
	books := []models.Book{
		{
			Title:       "The Great Gatsby",
			Author:      "F. Scott Fitzgerald",
			ISBN:        "978-0743273565",
			Description: "A classic novel of the Jazz Age",
			Category:    "Fiction",
			TotalPages:  180,
			Publisher:   "Scribner",
			PublisherAt: parseDate("1925-04-10"),
			Status:      "available",
		},
		{
			Title:       "To Kill a Mockingbird",
			Author:      "Harper Lee",
			ISBN:        "978-0061120084",
			Description: "A novel about racial inequality",
			Category:    "Fiction",
			TotalPages:  281,
			Publisher:   "J.B. Lippincott & Co.",
			PublisherAt: parseDate("1960-07-11"),
			Status:      "available",
		},
		{
			Title:       "1984",
			Author:      "George Orwell",
			ISBN:        "978-0452284234",
			Description: "Dystopian social science fiction",
			Category:    "Science Fiction",
			TotalPages:  328,
			Publisher:   "Secker & Warburg",
			PublisherAt: parseDate("1949-06-08"),
			Status:      "borrowed",
		},
	}

	for _, book := range books {
		var existingBook models.Book
		if err := config.DB.Where("isbn = ?", book.ISBN).First(&existingBook).Error; err != nil {
			if err := config.DB.Create(&book).Error; err != nil {
				log.Printf("Failed to create book: %v", err)
			}
		}
	}

	log.Println("Database seeding completed!")
}

func parseDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return time.Now()
	}
	return t
}
