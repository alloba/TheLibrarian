package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type EditionManagerService struct {
	repoManager *repository.RepoManager
	fileService *FileService
}

func NewEditionManagerService(repoManager *repository.RepoManager, fileService *FileService) *EditionManagerService {
	return &EditionManagerService{repoManager: repoManager, fileService: fileService}
}

func (service EditionManagerService) CreateNewEditionInNamedBook(bookName string, editionName string, systemPath string) (*database.Edition, error) {
	bookExist, err := service.repoManager.Book.ExistsByName(bookName)
	if err != nil {
		return nil, logTrace(err)
	}

	var book *database.Book
	if !bookExist {
		book = &database.Book{
			ID:           uuid.New().String(),
			Name:         bookName,
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}
	} else {
		foundBook, err := service.repoManager.Book.FindOneByName(bookName)
		if err != nil {
			return nil, logTrace(err)
		}
		book = foundBook
	}

	nextEdition, err := service.repoManager.Edition.FindNextEditionNumber(book.ID)
	if err != nil {
		return nil, logTrace(err)
	}
	var edition = &database.Edition{
		ID:            uuid.New().String(),
		Name:          editionName,
		EditionNumber: nextEdition,
		BookID:        book.ID,
		DateCreated:   time.Now(),
		DateModified:  time.Now(),
	}

	rootFileContainer, err := service.fileService.createFileContainer(systemPath)
	if err != nil {
		return nil, logTrace(err)
	}
	if !rootFileContainer.IsDir {
		return nil, fmt.Errorf("must provide a folder to create an edition, not a single file - %v", systemPath)
	}

	children, err := service.fileService.getChildrenContainers(rootFileContainer)
	if err != nil {
		return nil, logTrace(err)
	}

	records := make([]database.Record, 0)
	for _, child := range *children {
		rec, err := service.recordFromFileContainer(&child)
		if err != nil {
			return nil, logTrace(err)
		}
		if !child.IsDir {
			records = append(records, *rec)
		}
	}

	pages := make([]database.Page, 0)
	for _, record := range records {
		page := service.pageFromRecordAndEdition(&record, edition.ID)
		pages = append(pages, *page)
	}

	// it doesnt use the tx object directly, but this seems to work fine as a transaction.
	err = service.repoManager.Db.Transaction(func(tx *gorm.DB) error {
		for _, container := range *children {
			err := service.fileService.WriteContainerToArchive(&container)
			if err != nil {
				return logTrace(err)
			}
		}

		if !bookExist {
			err = service.repoManager.Book.CreateOne(book)
			if err != nil {
				return err
			}
		}

		err = service.repoManager.Record.UpsertAll(&records)
		if err != nil {
			return err
		}

		err = service.repoManager.Page.UpsertAll(&pages)
		if err != nil {
			return err
		}

		err = service.repoManager.Edition.CreateOne(edition)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, logTrace(err)
	}

	return edition, nil
}

func (service EditionManagerService) recordFromFileContainer(container *FileContainer) (*database.Record, error) {
	var record *database.Record
	exists, err := service.repoManager.Record.Exists(container.Hash)
	if err != nil {
		return nil, logTrace(err)
	}

	if exists {
		existingRecord, err := service.repoManager.Record.FindOne(container.Hash)
		if err != nil {
			return nil, logTrace(err)
		}
		record = existingRecord
	} else {
		record = &database.Record{
			ID:               container.Hash,
			FilePointer:      container.DestinationPath,
			Name:             container.OriginName,
			Extension:        container.OriginExt,
			DateFileModified: container.SourceFileInfo.ModTime(),
			DateCreated:      time.Now(),
			DateModified:     time.Now(),
		}
	}
	return record, nil
}

func (service EditionManagerService) pageFromRecordAndEdition(record *database.Record, editionId string) *database.Page {
	return &database.Page{
		ID:           uuid.New().String(),
		RecordID:     record.ID,
		EditionID:    editionId,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
}
