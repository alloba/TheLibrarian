package fileutil

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

func GetFileBinary(path string) (*os.File, error) {
	fullPath, err := GetQualifiedFilePath(path)
	if err != nil {
		return nil, fmt.Errorf("could not get filepath - %v", err.Error())
	}
	file, err := os.Open(fullPath)

	if err != nil {
		return nil, fmt.Errorf("unable to open file at path %v - %v", path, err.Error())
	}
	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		return nil, fmt.Errorf("cannot get file binary of directory %v - %v", path, err.Error())
	}
	return file, nil
}

func GetFileHash(file *os.File) (string, error) {

	var hasher = sha256.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		return "", fmt.Errorf("unable to get file hash for %v - %v", file.Name(), err.Error())
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func GetQualifiedFilePath(path string) (string, error) {
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
