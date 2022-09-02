package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/google/uuid"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type EditionManagerService struct {
	repoManager *repository.RepoManager
	fileService *FileService
}

func NewEditionManagerService(repoManager *repository.RepoManager, fileService *FileService) *EditionManagerService {
	return &EditionManagerService{repoManager: repoManager, fileService: fileService}
}

func (service EditionManagerService) CreateNewBook(bookName string) (*database.Book, error) {
	bookExist, err := service.repoManager.Book.ExistsByName(bookName)
	if err != nil {
		return nil, logTrace(err)
	}

	if bookExist {
		return nil, logTrace(fmt.Errorf("cannot create book that already exists with name %v", bookName))
	}

	var book = &database.Book{
		ID:           uuid.New().String(),
		Name:         bookName,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
	err = service.repoManager.Book.CreateOne(book)
	if err != nil {
		return nil, logTrace(err)
	}
	return book, nil
}

func (service EditionManagerService) CreateNewEdition(bookName string, chapterPaths ...string) (*database.Edition, error) {
	book, err := service.repoManager.Book.ExistAndFetchByName(bookName)
	if err != nil {
		return nil, logTrace(err)
	}

	nextEditionNum, err := service.repoManager.Edition.FindNextEditionNumber(book.ID)
	if err != nil {
		return nil, logTrace(err)
	}

	edition := &database.Edition{
		ID:            uuid.New().String(),
		Name:          "",
		EditionNumber: nextEditionNum,
		BookID:        book.ID,
		DateCreated:   time.Now(),
		DateModified:  time.Now(),
	}
	err = service.repoManager.Edition.CreateOne(edition)
	if err != nil {
		return nil, logTrace(err)
	}

	if len(chapterPaths) == 0 {
		return edition, nil
	}

	for _, path := range chapterPaths {
		_, err = service.CreateNewChapter(path, edition.ID)
		if err != nil {
			return nil, logTrace(err)
		}
	}

	return edition, nil
}

func (service EditionManagerService) CreateNewChapter(chapterPath string, editionId string) (*database.Chapter, error) {
	exist, err := service.repoManager.Edition.Exists(editionId)
	if err != nil {
		return nil, logTrace(err)
	}
	if !exist {
		return nil, logTrace(fmt.Errorf("edition [%v] does not exist", editionId))
	}

	fileContainer, err := service.fileService.createFileContainer(chapterPath)
	if err != nil {
		return nil, logTrace(err)
	}
	if !fileContainer.IsDir {
		return nil, logTrace(fmt.Errorf("cannot create chapter out of a file [%v]", chapterPath))
	}

	children, err := service.fileService.getChildrenContainers(fileContainer, true)
	if err != nil {
		return nil, logTrace(err)
	}

	dirPieces := strings.Split(fileContainer.OriginName, string(filepath.Separator))
	closestBareDirName := dirPieces[len(dirPieces)-1]
	closestBareDirName = strings.TrimSuffix(strings.TrimPrefix(closestBareDirName, string(filepath.Separator)), string(filepath.Separator))

	chapter := &database.Chapter{
		ID:           uuid.New().String(),
		EditionID:    editionId,
		RootPath:     closestBareDirName,
		Name:         "",
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
	err = service.repoManager.Chapter.CreateOne(chapter)
	if err != nil {
		return nil, logTrace(err)
	}

	for _, child := range *children {
		err := service.fileService.WriteContainerToArchive(&child)
		if err != nil {
			return nil, logTrace(err)
		}
		_, err = service.CreateNewPageAndAttachRecord(chapter, &child)
		if err != nil {
			return nil, logTrace(err)
		}
	}
	return chapter, nil
}

func (service EditionManagerService) CreateNewPageAndAttachRecord(chapter *database.Chapter, fileContainer *FileContainer) (*database.Page, error) {
	if fileContainer.IsDir {
		return nil, fmt.Errorf("cannot create page/record out of directory [%v]", fileContainer.OriginPath)
	}

	rec, err := service.CreateOrFindRecordForContainer(fileContainer)
	if err != nil {
		return nil, logTrace(err)
	}

	page := &database.Page{
		ID:           uuid.New().String(),
		RecordID:     rec.ID,
		ChapterID:    chapter.ID,
		RelativePath: strings.Split(fileContainer.OriginPath, string(filepath.Separator)+chapter.RootPath+string(filepath.Separator))[1], //hacky.
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}

	err = service.repoManager.Page.CreateOne(page)
	if err != nil {
		return nil, logTrace(err)
	}

	return page, nil

}

func (service EditionManagerService) CreateOrFindRecordForContainer(container *FileContainer) (*database.Record, error) {
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
		err := service.repoManager.Record.CreateOne(record)
		if err != nil {
			return nil, logTrace(err)
		}
	}
	return record, nil
}

func (service EditionManagerService) DownloadEdition(bookName string, editionNum int, destinationFolder string) error {
	bookExist, err := service.repoManager.Book.ExistsByName(bookName)
	if err != nil {
		return logTrace(err)
	}
	if !bookExist {
		return logTrace(fmt.Errorf("cannot locate book %v", bookName))
	}

	book, err := service.repoManager.Book.ExistAndFetchByName(bookName)
	if err != nil {
		return logTrace(err)
	}

	editionExist, err := service.repoManager.Edition.ExistByBookIdAndEditionNumber(book.ID, editionNum)
	if err != nil {
		return logTrace(err)
	}
	if !editionExist {
		return logTrace(fmt.Errorf("cannot locate book %v", bookName))
	}

	edition, err := service.repoManager.Edition.FindByBookIdAndEditionNumber(book.ID, editionNum)
	if err != nil {
		return logTrace(err)
	}

	chapters, err := service.repoManager.Chapter.FindAllByEditionId(edition.ID)
	if err != nil {
		return logTrace(err)
	}

	for _, chapter := range *chapters {
		pagerecpair := make([]PageRecordPair, 0)
		chapterPages, err := service.repoManager.Page.FindAllByChapterId(chapter.ID)
		if err != nil {
			return logTrace(err)
		}

		for _, page := range *chapterPages {
			record, err := service.repoManager.Record.FindOne(page.RecordID)
			if err != nil {
				return logTrace(err)
			}
			pageMem := page
			pagerecpair = append(pagerecpair, PageRecordPair{page: &pageMem, record: record})
		}
		subfolderName := fmt.Sprintf("%v - ed[%v-%v] ch[%v]", book.Name, strconv.Itoa(edition.EditionNumber), edition.Name, chapter.Name)
		err = service.fileService.WritePageAssociationsToDestination(&pagerecpair, destinationFolder, subfolderName+string(filepath.Separator))
		if err != nil {
			return logTrace(err)
		}
	}

	return nil
}
