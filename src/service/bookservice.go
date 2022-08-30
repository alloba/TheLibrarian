package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/google/uuid"
	"time"
)

//TODO check all structs and make sure the dependencies are private level, not public. only the methods should be public.
type BookService struct {
	RecordService *RecordService
	FileService   *FileService
	BookRepo      *database.BookRepo
	PageRepo      *database.PageRepo
	EditionRepo   *database.EditionRepo
}

func NewBookService(recordService *RecordService, fileService *FileService, bookRepo *database.BookRepo, editionRepo *database.EditionRepo) *BookService {
	return &BookService{
		RecordService: recordService,
		FileService:   fileService,
		BookRepo:      bookRepo,
		EditionRepo:   editionRepo,
	}
}

func (service BookService) CreateNewBook(bookName string, directoryPath string) (*database.Edition, error) {
	bookExists, err := service.BookRepo.ExistsByName(bookName)
	if err != nil {
		return nil, fmt.Errorf("could not look for existing book - %v", err.Error())
	}
	if bookExists {
		return nil, fmt.Errorf("book with name %v already exists", bookName)
	}

	fileContainer, err := service.FileService.CreateFileObjectContainer(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("could not create file container - %v", err.Error())
	}
	defer fileContainer.Close()

	newRecords, err := service.getRecordList(fileContainer)
	if err != nil {
		return nil, fmt.Errorf("couldnt form records - %v", err.Error())
	}

	book := database.Book{
		Id:           uuid.New().String(),
		Name:         bookName,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
	edition := database.Edition{
		Id:            uuid.New().String(),
		EditionNumber: 0,
		BookId:        book.Id,
		Book:          &book,
		DateCreated:   time.Now(),
		DateModified:  time.Now(),
		Pages:         nil,
	}
	newPages := make([]database.Page, 0)
	for _, rec := range *newRecords {
		p := database.Page{
			Id:           uuid.New().String(),
			RecordId:     rec.RecordItem.Id,
			Record:       *rec.RecordItem,
			EditionId:    edition.Id,
			Edition:      nil,
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}
		newPages = append(newPages, p)
	}
	edition.Pages = &newPages

	for _, rec := range *newRecords {
		service.FileService.SaveToArchiveBin(&rec)
	}

	err = service.EditionRepo.SaveOne(&edition)
	if err != nil {
		return nil, fmt.Errorf("could not save edition to db - %v", err.Error())
	}
	ackEdition, err := service.EditionRepo.FindOne(edition.Id)
	if err != nil {
		return nil, fmt.Errorf("could not ack edition %v - %v", edition.Id, err.Error())
	}
	return ackEdition, nil
}

func (service BookService) getRecordList(fileContainer *FileObjectContainer) (*[]RecordZip, error) {
	var records = make([]RecordZip, 0)

	if !fileContainer.isDir {
		record, err := service.formRecord(fileContainer)
		if err != nil {
			return nil, fmt.Errorf("could not form record - %v", err.Error())
		}
		records = append(records, RecordZip{
			OriginPath: fileContainer.qualifiedPath,
			RecordItem: record,
		})
	} else {
		allFileNames, err := service.FileService.GetAllNestedFilePaths(fileContainer.qualifiedPath)
		if err != nil {
			return nil, fmt.Errorf("couldnt get files from directory - %v", err.Error())
		}
		for _, fname := range *allFileNames {
			subContainer, err := service.FileService.CreateFileObjectContainer(fname)
			if err != nil {
				return nil, fmt.Errorf("Could not process file %v - %v", fname, err.Error())
			}
			r, err := service.formRecord(subContainer)
			if err != nil {
				return nil, fmt.Errorf("could not form record %v - %v", fname, err.Error())
			}
			records = append(records, RecordZip{
				OriginPath: subContainer.qualifiedPath,
				RecordItem: r,
			})
		}
	}
	return &records, nil
}

func (service BookService) formRecord(fileContainer *FileObjectContainer) (*database.Record, error) {
	if fileContainer.isDir {
		return nil, fmt.Errorf("cannot form record from folder")
	}
	recordExists, err := service.RecordService.Exists(fileContainer.hash)
	if err != nil {
		return nil, fmt.Errorf("could not check if record exists %v - %v", fileContainer.hash, err.Error())
	}

	if recordExists {
		rec, err := service.RecordService.GetByHash(fileContainer.hash)
		if err != nil {
			return nil, fmt.Errorf("could not load existing record %v - %v", fileContainer.hash, err.Error())
		}

		return rec, nil

	} else {
		return &database.Record{
			Id:               fileContainer.hash,
			FilePointer:      fileContainer.archivePath,
			Name:             fileContainer.name,
			Extension:        fileContainer.extension,
			DateFileModified: (*fileContainer.stat).ModTime(),
			DateCreated:      time.Now(),
			DateModified:     time.Now(),
		}, nil
	}
}
