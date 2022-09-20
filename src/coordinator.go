package main

import (
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/alloba/TheLibrarian/logging"
	"github.com/alloba/TheLibrarian/service"
	"gorm.io/gorm"
)

type ActionCoordinator struct {
	db          *gorm.DB
	archivePath string
}

type actionInstance struct {
	db             *gorm.DB
	pageService    *service.PageService
	editionService *service.EditionService
	recordService  *service.RecordService
	bookService    *service.BookService
	fileService    *service.FileService
	archivePath    string
}

func NewActionCoordinator(db *gorm.DB, archivePath string) *ActionCoordinator {
	return &ActionCoordinator{
		db:          db,
		archivePath: archivePath,
	}
}

func (service ActionCoordinator) SubmitNewEdition(bookName string, editionName string, sourceDirectory string) error {
	tx := service.db.Begin()
	transactionCoordinator := newActionInstance(tx, service.archivePath)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err := transactionCoordinator.submitNewEdition(bookName, editionName, sourceDirectory)
	if err != nil {
		tx.Rollback()
		return logging.LogTrace(err)
	}
	return tx.Commit().Error
}

func newActionInstance(db *gorm.DB, archivePath string) *actionInstance {
	pagerepo := repository.NewPageRepo(db)
	editionrepo := repository.NewEditionRepo(db)
	recordrepo := repository.NewRecordRepo(db)
	bookrepo := repository.NewBookRepo(db)

	recordservice := service.NewRecordService(recordrepo)
	bookservice := service.NewBookService(bookrepo)
	editionservice := service.NewEditionService(editionrepo, bookservice)
	pageservice := service.NewPageService(pagerepo, recordservice, editionservice)
	fileservice := service.NewFileService(archivePath)

	return &actionInstance{
		db:             db,
		pageService:    pageservice,
		editionService: editionservice,
		recordService:  recordservice,
		bookService:    bookservice,
		fileService:    fileservice,
		archivePath:    archivePath,
	}
}

func (service actionInstance) submitNewEdition(bookName string, editionName string, sourceDirectory string) error {
	directoryContainer, err := service.fileService.CreateFileContainer(sourceDirectory)
	if err != nil {
		return logging.LogTrace(err)
	}
	allFilesContainers, err := service.fileService.GetChildrenContainers(directoryContainer, true)
	if err != nil {
		return logging.LogTrace(err)
	}

	bookExist, err := service.bookService.ExistByName(bookName)
	if err != nil {
		return logging.LogTrace(err)
	}
	if !bookExist {
		err = service.bookService.CreateBook(bookName)
		if err != nil {
			return logging.LogTrace(err)
		}
	}
	book, err := service.bookService.GetBookByName(bookName)
	if err != nil {
		return logging.LogTrace(err)
	}
	err = service.editionService.CreateEdition(book.ID, editionName)
	if err != nil {
		return logging.LogTrace(err)
	}
	edition, err := service.editionService.GetMostRecentEditionForBook(book.ID)
	if err != nil {
		return logging.LogTrace(err)
	}

	for _, container := range *allFilesContainers {
		err = service.fileService.WriteContainerToArchive(&container)
		if err != nil {
			return logging.LogTrace(err)
		}
		recordExist, err := service.recordService.ExistsByID(container.Hash)
		if err != nil {
			return logging.LogTrace(err)
		}
		if !recordExist {
			err = service.recordService.CreateRecord(&container)
			if err != nil {
				return logging.LogTrace(err)
			}
		}
		record, err := service.recordService.GetRecord(container.Hash)
		if err != nil {
			return logging.LogTrace(err)
		}

		relativePath, err := service.fileService.GetSubpathOfContainer(sourceDirectory, &container)
		if err != nil {
			return logging.LogTrace(err)
		}
		err = service.pageService.CreatePage(record.ID, edition.ID, relativePath)
		if err != nil {
			return logging.LogTrace(err)
		}
	}
	//return fmt.Errorf("testfail")
	return nil
}

//
//func (service actionInstance) DownloadPage(pageId string, destinationFolder string) error   {}
//func (service actionInstance) DownloadEdition(bookName string, editionNum int, destinationFolder string) error {
//}
//func (service actionInstance) DownloadBook(bookName string, destinationFolder string) error {}

//func GetBookInformation(bookName string) {}
//func GetEditionInformation(bookName string, editionNum int){}
