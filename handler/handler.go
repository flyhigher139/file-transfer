package handler

import (
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// StoragePath is the directory where files are stored.
var StoragePath string

// ListFiles lists all uploaded files.
func ListFiles(c *gin.Context) {
	files, err := os.ReadDir(StoragePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read files directory"})
		return
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	c.JSON(http.StatusOK, gin.H{"files": filenames})
}

// UploadFile handles chunked file uploads.
func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("get form err: %s", err.Error())})
		return
	}
	defer file.Close()

	chunkNumber := c.PostForm("chunkNumber")
	fileIdentifier := c.PostForm("fileIdentifier")

	chunkDir := filepath.Join(StoragePath, "tmp", fileIdentifier)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temp directory"})
		return
	}

	dst, err := os.Create(filepath.Join(chunkDir, chunkNumber))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save chunk"})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not write chunk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Chunk %s for %s uploaded successfully", chunkNumber, header.Filename)})
}

// MergeFile merges the uploaded chunks into a single file.
func MergeFile(c *gin.Context) {
	filename := c.PostForm("filename")
	fileIdentifier := c.PostForm("fileIdentifier")
	

	chunkDir := filepath.Join(StoragePath, "tmp", fileIdentifier)
	defer os.RemoveAll(chunkDir)

	finalFilePath := filepath.Join(StoragePath, filename)
	finalFile, err := os.Create(finalFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create final file"})
		return
	}
	defer finalFile.Close()

	totalChunksStr := c.PostForm("totalChunks")
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid totalChunks"})
		return
	}

	for i := 1; i <= totalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not open chunk %d", i)})
			return
		}

		_, err = io.Copy(finalFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not merge chunk %d", i)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "File merged successfully", "filename": filename})
}

// DownloadFile handles file downloads with support for Range requests.
// SimpleUploadFile handles a single file upload.
func SimpleUploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("get form err: %s", err.Error())})
		return
	}
	defer file.Close()

	filename := header.Filename
	dst, err := os.Create(filepath.Join(StoragePath, filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create file"})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s uploaded successfully", filename)})
}

// SimpleDownloadFile handles a simple file download.
func SimpleDownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(StoragePath, filename)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

// DownloadFile handles file downloads with support for Range requests.
func DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(StoragePath, filename)

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	// http.ServeContent will handle Range requests
	http.ServeContent(c.Writer, c.Request, filename, fileInfo.ModTime(), file)

	_ = fileSize // This is to avoid unused variable error, as ServeContent uses it internally.
}