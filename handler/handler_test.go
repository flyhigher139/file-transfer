package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"bytes"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTest a helper function to set up a temporary storage directory and gin router for testing.
func setupTest(t *testing.T) (*gin.Engine, string, func()) {
	tmpDir, err := os.MkdirTemp("", "test-files-*")
	assert.NoError(t, err)

	StoragePath = tmpDir

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	teardown := func() {
		os.RemoveAll(tmpDir)
	}

	return r, tmpDir, teardown
}

// TestListFiles tests the ListFiles handler.
func TestListFiles(t *testing.T) {
	r, tmpDir, teardown := setupTest(t)
	defer teardown()

	// Create some dummy files
	_, err := os.Create(filepath.Join(tmpDir, "file1.txt"))
	assert.NoError(t, err)
	_, err = os.Create(filepath.Join(tmpDir, "file2.txt"))
	assert.NoError(t, err)

	r.GET("/api/files", ListFiles)

	req, _ := http.NewRequest(http.MethodGet, "/api/files", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response["files"], "file1.txt")
	assert.Contains(t, response["files"], "file2.txt")
}

// TestSimpleUploadFile tests the SimpleUploadFile handler.
func TestSimpleUploadFile(t *testing.T) {
	r, tmpDir, teardown := setupTest(t)
	defer teardown()

	r.POST("/api/simple/upload", SimpleUploadFile)

	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	file, err := mw.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)
	file.Write([]byte("hello world"))
	mw.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/simple/upload", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the file was created
	_, err = os.Stat(filepath.Join(tmpDir, "test.txt"))
	assert.NoError(t, err)
}

// TestSimpleDownloadFile tests the SimpleDownloadFile handler.
func TestSimpleDownloadFile(t *testing.T) {
	r, tmpDir, teardown := setupTest(t)
	defer teardown()

	// Create a dummy file
	fileContent := "hello download"
	filePath := filepath.Join(tmpDir, "download.txt")
	os.WriteFile(filePath, []byte(fileContent), 0644)

	r.GET("/api/simple/download/:filename", SimpleDownloadFile)

	req, _ := http.NewRequest(http.MethodGet, "/api/simple/download/download.txt", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "attachment; filename=download.txt", w.Header().Get("Content-Disposition"))
	assert.Equal(t, fileContent, w.Body.String())
}

// TestChunkedUploadAndMerge tests the chunked upload and merge functionality.
func TestChunkedUploadAndMerge(t *testing.T) {
	r, tmpDir, teardown := setupTest(t)
	defer teardown()

	r.POST("/api/files", UploadFile)
	r.POST("/api/files/merge", MergeFile)

	fileIdentifier := "test-identifier"
	chunkDir := filepath.Join(tmpDir, "tmp", fileIdentifier)
	os.MkdirAll(chunkDir, 0755)

	// Simulate chunk uploads
	chunk1Content := "hello "
	os.WriteFile(filepath.Join(chunkDir, "1"), []byte(chunk1Content), 0644)
	chunk2Content := "world"
	os.WriteFile(filepath.Join(chunkDir, "2"), []byte(chunk2Content), 0644)

	// Simulate merge request
	mergeBody := new(bytes.Buffer)
	mw := multipart.NewWriter(mergeBody)
	mw.WriteField("filename", "merged.txt")
	mw.WriteField("fileIdentifier", fileIdentifier)
	mw.WriteField("totalChunks", "2")
	mw.Close()

	mergeReq, _ := http.NewRequest(http.MethodPost, "/api/files/merge", mergeBody)
	mergeReq.Header.Set("Content-Type", mw.FormDataContentType())

	mergeW := httptest.NewRecorder()
	r.ServeHTTP(mergeW, mergeReq)

	assert.Equal(t, http.StatusOK, mergeW.Code)

	// Check if the merged file was created and has the correct content
	mergedFilePath := filepath.Join(tmpDir, "merged.txt")
	mergedContent, err := os.ReadFile(mergedFilePath)
	assert.NoError(t, err)
	assert.Equal(t, chunk1Content+chunk2Content, string(mergedContent))

	// Check if the temp chunk directory was removed
	_, err = os.Stat(chunkDir)
	assert.True(t, os.IsNotExist(err))
}

// TestDownloadFileWithRange tests the DownloadFile handler with a Range request.
func TestDownloadFileWithRange(t *testing.T) {
	r, tmpDir, teardown := setupTest(t)
	defer teardown()

	// Create a dummy file
	fileContent := "hello range request"
	filePath := filepath.Join(tmpDir, "range.txt")
	os.WriteFile(filePath, []byte(fileContent), 0644)

	r.GET("/api/files/:filename", DownloadFile)

	req, _ := http.NewRequest(http.MethodGet, "/api/files/range.txt", nil)
	req.Header.Set("Range", "bytes=6-10")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPartialContent, w.Code)
	assert.Equal(t, "range", w.Body.String())
}