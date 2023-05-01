package limited_file_loader

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
	"time"
)

const (
	loaderName        = "limited file loader"
	defaultSpeedLimit = 1_000 // 1 mb per second
)

func New(speedLimitInKb int64) *FileLoader {
	if speedLimitInKb == 0 {
		speedLimitInKb = defaultSpeedLimit
	}
	return &FileLoader{speedLimit: speedLimitInKb * 1_000}
}

type FileLoader struct {
	speedLimit int64 //represent speed limit in kb per second
}

// Download allow to download file by file link and save it into root file di
func (f *FileLoader) Download(fileLink string, fileRootPath string) error {
	_, err := url.Parse(fileLink)
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't parse file link", loaderName), err)
	}

	fileName := filepath.Base(fileLink)
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

func (f *FileLoader) createAndLoadFile(response *http.Response, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't file to load data", loaderName), err)
	}
	defer file.Close()

	reader := bufio.NewReader(response.Body)
	buffer := make([]byte, f.speedLimit/1000)

	timer := time.NewTimer(time.Millisecond * 2)
	defer timer.Stop()

	for {
		timer.Reset(time.Millisecond * 2)

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

		<-timer.C
	}

	return nil
}

func (f *FileLoader) continueFileLoading(response *http.Response, filePath string) error {
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
	factFileSizeInString := response.Header.Get("Content-Length")
	factFileSize, err := strconv.ParseInt(factFileSizeInString, 10, 64)
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't parse content length from response [response data: %s]", loaderName, factFileSizeInString), err)
	}

	if currentFileSize == factFileSize {
		return nil
	}

	reader := bufio.NewReader(response.Body)
	buffer := make([]byte, f.speedLimit/1000)

	_, err = reader.Discard(int(currentFileSize))
	if err != nil {
		return errors.Join(fmt.Errorf("%s: can't skipe loaded bytes", loaderName), err)
	}

	timer := time.NewTimer(time.Millisecond * 2)
	defer timer.Stop()

	for {
		timer.Reset(time.Millisecond * 2)

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

		<-timer.C
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
