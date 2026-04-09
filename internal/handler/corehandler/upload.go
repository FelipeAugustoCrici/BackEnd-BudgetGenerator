package corehandler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"budgetgen/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

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
		context.Background(),
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

	publicURL := os.Getenv("MINIO_PUBLIC_URL")
	url := fmt.Sprintf("%s/%s/%s", publicURL, storage.Bucket, objectName)

	c.JSON(http.StatusOK, gin.H{"url": url})
}
