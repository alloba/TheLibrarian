package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/alloba/TheLibrarian/logging"
	"github.com/google/uuid"
	"time"
)

type BookService struct {
	bookRepo *repository.BookRepo
}

func NewBookService(bookRepo *repository.BookRepo) *BookService {
	if bookRepo == nil {
		panic("null book repo")
	}
	return &BookService{bookRepo: bookRepo}
}

func (service BookService) CreateBook(bookName string) error {
	exist, err := service.bookRepo.ExistsByName(bookName)
	if err != nil {
		return logging.LogTrace(err)
	}
	if exist {
		return logging.LogTrace(fmt.Errorf("book with bookName %v already exists", bookName))
	}

	book := &database.Book{
		ID:           uuid.NewString(),
		Name:         bookName,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	err = service.bookRepo.CreateOne(book)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (service BookService) GetBookByName(bookName string) (*database.Book, error) {
	exist, err := service.bookRepo.ExistsByName(bookName)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	if !exist {
		return nil, fmt.Errorf("no book with name %v found", bookName)
	}

	book, err := service.bookRepo.FindOneByName(bookName)

	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return book, nil
}

func (service BookService) GetBookByID(bookId string) (*database.Book, error) {
	exist, err := service.bookRepo.Exists(bookId)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	if !exist {
		return nil, fmt.Errorf("no book with id %v found", bookId)
	}

	book, err := service.bookRepo.FindOneByID(bookId)

	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return book, nil
}
func (service BookService) DeleteBookByName(bookName string) error {
	exist, err := service.bookRepo.ExistsByName(bookName)
	if err != nil {
		return logging.LogTrace(err)
	}
	if !exist {
		return logging.LogTrace(fmt.Errorf("no book with name %v found", bookName))
	}
	book, err := service.bookRepo.FindOneByName(bookName)
	if err != nil {
		return logging.LogTrace(err)
	}

	err = service.bookRepo.DeleteOne(book.ID)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (service BookService) DeleteBookByID(bookId string) error {
	exist, err := service.bookRepo.Exists(bookId)
	if err != nil {
		return logging.LogTrace(err)
	}
	if !exist {
		return logging.LogTrace(fmt.Errorf("no book with id %v found", bookId))
	}

	err = service.bookRepo.DeleteOne(bookId)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}
