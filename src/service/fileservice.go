package service

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type FileService struct {
	archiveBasePath string
}

func NewFileService(arhiveBasePath string) *FileService {
	return &FileService{
		archiveBasePath: arhiveBasePath,
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
	FileInfo        os.FileInfo
}

func (service FileService) getChildrenContainers(container *FileContainer) (*[]FileContainer, error) {
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
		x = append(x, *childContainer)
		return nil
	})
	if err != nil {
		return nil, logTrace(err)
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
		FileInfo:        fileInfo,
	}, nil
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
		return "", logTrace(err)
	}

	if stat.IsDir() {
		path = path + string(os.PathSeparator)
	}
	return path, nil
}
