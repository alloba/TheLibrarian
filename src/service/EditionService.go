package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/alloba/TheLibrarian/logging"
	"github.com/google/uuid"
	"time"
)

type EditionService struct {
	editionRepo *repository.EditionRepo
	bookService *BookService
}

func NewEditionService(editionRepo *repository.EditionRepo, bookService *BookService) *EditionService {
	if editionRepo == nil {
		panic("nil edition repo in edition service")
	}
	if bookService == nil {
		panic("nil book service in edition service")
	}
	return &EditionService{
		editionRepo: editionRepo,
		bookService: bookService,
	}
}

func (service EditionService) CreateEdition(bookId string, editionName string) error {
	bookExist, err := service.bookService.ExistsByID(bookId)
	if err != nil {
		return logging.LogTrace(err)
	}
	if !bookExist {
		return logging.LogTrace(fmt.Errorf("cannot create edition for book that does not exist %v", bookId))
	}

	nextEditionNum, err := service.editionRepo.FindNextEditionNumber(bookId)
	if err != nil {
		return logging.LogTrace(err)
	}

	edition := &database.Edition{
		ID:            uuid.NewString(),
		Name:          editionName,
		EditionNumber: nextEditionNum,
		BookID:        bookId,
		DateCreated:   time.Now(),
		DateModified:  time.Now(),
	}

	err = service.editionRepo.CreateOne(edition)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (service EditionService) GetMostRecentEditionForBook(bookId string) (*database.Edition, error) {
	bookExist, err := service.bookService.ExistsByID(bookId)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	if !bookExist {
		return nil, logging.LogTrace(fmt.Errorf("book with id %v does not exist", bookId))
	}

	nextEdition, err := service.editionRepo.FindNextEditionNumber(bookId)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	edition, err := service.editionRepo.FindByBookIdAndEditionNumber(bookId, nextEdition-1) //should never be 0, so subtraction is fine
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return edition, nil
}

func (service EditionService) GetEditionByID(editionId string) (*database.Edition, error) {
	exist, err := service.editionRepo.Exists(editionId)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	if !exist {
		return nil, logging.LogTrace(fmt.Errorf("edition does not exist %v", editionId))
	}

	edition, err := service.editionRepo.FindByID(editionId)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return edition, nil
}
