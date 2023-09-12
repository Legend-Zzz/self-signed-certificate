package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type FileData struct {
	FileName   string
	CreateDate time.Time
}

func getFilesAndDates(folderPath string) ([]FileData, error) {
	fileDataList := []FileData{}

	err := filepath.Walk(folderPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			relativePath, _ := filepath.Rel(folderPath, filePath)

			fileData := FileData{
				FileName:   relativePath,
				CreateDate: fileInfo.ModTime(),
			}
			fileDataList = append(fileDataList, fileData)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return fileDataList, nil
}

func GetFiles(c *gin.Context, outPath string) {
	file, err := getFilesAndDates(outPath)
	if os.IsNotExist(err) {
		c.HTML(http.StatusOK, "files.html", nil)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "files.html", gin.H{
		"FileList": file,
	})
}

func ViewFiles(c *gin.Context, outPath string) {
	filename := c.Param("filename")

	filePath := filepath.Join(outPath, filename)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fileContent = []byte("Unable to view file")
	}

	if len(fileContent) == 0 {
		fileContent = []byte("File content is empty")
	}

	c.String(http.StatusOK, string(fileContent))
}

func DeleteFiles(c *gin.Context, outPath string) {
	filename := c.Param("filename")

	filePath := outPath + filename

	err := os.Remove(filePath)
	if err != nil {
		return
	}
}

func DownloadFiles(c *gin.Context, outPath string) {
	filename := c.Param("filename")

	filePath := outPath + filename

	c.File(filePath)
}

func createFilesIfNotPresent(c *gin.Context, fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to open file"})
		return
	}
	defer file.Close()
}

func ViewResult(c *gin.Context, fileName string) {
	createFilesIfNotPresent(c, fileName)
	Data, err := os.ReadFile(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read file"})
		return
	}
	c.HTML(http.StatusOK, "result.html", gin.H{
		"Content": string(Data),
	})
}

func WriteFiles(c *gin.Context, content string, fileName string) {
	createFilesIfNotPresent(c, fileName)

	oldData, err := os.ReadFile(fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read file"})
		return
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to open file"})
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error writing file newData"})
	}

	_, err = file.Write(oldData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error writing file oldData"})
		return
	}
}
