package service

import (
	"crypto/sha256"
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type FileService struct {
	archiveBasePath string
}

func NewFileService(archiveBasePath string) *FileService {
	return &FileService{
		archiveBasePath: archiveBasePath,
	}
}

type FileContainer struct {
	Hash            string
	IsDir           bool
	OriginPath      string
	OriginName      string
	OriginExt       string
	DestinationPath string
	DestinationName string
	DestinationExt  string
	SourceFileInfo  os.FileInfo
}

func (service FileService) WriteContainerToArchive(container *FileContainer) error {
	childContainers := make([]FileContainer, 0)
	if container.IsDir {
		childs, err := service.getChildrenContainers(container, true)
		if err != nil {
			return logTrace(err)
		}
		childContainers = *childs
	} else {
		err := service.copyFileToArchive(container)
		if err != nil {
			return logTrace(err)
		}
	}

	for _, child := range childContainers {
		if !child.IsDir {
			err := service.copyFileToArchive(&child)
			if err != nil {
				return logTrace(err)
			}
		}
	}
	return nil
}

func (service FileService) checkArchiveForExistinHash(hash string) (bool, error) {
	var exist = true
	err := filepath.Walk(service.archiveBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		//fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		if !info.IsDir() {
			if hash == strings.Split(info.Name(), ".")[0] {
				exist = false //TODO: should additionally verify the actual hash of the found file.
			}
		}
		return nil
	})
	if err != nil {
		return false, logTrace(fmt.Errorf("could not form child containers - %v", err.Error()))
	}
	return !exist, nil //todo i flipped this logic from an inverse operation, i should rewrite things to assume false instead of true for clarity.
}

func (service FileService) copyFileToArchive(container *FileContainer) error {
	if !container.SourceFileInfo.Mode().IsRegular() {
		return logTrace(fmt.Errorf("specified file is not regular [%v]", container.OriginPath))
	}

	exist, err := service.checkArchiveForExistinHash(container.Hash)
	if err != nil {
		return logTrace(err)
	}
	if exist {
		return nil // im not going to throw an error for not copying a file that already supposedly exists. just skip it.
		// TODO: technically possible to mismatch filename vs actual hashed file contents. should verify.
	}

	err = copyFile(container.OriginPath, container.DestinationPath)
	if err != nil {
		return logTrace(err)
	}
	return nil
}

func copyFile(sourcePath string, destinationPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return logTrace(err)
	}
	defer func(file *os.File) {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}(source)

	destination, err := os.Create(destinationPath)
	if err != nil {
		return logTrace(err)
	}
	defer func(file *os.File) {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}(destination)

	_, err = io.Copy(destination, source)
	return nil
}

func (service FileService) getChildrenContainers(container *FileContainer, onlyFiles bool) (*[]FileContainer, error) {
	x := make([]FileContainer, 0)
	if !container.IsDir {
		return &x, nil
	}

	err := filepath.Walk(container.OriginPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		//fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		childContainer, err := service.createFileContainer(path)
		if err != nil {
			return err
		}
		if (childContainer.IsDir && onlyFiles == false) || (!childContainer.IsDir && onlyFiles == true) {
			x = append(x, *childContainer)
		}
		return nil
	})
	if err != nil {
		return nil, logTrace(fmt.Errorf("could not form child containers - %v", err.Error()))
	}
	return &x, nil
}

func (service FileService) createFileContainer(path string) (*FileContainer, error) {
	originPath, err := getQualifiedPath(path)
	if err != nil {
		return nil, logTrace(err)
	}

	fileInfo, err := os.Stat(originPath)
	if err != nil {
		return nil, logTrace(err)
	}

	originSplit := strings.Split(originPath, string(os.PathSeparator))
	var originName string
	if originSplit[len(originSplit)-1] == "" && len(originSplit) > 2 {
		originName = originSplit[len(originSplit)-2]
	} else {
		originName = originSplit[len(originSplit)-1]
	}
	extSplit := strings.Split(originName, ".")
	var originExt string //some files dont have an extension. check here.
	if len(extSplit) == 1 {
		originExt = ""
	} else {
		originExt = "." + extSplit[len(extSplit)-1]
	}

	// if the container is a folder, cannot get hash or destination
	hash := ""
	destName := ""
	destPath := ""
	destExt := ""
	if !fileInfo.IsDir() {
		file, err := os.Open(originPath)
		if err != nil {
			return nil, logTrace(err)
		}
		defer func(file *os.File) {
			cerr := file.Close()
			if err == nil {
				err = cerr
			}
		}(file)
		hash, err = calculateHash(file)
		if err != nil {
			return nil, logTrace(err)
		}

		destName = hash + ".bin"

		destPath, err = getQualifiedPath(service.archiveBasePath)
		if err != nil {
			return nil, logTrace(err)
		}

		destPath = destPath + destName

		destExt = ".bin"
	}

	return &FileContainer{
		Hash:            hash,
		IsDir:           fileInfo.IsDir(),
		OriginPath:      originPath,
		OriginName:      originName,
		OriginExt:       originExt,
		DestinationPath: destPath,
		DestinationName: destName,
		DestinationExt:  destExt,
		SourceFileInfo:  fileInfo,
	}, nil
}

type PageRecordPair struct {
	page   *database.Page
	record *database.Record
}

func (service FileService) WritePageAssociationsToDestination(pageRecPairs *[]PageRecordPair, destination string, subfolder string) error {
	//verify directory exists
	qPath, err := getQualifiedPath(destination)
	if err != nil {
		return logTrace(err)
	}

	destStats, err := os.Stat(qPath)
	if err != nil {
		return logTrace(err)
	}
	if !destStats.IsDir() {
		return logTrace(fmt.Errorf("destination is not a directory [%v]", qPath))
	}

	subfolderQualified, err := getQualifiedPath(qPath + subfolder)
	if err != nil {
		return logTrace(err)
	}
	err = os.MkdirAll(subfolderQualified, 755) //mkdir all does not fail on existing directories
	if err != nil {
		return logTrace(err)
	}

	for _, pair := range *pageRecPairs {
		source := pair.record.FilePointer
		dest := subfolderQualified + string(filepath.Separator) + pair.page.RelativePath
		lastIndexOfSep := strings.LastIndex(dest, string(filepath.Separator))
		err = os.MkdirAll(dest[:lastIndexOfSep], 755)
		if err != nil {
			return logTrace(err)
		}
		err = copyFile(source, dest)
		if err != nil {
			return logTrace(err)
		}
	}
	return nil
}

func calculateHash(file *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", logTrace(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func getQualifiedPath(path string) (string, error) {
	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		usr, err := user.Current()
		if err != nil {
			return "", logTrace(err)
		}
		path = usr.HomeDir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		usr, err := user.Current()
		if err != nil {
			return "", logTrace(err)
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", logTrace(err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return path, nil //dont fail on a path that doesnt exist. but if it is there, continue on to see if a path separator should be appended as well.
	}
	if stat.IsDir() {
		path = path + string(os.PathSeparator)
	}
	return path, nil
}
