package corehandler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"budgetgen/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type presignResponse struct {
	UploadURL string `json:"uploadUrl"`
	PublicURL string `json:"publicUrl"`
	ObjectKey string `json:"objectKey"`
}

func PresignUpload(c *gin.Context) {
	if storage.Client == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage não disponível"})
		return
	}

	var body struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"contentType"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ext := filepath.Ext(body.Filename)
	objectKey := fmt.Sprintf("%s%s", uuid.NewString(), ext)

	contentType := body.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Generate presigned PUT URL — valid for 10 minutes
	presignedURL, err := storage.Client.PresignedPutObject(
		c.Request.Context(),
		storage.Bucket,
		objectKey,
		10*time.Minute,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao gerar URL de upload"})
		return
	}

	publicURL := os.Getenv("S3_PUBLIC_URL")
	var fileURL string
	if publicURL != "" {
		fileURL = fmt.Sprintf("%s/%s", publicURL, objectKey)
	} else {
		endpoint := os.Getenv("S3_ENDPOINT")
		fileURL = fmt.Sprintf("%s/%s/%s", endpoint, storage.Bucket, objectKey)
	}

	c.JSON(http.StatusOK, presignResponse{
		UploadURL: presignedURL.String(),
		PublicURL: fileURL,
		ObjectKey: objectKey,
	})
}

// Keep direct upload as fallback
func Upload(c *gin.Context) {
	if storage.Client == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "storage não disponível"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo não encontrado"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	objectName := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = storage.Client.PutObject(
		c.Request.Context(),
		storage.Bucket,
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar arquivo"})
		return
	}

	publicURL := os.Getenv("S3_PUBLIC_URL")
	var url string
	if publicURL != "" {
		url = fmt.Sprintf("%s/%s", publicURL, objectName)
	} else {
		endpoint := os.Getenv("S3_ENDPOINT")
		url = fmt.Sprintf("%s/%s/%s", endpoint, storage.Bucket, objectName)
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
