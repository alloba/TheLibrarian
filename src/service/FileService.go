package service

import (
	"crypto/sha256"
	"fmt"
	"github.com/alloba/TheLibrarian/database"
	"github.com/alloba/TheLibrarian/logging"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type FileService struct {
	archiveBasePath string
	existingHashes  []string
	archiveDirty    bool
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

func NewFileService(archiveBasePath string) *FileService {
	archiveFullPath, err := getQualifiedPath(archiveBasePath)
	if err != nil {
		panic(logging.LogTrace(err))
	}
	return &FileService{
		archiveBasePath: archiveFullPath,
		existingHashes:  make([]string, 0),
		archiveDirty:    true,
	}
}

func (service *FileService) WriteContainerToArchive(container *FileContainer) error {
	childContainers := make([]FileContainer, 0)
	if container.IsDir {
		childs, err := service.GetChildrenContainers(container, true)
		if err != nil {
			return logging.LogTrace(err)
		}
		childContainers = *childs
	} else {
		err := service.copyFileToArchive(container)
		if err != nil {
			return logging.LogTrace(err)
		}
	}

	for _, child := range childContainers {
		if !child.IsDir {
			err := service.copyFileToArchive(&child)
			if err != nil {
				return logging.LogTrace(err)
			}
		}
	}
	return nil
}

// Write the file represented by the passed in container to disk.
//
func (service *FileService) copyFileToArchive(container *FileContainer) error {
	if !container.SourceFileInfo.Mode().IsRegular() {
		return logging.LogTrace(fmt.Errorf("specified file is not regular [%v]", container.OriginPath))
	}

	exist, err := service.doesFileExistInArchive(container.Hash)
	if err != nil {
		return logging.LogTrace(err)
	}
	if exist {
		return nil // existing is allowed. just quietly skip.
	}

	err = copyFile(container.OriginPath, container.DestinationPath)
	if err != nil {
		return logging.LogTrace(err)
	}
	service.existingHashes = append(service.existingHashes, container.Hash)
	return nil
}

// Check the archive folder for anything that matches the existing hash.
// This relies on the stored files being named after their hash, to avoid needing to read all files in the directory repeatedly.
// So here is yet another good reason to never mess with archive file names (not to mention db row association).
func (service *FileService) doesFileExistInArchive(hash string) (bool, error) {
	if service.archiveDirty {
		err := service.calculateArchiveHashes()
		if err != nil {
			return false, logging.LogTrace(err)
		}
	}
	for _, hashStore := range service.existingHashes {
		if hash == hashStore {
			return true, nil
		}
	}
	return false, nil
}

func (service *FileService) calculateArchiveHashes() error {
	fmt.Printf("Calculating archive hashes\n")
	err := filepath.Walk(service.archiveBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash := strings.Split(info.Name(), ".")[0]
			service.existingHashes = append(service.existingHashes, hash)
		}
		return nil
	})
	if err != nil {
		return logging.LogTrace(fmt.Errorf("error while gathering hashes - %v", err.Error()))
	}
	service.archiveDirty = false
	return nil
}

