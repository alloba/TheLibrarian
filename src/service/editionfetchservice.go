package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/alloba/TheLibrarian/logging"
	"path/filepath"
	"strconv"
)

type EditionFetchService struct {
	repoManager *repository.RepoManager
	fileService *FileService
}

type

func NewEditionFetchService(repoManager *repository.RepoManager, fileService *FileService) *EditionFetchService {
	return &EditionFetchService{repoManager: repoManager, fileService: fileService}
}

func (service EditionFetchService) DownloadEdition(bookName string, editionNum int, destinationFolder string) error {
	bookExist, err := service.repoManager.Book.ExistsByName(bookName)
	if err != nil {
		return logging.LogTrace(err)
	}
	if !bookExist {
		return logging.LogTrace(fmt.Errorf("cannot locate book %v", bookName))
	}

	book, err := service.repoManager.Book.ExistAndFetchByName(bookName)
	if err != nil {
		return logging.LogTrace(err)
	}

	editionExist, err := service.repoManager.Edition.ExistByBookIdAndEditionNumber(book.ID, editionNum)
	if err != nil {
		return logging.LogTrace(err)
	}
	if !editionExist {
		return logging.LogTrace(fmt.Errorf("cannot locate book %v", bookName))
	}

	edition, err := service.repoManager.Edition.FindByBookIdAndEditionNumber(book.ID, editionNum)
	if err != nil {
		return logging.LogTrace(err)
	}

	chapters, err := service.repoManager.Chapter.FindAllByEditionId(edition.ID)
	if err != nil {
		return logging.LogTrace(err)
	}

	for _, chapter := range *chapters {
		pagerecpair := make([]PageRecordPair, 0)
		chapterPages, err := service.repoManager.Page.FindAllByChapterId(chapter.ID)
		if err != nil {
			return logging.LogTrace(err)
		}

		for _, page := range *chapterPages {
			record, err := service.repoManager.Record.FindOne(page.RecordID)
			if err != nil {
				return logging.LogTrace(err)
			}
			pageMem := page
			pagerecpair = append(pagerecpair, PageRecordPair{page: &pageMem, record: record})
		}
		subfolderName := fmt.Sprintf("%v - ed[%v-%v] ch[%v]", book.Name, strconv.Itoa(edition.EditionNumber), edition.Name, chapter.Name)
		err = service.fileService.WritePageAssociationsToDestination(&pagerecpair, destinationFolder, subfolderName+string(filepath.Separator))
		if err != nil {
			return logging.LogTrace(err)
		}
	}

	return nil
}


