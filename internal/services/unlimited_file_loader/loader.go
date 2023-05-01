package unlimited_file_loader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

const (
	loaderName = "unlimited file loader"
)

type FileLoader struct {
}

// Download allow to download file by file link and save it into root file di
func (f *FileLoader) Download(fileLink string, fileRootPath string) error {
	fileUrl, err := url.Parse(fileLink)
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't parse file link", loaderName), err)
	}

	fileName := filepath.Base(fileUrl.Path)
	isExists, err := checkIfFileExists(fileName, fileRootPath)
	if err != nil {
		return err
	}

	resp, err := http.Get(fileLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if isExists {
		return f.continueFileLoading(resp, filepath.Join(fileRootPath, fileName))
	} else {
		return f.createAndLoadFile(resp, filepath.Join(fileRootPath, fileName))
	}
}

func (f *FileLoader) createAndLoadFile(resp *http.Response, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't file to load data", loaderName), err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	_, copyingErr := io.CopyBuffer(file, resp.Body, buffer)
	if copyingErr != nil {
		return copyingErr
	}

	return nil
}

func (f *FileLoader) continueFileLoading(resp *http.Response, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't open file [path: %s] to load data", loaderName, filePath), err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't take file statistics [path: %s]", loaderName, filePath), err)
	}

	currentFileSize := fileInfo.Size()
	factFileSizeInString := resp.Header.Get("Content-Length")
	factFileSize, err := strconv.ParseInt(factFileSizeInString, 10, 64)
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't parse content length from response [response data: %s]", loaderName, factFileSizeInString), err)
	}

	if currentFileSize == factFileSize {
		return nil
	}

	reader := bufio.NewReader(resp.Body)
	_, err = reader.Discard(int(currentFileSize))
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't skipe loaded bytes", loaderName), err)
	}

	buffer := make([]byte, 1024)
	for {
		bytesRead, readErr := reader.Read(buffer)
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return readErr
		}

		_, writeErr := file.Write(buffer[:bytesRead])
		if writeErr != nil {
			return writeErr
		}
	}

	return nil
}

func checkIfFileExists(fileName string, fileRootPath string) (bool, error) {
	fileDir, err := os.Open(fileRootPath)
	if err != nil {
		return false, errors.Join(fmt.Errorf("%s: can't open dirrectory to load file", loaderName), err)
	}
	defer fileDir.Close()

	fileInfos, err := fileDir.Readdir(0)
	if err != nil {
		return false, errors.Join(fmt.Errorf("%s: can't read files from dirrectory", loaderName), err)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		if fileInfo.Name() == fileName {
			return true, nil
		}
	}

	return false, nil
}
