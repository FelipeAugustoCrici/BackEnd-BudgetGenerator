package corehandler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ImageProxy(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "public, max-age=86400")
	c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, resp.Body, nil)
}
