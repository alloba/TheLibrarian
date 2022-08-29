package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/fileutil"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RecordService struct {
	recordRepo *database.RecordRepo
	basePath   string
}

func New(repo *database.RecordRepo, archiveBasePath string) *RecordService {
	absPath, err := filepath.Abs(archiveBasePath)
	if err != nil {
		log.Fatalf("could not load path for %v. does directory exist?", archiveBasePath)
	}
	stat, err := os.Stat(absPath)
	if err != nil || !stat.IsDir() {
		log.Fatalf("either couldnt read file stats for absPath or it isn't a directory.")
	}

	return &RecordService{
		recordRepo: repo,
		basePath:   absPath + string(os.PathSeparator),
	}
}

func (service RecordService) CreateRecordData(filepath string) (*database.Record, error) {
	file, err := fileutil.GetFileBinary(filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to get file binary - %v", err.Error())
	}
	defer file.Close()

	fileHash, err := fileutil.GetFileHash(file)
	if err != nil {
		return nil, fmt.Errorf("unable to get file hash - %v", err.Error())
	}

	qualifiedPath, err := fileutil.GetQualifiedFilePath(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not get full path of file - %v", err.Error())
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not load file stats - %v", err.Error())
	}

	var pathBits = strings.Split(qualifiedPath, string(os.PathSeparator))
	var fname = pathBits[len(pathBits)-1]
	var ext = "." + strings.Split(fname, ".")[len(strings.Split(fname, "."))-1]

	var recordObj = &database.Record{
		Hash:             fileHash,
		FilePointer:      service.basePath + generateUniqueSubPath(fileHash) + ".bin",
		Name:             fname,
		Extension:        ext,
		DateFileModified: stat.ModTime(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}

	return recordObj, nil
}

func (service RecordService) PersistRecordIfUnique(record *database.Record, originPath string) (bool, error) {
	recordExists, err := service.recordRepo.Exists(record.Hash)
	if err != nil {
		return false, fmt.Errorf("unable to query db for record exists - %v", err.Error())
	}
	if recordExists {
		return false, nil
	}

	file, err := fileutil.GetFileBinary(originPath)
	if err != nil {
		return false, fmt.Errorf("unable to get file binary - %v", err.Error())
	}
	defer file.Close()

	fileDestination, err := os.Create(record.FilePointer)
	if err != nil {
		return false, fmt.Errorf("unable to access record filepointer - %v", err.Error())
	}
	defer fileDestination.Close()

	_, err = io.Copy(fileDestination, file)
	if err != nil {
		return false, fmt.Errorf("unable to write bytes to record filepointer - %v", err.Error())
	}

	err = service.recordRepo.SaveOne(record)
	if err != nil {
		quietAttemptDelete(record.FilePointer)
		return false, fmt.Errorf("unable to save new record to database - %v", err.Error())
	}

	_, err = service.recordRepo.FindByHash(record.Hash)
	if err != nil {
		quietAttemptDelete(record.FilePointer)
		return false, fmt.Errorf("could not ack record in database: %v", record.Hash)
	}

	return true, nil
}

func generateUniqueSubPath(hash string) string {
	return fmt.Sprintf("%v_%v", time.Now().Unix(), hash)
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
