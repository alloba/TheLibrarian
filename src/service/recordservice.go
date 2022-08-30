package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"io"
	"log"
	"os"
	"time"
)

type RecordService struct {
	recordRepo  *database.RecordRepo
	fileService *FileService
}

type RecordZip struct {
	OriginPath string
	RecordItem *database.Record
}

func NewRecordService(repo *database.RecordRepo, fileService *FileService) *RecordService {
	return &RecordService{
		recordRepo:  repo,
		fileService: fileService,
	}
}

func (service RecordService) CreateRecordData(filepath string) (*database.Record, error) {
	record, err := service.fileService.CreateFileObjectContainer(filepath)
	if err != nil || record.isDir {
		return nil, fmt.Errorf("unable to form file object for path %v", filepath)
	}
	defer record.Close()

	var recordObj = &database.Record{
		Id:               record.hash,
		FilePointer:      record.archivePath,
		Name:             record.name,
		Extension:        record.extension,
		DateFileModified: (*record.stat).ModTime(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}

	return recordObj, nil
}

func (service RecordService) PersistRecordIfUnique(recordZip *RecordZip) (bool, error) {
	recordExists, err := service.recordRepo.Exists(recordZip.RecordItem.Id)
	if err != nil {
		return false, fmt.Errorf("unable to query db for record exists - %v", err.Error())
	}
	if recordExists {
		return false, nil
	}

	filecontainer, err := service.fileService.CreateFileObjectContainer(recordZip.OriginPath)
	if err != nil || filecontainer.isDir {
		return false, fmt.Errorf("unable to get file container")
	}
	defer filecontainer.Close()

	fileDestination, err := os.Create(recordZip.RecordItem.FilePointer)
	if err != nil {
		return false, fmt.Errorf("unable to access record filepointer - %v", err.Error())
	}
	defer fileDestination.Close()

	_, err = io.Copy(fileDestination, filecontainer.binary)
	if err != nil {
		return false, fmt.Errorf("unable to write bytes to record filepointer - %v", err.Error())
	}

	err = service.recordRepo.SaveOne(recordZip.RecordItem)
	if err != nil {
		quietAttemptDelete(recordZip.RecordItem.FilePointer)
		return false, fmt.Errorf("unable to save new record to database - %v", err.Error())
	}

	_, err = service.recordRepo.FindByHash(recordZip.RecordItem.Id)
	if err != nil {
		quietAttemptDelete(recordZip.RecordItem.FilePointer)
		return false, fmt.Errorf("could not ack record in database: %v", recordZip.RecordItem.Id)
	}

	return true, nil
}

func (service RecordService) PersistAllRecordsIfUnique(records *[]RecordZip) (bool, error) {
	var allFullOps = false
	for _, item := range *records {
		fullAction, err := service.PersistRecordIfUnique(&item)
		if err != nil {
			return false, fmt.Errorf("unable to save all records, failed at %v - %v", item.OriginPath, err.Error())
		}
		if fullAction == false {
			allFullOps = false
		}
	}
	return allFullOps, nil
}

func (service RecordService) Exists(hash string) (bool, error) {
	exist, err := service.recordRepo.Exists(hash)
	if err != nil {
		return false, fmt.Errorf("could not check if record exists %v - %v", hash, err.Error())
	}
	return exist, nil
}

func (service RecordService) GetByHash(hash string) (*database.Record, error) {
	exist, err := service.Exists(hash)
	if err != nil {
		return nil, fmt.Errorf("could not check if file exists -%v", err.Error())
	}
	if !exist {
		return nil, fmt.Errorf("no record found %v", hash)
	}

	rec, err := service.recordRepo.FindByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to load record %v - %v", hash, err.Error())
	}

	return rec, nil
}

// an attempt at file cleanup when other operations fail and the copy operation needs to be backed out.
// a failure to delete something from the provided path should be rare (hopefully non-existent), and isn't something that
// can be easily handled other than to shout about it.
func quietAttemptDelete(path string) {
	e := os.Remove(path)
	if e != nil {
		log.Printf("WARNING - Orphan file at '%v'", path)
	}
}
