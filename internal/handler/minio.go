package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUploadURL(c *gin.Context) {
	objectName := c.Query("file_name")

	url, err := h.s.Minio.PresignedUpload(c.Request.Context(), objectName, time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url})
}

func (h *Handler) GetDownloadURL(c *gin.Context) {
	fileName := c.Query("file_name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_name is required"})
		return
	}

	url, err := h.s.Minio.PresignedDownload(c.Request.Context(), fileName, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate download url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"download_url": url,
	})
}