// Create file container representation of the provided path.
// The file container is mainly oriented towards actual file objects,
//   with a flag that marks if it is a directory (this flag must be checked manually before file operations take place)
// This function represents a relatively expensive operation, since the file must be fully scanned to calculate a hash value for the container.
// Target destination is preemptively assigned to the container as well, using the archive path passed to the service during initialization.
func (service *FileService) CreateFileContainer(path string) (*FileContainer, error) {
	originPath, err := getQualifiedPath(path)
	if err != nil {
		return nil, logging.LogTrace(err)
	}

	fileInfo, err := os.Stat(originPath)
	if err != nil {
		return nil, logging.LogTrace(err)
	}

	originSplit := strings.Split(originPath, string(os.PathSeparator))
	var originName string
	if originSplit[len(originSplit)-1] == "" && len(originSplit) > 2 {
		originName = originSplit[len(originSplit)-2]
	} else {
		originName = originSplit[len(originSplit)-1]
	}
	extSplit := strings.Split(originName, ".")
	var originExt string //some files don't have an extension. check here.
	if len(extSplit) == 1 {
		originExt = ""
	} else {
		originExt = "." + extSplit[len(extSplit)-1]
	}

	hash := ""
	destName := ""
	destPath := ""
	destExt := ""
	// if the container is a folder, cannot get hash or destination. so skip.
	if !fileInfo.IsDir() {
		file, err := os.Open(originPath)
		if err != nil {
			return nil, logging.LogTrace(err)
		}
		defer func(file *os.File) {
			cerr := file.Close()
			if err == nil {
				err = cerr
			}
		}(file)
		hash, err = calculateHash(file)
		if err != nil {
			return nil, logging.LogTrace(err)
		}

		destName = hash + ".bin"

		destPath, err = getQualifiedPath(service.archiveBasePath)
		if err != nil {
			return nil, logging.LogTrace(err)
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

// Create containers for all files that exist within the directory specified by the input container.
// This is a comprehensive scan - all subdirectories are also examined for files.
// Chains together calls to base function for creating a single file container
func (service *FileService) GetChildrenContainers(container *FileContainer, onlyFiles bool) (*[]FileContainer, error) {
	paths := make([]string, 0)
	x := make([]FileContainer, 0)
	if !container.IsDir {
		return &x, logging.LogTrace(fmt.Errorf("the provided container does not represent a directory"))
	}

	err := filepath.Walk(container.OriginPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, logging.LogTrace(fmt.Errorf("could not form child containers - %v", err.Error()))
	}

	for index, pth := range paths {
		fmt.Printf("Creating source containers %v/%v\n", index, len(paths))
		childContainer, err := service.CreateFileContainer(pth)
		if err != nil {
			return nil, err
		}
		if (childContainer.IsDir && onlyFiles == false) || (!childContainer.IsDir && onlyFiles == true) {
			x = append(x, *childContainer)
		}
	}
	return &x, nil
}

// Read the file and calculate the sha256 hash.
// This is used when checking to see if the file already exists in the archive, and when creating records in the database.
func calculateHash(file *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", logging.LogTrace(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Paranoid version of path sanitization.
// Given a path string, create a normalized absolute path.
// If the path is a directory, ensure that a standard path separator is included on the end.
func getQualifiedPath(path string) (string, error) {
	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		usr, err := user.Current()
		if err != nil {
			return "", logging.LogTrace(err)
		}
		path = usr.HomeDir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		usr, err := user.Current()
		if err != nil {
			return "", logging.LogTrace(err)
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", logging.LogTrace(err)
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

// A supposedly "correct" version of writing files in go.
// Basically just adds error checking and deferred operations that aren't intuitive.
// This function should only be used as support functionality for filecontainer oriented functions.
func copyFile(sourcePath string, destinationPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return logging.LogTrace(err)
	}
	defer func(file *os.File) {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}(source)

	destination, err := os.Create(destinationPath)
	if err != nil {
		return logging.LogTrace(err)
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

func (service *FileService) GetSubpathOfContainer(sourceDir string, container *FileContainer) (string, error) {
	absOrigin, err := getQualifiedPath(sourceDir)
	if err != nil {
		return "", logging.LogTrace(err)
	}
	absFile, err := getQualifiedPath(container.OriginPath)
	if err != nil {
		return "", logging.LogTrace(err)
	}
	return strings.TrimPrefix(absFile, absOrigin), nil
}

func (service *FileService) DownloadPageRecord(destinationFolder string, book *database.Book, edition *database.Edition, page *database.Page, record *database.Record) error {
	qualifiedDestination, err := getQualifiedPath(destinationFolder)
	if err != nil {
		return logging.LogTrace(err)
	}
	qualifiedDestination = qualifiedDestination + book.Name + " - " + strconv.Itoa(edition.EditionNumber) + string(os.PathSeparator)

	fullPath := qualifiedDestination + page.RelativePath
	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(fullPath), os.ModeDir)
		if err != nil {
			return logging.LogTrace(err)
		}
	}
	err = copyFile(service.archiveBasePath+record.FilePointer, fullPath)
	if err != nil {
		return logging.LogTrace(err)
	}
	return nil
}
