package service

import (
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/fileutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RecordService struct {
	recordRepo       *database.RecordRepo
	CreateRecordData processFile
}

type processFile func(filepath string) (*database.Record, error)

var basePath string

func New(repo *database.RecordRepo, archiveBasePath string) *RecordService {
	absPath, err := filepath.Abs(archiveBasePath)
	if err != nil {
		log.Fatalf("could not load path for %v. does directory exist?", archiveBasePath)
	}
	stat, err := os.Stat(absPath)
	if err != nil || !stat.IsDir() {
		log.Fatalf("either couldnt read file stats for absPath or it isn't a directory.")
	}
	basePath = absPath + string(os.PathSeparator)

	return &RecordService{
		recordRepo:       repo,
		CreateRecordData: createRecordFromFile,
	}
}

func createRecordFromFile(filepath string) (*database.Record, error) {
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
		FilePointer:      basePath + generateUniqueSubPath(fileHash) + ".bin",
		Name:             fname,
		Extension:        ext,
		DateFileModified: stat.ModTime(),
		DateCreated:      time.Now(),
		DateModified:     time.Now(),
	}

	//TODO actually copy file into destination here before returning the record.
	return recordObj, nil
}

func generateUniqueSubPath(hash string) string {
	return fmt.Sprintf("%v_%v", time.Now().Unix(), hash)
}
