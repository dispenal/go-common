package common_utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func GetCurrentDir() string {
	dir, _ := os.Getwd()
	return dir
}

func GetExt(pathOrFilename string) string {
	return mime.TypeByExtension(filepath.Ext(pathOrFilename))
}

func CreateFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func CheckIfFileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DownloadFile(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func FetchFile(url, filePath string) (*os.File, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file: %s", response.Status)
	}

	tmpFile, err := ioutil.TempFile("", "downloaded")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}

	return tmpFile, nil
}
