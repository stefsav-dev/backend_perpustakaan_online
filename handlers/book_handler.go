package handlers

import (
	"backend_perpustakaan_online/models"
	"backend_perpustakaan_online/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	bookRepo *repositories.BookRepository
}

func NewBookHandler() *BookHandler {
	return &BookHandler{
		bookRepo: repositories.NewBookRepository(),
	}
}

func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")
	category := c.Query("category")

	filter := repositories.BookFilter{
		Search:   search,
		Status:   status,
		Category: category,
		Page:     page,
		Limit:    limit,
	}

	books, pagination, err := h.bookRepo.GetAll(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Gagal mengambil data books",
		})
	}

	var bookResponses []models.BookResponse
	for _, book := range books {
		bookResponses = append(bookResponses, models.BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Author:      book.Author,
			ISBN:        book.ISBN,
			Description: book.Description,
			Category:    book.Category,
			TotalPages:  book.TotalPages,
			Publisher:   book.Publisher,
			PublisherAt: book.PublisherAt,
			Status:      book.Status,
			CreatedAt:   book.CreatedAt,
			UpdateAt:    book.UpdateAt,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookResponses,
		"meta":    pagination,
	})
}

func (h *BookHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Book ID",
		})
	}

	book, err := h.bookRepo.GetBookByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Buku Tidak ada",
		})
	}

	bookResponse := models.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		ISBN:        book.ISBN,
		Description: book.Description,
		Category:    book.Category,
		TotalPages:  book.TotalPages,
		Publisher:   book.Publisher,
		PublisherAt: book.PublisherAt,
		Status:      book.Status,
		CreatedAt:   book.CreatedAt,
		UpdateAt:    book.UpdateAt,
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookResponse,
	})
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var bookReq models.BookRequest
	if err := c.BodyParser(&bookReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Request Body",
		})
	}

	if bookReq.Title == "" || bookReq.Author == "" || bookReq.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Title, Author and ISBN are required",
		})
	}

	exists, err := h.bookRepo.CheckISBNExists(bookReq.ISBN, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal cek ISBN",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error":   "Buku dengan ISBN sudah ada",
		})
	}

	book := models.Book{
		Title:       bookReq.Title,
		Author:      bookReq.Author,
		ISBN:        bookReq.ISBN,
		Description: bookReq.Description,
		Category:    bookReq.Category,
		TotalPages:  bookReq.TotalPages,
		Publisher:   bookReq.Publisher,
		PublisherAt: bookReq.PublisherAt,
	}

	if bookReq.Status != "" {
		book.Status = bookReq.Status
	}

	if err := h.bookRepo.CreateDataBook(&book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal Membuat data buku",
		})
	}

	bookResponse := models.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		ISBN:        book.ISBN,
		Description: book.Description,
		Category:    book.Category,
		TotalPages:  book.TotalPages,
		Publisher:   book.Publisher,
		PublisherAt: book.PublisherAt,
		Status:      book.Status,
		CreatedAt:   book.CreatedAt,
		UpdateAt:    book.UpdateAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    bookResponse,
	})
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Book ID",
		})
	}

	existingBook, err := h.bookRepo.GetBookByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Buku tidak ada",
		})
	}

	var bookReq models.BookRequest
	if err := c.BodyParser(&bookReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Request Body",
		})
	}

	if bookReq.ISBN != "" && bookReq.ISBN != existingBook.ISBN {
		exists, err := h.bookRepo.CheckISBNExists(bookReq.ISBN, uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Gagal cek ISBN",
			})
		}
		if exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"error":   "Buku dengan ISBN sudah ada",
			})
		}
		existingBook.ISBN = bookReq.ISBN
	}

	if bookReq.Title != "" {
		existingBook.Title = bookReq.Title
	}
	if bookReq.Author != "" {
		existingBook.Author = bookReq.Author
	}
	if bookReq.Description != "" {
		existingBook.Description = bookReq.Description
	}
	if bookReq.Category != "" {
		existingBook.Category = bookReq.Category
	}
	if bookReq.TotalPages > 0 {
		existingBook.TotalPages = bookReq.TotalPages
	}
	if bookReq.Publisher != "" {
		existingBook.Publisher = bookReq.Publisher
	}
	if !bookReq.PublisherAt.IsZero() {
		existingBook.PublisherAt = bookReq.PublisherAt
	}
	if bookReq.Status != "" {
		existingBook.Status = bookReq.Status
	}

	if err := h.bookRepo.UpdateDataBook(existingBook); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "gagal update data buku",
		})
	}
	bookResponse := models.BookResponse{
		ID:          existingBook.ID,
		Title:       existingBook.Title,
		Author:      existingBook.Author,
		ISBN:        existingBook.ISBN,
		Description: existingBook.Description,
		Category:    existingBook.Category,
		TotalPages:  existingBook.TotalPages,
		Publisher:   existingBook.Publisher,
		PublisherAt: existingBook.PublisherAt,
		Status:      existingBook.Status,
		CreatedAt:   existingBook.CreatedAt,
		UpdateAt:    existingBook.UpdateAt,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookResponse,
	})
}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid Book ID",
		})
	}

	if _, err := h.bookRepo.GetBookByID(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Buku tidak ada",
		})
	}

	if err := h.bookRepo.DeleteDataBook(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "gagal delete data buku",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"success": true,
		"message": "Data Buku berhasil dihapus",
	})
}

func (h *BookHandler) UpdateBookStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid book ID",
		})
	}

	var request struct {
		Status string `json:"status" validate:"required,oneof=available borrowed maintenance"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	validStatuses := map[string]bool{
		"available":   true,
		"borrowed":    true,
		"maintenance": true,
	}

	if !validStatuses[request.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid status. Must be: available, borrowed, or maintenance",
		})
	}

	// Check if book exists
	if _, err := h.bookRepo.GetBookByID(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Book not found",
		})
	}

	if err := h.bookRepo.UpdateStatus(uint(id), request.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update book status",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Book status updated successfully",
		"status":  request.Status,
	})
}
