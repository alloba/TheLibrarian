package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/database/repository"
	"github.com/alloba/TheLibrarian/logging"
	"time"
)

type RecordService struct {
	RecordRepo *repository.RecordRepo
}

func NewRecordService(recordRepo *repository.RecordRepo) *RecordService {
	if recordRepo == nil {
		panic("cannot pass nil record repo")
	}
	return &RecordService{RecordRepo: recordRepo}
}

// CreateRecord will insert a new row into the database, based off of the provided FileContainer metadata.
// Attempting to write a duplicate record will fail (although in the future it might be nice to simply skip the operation)
func (service RecordService) CreateRecord(fileContainer *FileContainer) error {
	exist, err := service.RecordRepo.Exists(fileContainer.Hash)
	if err != nil {
		return logging.LogTrace(err)
	}
	if exist {
		return logging.LogTrace(fmt.Errorf("file with hash %v already exists in archive", fileContainer.Hash))
	}

	record := database.Record{
		ID:               fileContainer.Hash,
		FilePointer:      fileContainer.DestinationName,
		Name:             fileContainer.OriginName,
		Extension:        fileContainer.OriginExt,
		DateFileModified: time.Now(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}

	err = service.RecordRepo.CreateOne(&record)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

// GetRecord will return the record from the database with the ID that matches the provided file hash.
// If no such record exists, an error will be returned.
func (service RecordService) GetRecord(fileHash string) (*database.Record, error) {
	exist, err := service.RecordRepo.Exists(fileHash)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	if !exist {
		return nil, logging.LogTrace(fmt.Errorf("no record found for hash %v", fileHash))
	}

	record, err := service.RecordRepo.FindOne(fileHash)
	if err != nil {
		return nil, logging.LogTrace(err)
	}
	return record, nil
}

// DeleteRecord removes an entry from the database based off the provided file hash.
// Attempting to delete a record that does not exist can be considered a no-op.
func (service RecordService) DeleteRecord(fileHash string) error {
	exist, err := service.RecordRepo.Exists(fileHash)
	if err != nil {
		return logging.LogTrace(err)
	}
	if !exist {
		return nil
	}

	err = service.RecordRepo.DeleteOne(fileHash)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}

func (service RecordService) ExistsByID(recordId string) (bool, error) {
	exist, err := service.RecordRepo.Exists(recordId)
	if err != nil {
		return false, logging.LogTrace(err)
	}
	return exist, nil
}
