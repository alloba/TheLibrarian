package service

import (
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/alloba/TheLibrarian/logging"
	"github.com/google/uuid"
	"time"
)

type PageService struct {
	pageRepo       *repository.PageRepo
	recordService  *RecordService
	editionService *EditionService
}

func NewPageService(pageRepo *repository.PageRepo, recordService *RecordService, editionService *EditionService) *PageService {
	if pageRepo == nil {
		panic("nil pagerepo for pageservice")
	}
	if recordService == nil {
		panic("nil recordservice for pageservice")
	}
	if editionService == nil {
		panic("nil editionservice for pageservice")
	}

	return &PageService{
		pageRepo:       pageRepo,
		recordService:  recordService,
		editionService: editionService,
	}
}

func (service PageService) CreatePage(recordId string, editionId string, relativeFilePath string) error {
	//verify that record and edition exist in db
	_, err := service.recordService.GetRecord(recordId)
	if err != nil {
		return logging.LogTrace(err)
	}
	_, err = service.editionService.GetEditionByID(editionId)
	if err != nil {
		return logging.LogTrace(err)
	}

	page := &database.Page{
		ID:           uuid.NewString(),
		RecordID:     recordId,
		EditionID:    editionId,
		RelativePath: relativeFilePath,
		DateCreated:  time.Now(),
		DateModified: time.Now(),
	}
	err = service.pageRepo.CreateOne(page)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (service PageService) FindByRecordAndEdition(recordId string, editionId string) (*database.Page, error) {
	page, err := service.pageRepo.FindOneByRecordAndEdition(recordId, editionId)
	if err != nil {
		return nil, logging.LogTrace(err)
	}

	return page, nil
}

func (service PageService) FindAllByEditionId(editionId string) (*[]database.Page, error) {
	pages, err := service.pageRepo.FindAllByEditionId(editionId)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return pages, nil
}
