package service

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type FileService struct {
	basePath string
}

type FileObjectContainer struct {
	binary        *os.File
	stat          *os.FileInfo
	archivePath   string
	qualifiedPath string
	name          string
	extension     string
	isDir         bool
	hash          string
}

func (container FileObjectContainer) Close() {
	container.binary.Close()
}

func NewFileService(basePath string) *FileService {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		log.Fatalf("could not load path for %v. does directory exist?", basePath)
	}
	stat, err := os.Stat(absPath)
	if err != nil || !stat.IsDir() {
		log.Fatalf("either couldnt read file stats for absPath or it isn't a directory.")
	}

	if strings.Split(absPath, "")[len(strings.Split(absPath, ""))-1] != string(os.PathSeparator) {
		absPath = absPath + string(os.PathSeparator)
	}

	return &FileService{basePath: absPath}
}

func (service FileService) CreateFileObjectContainer(path string) (*FileObjectContainer, error) {
	absPath, err := service.getQualifiedFilePath(path)
	if err != nil {
		return nil, fmt.Errorf("could not format file path %v - %v", path, err.Error())
	}

	stat, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("could not get file stats - %v", err.Error())
	}

	binary, err := service.getFileBinary(absPath)
	if err != nil && !stat.IsDir() {
		return nil, fmt.Errorf("could not get file binary - %v", err.Error())
	}

	name, err := service.getFileName(absPath)
	if err != nil {
		return nil, fmt.Errorf("could not get file name - %v", err.Error())
	}

	ext, err := service.getFileExtension(absPath)
	if err != nil {
		return nil, fmt.Errorf("could not get file extension - %v", err.Error())
	}

	var hash = ""
	if !stat.IsDir() {
		hashAck, err := service.getFileHash(binary)
		if err != nil {
			return nil, fmt.Errorf("could not get file hash - %v", err.Error())
		}
		hash = hashAck
	}

	return &FileObjectContainer{
		binary:        binary,
		stat:          &stat,
		archivePath:   service.createRecordFilePath(hash),
		qualifiedPath: absPath,
		name:          name,
		extension:     ext,
		isDir:         stat.IsDir(),
		hash:          hash,
	}, nil
}

func (service FileService) getFileBinary(path string) (*os.File, error) {
	fullPath, err := service.getQualifiedFilePath(path)
	if err != nil {
		return nil, fmt.Errorf("could not get filepath - %v", err.Error())
	}
	file, err := os.Open(fullPath)

	if err != nil {
		return nil, fmt.Errorf("unable to open file at path %v - %v", path, err.Error())
	}
	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		return nil, fmt.Errorf("cannot get file binary of directory %v", path)
	}
	return file, nil
}

func (service FileService) SaveToArchiveBin(zip *RecordZip) (bool, error) {

	filecontainer, err := service.CreateFileObjectContainer(zip.OriginPath)
	if err != nil {
		fmt.Errorf("couldnt do the thing - %v", err.Error())
	}

	fileDestination, err := os.Create(zip.RecordItem.FilePointer)
	if err != nil {
		return false, fmt.Errorf("unable to access record filepointer - %v", err.Error())
	}
	defer fileDestination.Close()

	_, err = io.Copy(fileDestination, filecontainer.binary)
	if err != nil {
		return false, fmt.Errorf("unable to write bytes to record filepointer - %v", err.Error())
	}

	//err = service.recordRepo.SaveOne(recordZip.RecordItem)
	//if err != nil {
	//	quietAttemptDelete(recordZip.RecordItem.FilePointer)
	//	return false, fmt.Errorf("unable to save new record to database - %v", err.Error())
	//}
	return true, nil
}

func (service FileService) getFileHash(file *os.File) (string, error) {
	var hasher = sha256.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		return "", fmt.Errorf("unable to get file hash for %v - %v", file.Name(), err.Error())
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func (service FileService) getQualifiedFilePath(path string) (string, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if path == "~" {
		return dir, nil
	} else {
		finalPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("could not form path from %v - %v", path, err.Error())
		}
		return finalPath, nil
	}
}

func (service FileService) getFileName(path string) (string, error) {
	qualifiedPath, err := service.getQualifiedFilePath(path)
	if err != nil {
		return "", fmt.Errorf("could not get qualified path %v - %v", path, err.Error())
	}
	var pathBits = strings.Split(qualifiedPath, string(os.PathSeparator))
	var fname = pathBits[len(pathBits)-1]
	return fname, nil
}

func (service FileService) getFileExtension(path string) (string, error) {
	name, err := service.getFileName(path)
	if err != nil {
		return "", fmt.Errorf("could not get filename - %v", err.Error())
	}

	split := strings.Split(name, ".")
	if len(split) == 1 { //some files have no dot for the extension. the name and the extension are the same.
		return split[0], nil
	} else {
		return split[len(strings.Split(name, "."))-1], nil
	}
}

func (service FileService) createRecordFilePath(hash string) string {
	return service.basePath + hash + ".bin"
}

func (service FileService) GetAllNestedFilePaths(dirPath string) (*[]string, error) {
	directoryPath, err := service.getQualifiedFilePath(dirPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load dirPath %v - %v", dirPath, err.Error())
	}
	stats, err := os.Stat(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load stats for %v - %v", directoryPath, err.Error())
	}

	if !stats.IsDir() {
		return nil, fmt.Errorf("provided path is not a directory %v", directoryPath)
	}

	var paths = make([]string, 0)
	err = filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("could not walk directory %v - %v", directoryPath, err.Error())
	}

	return &paths, nil
}
