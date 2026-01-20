package miniohandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
)

type MinioHandler struct {
	minio minioservice.Minio
}

func NewMinioHandler(minio minioservice.Minio) *MinioHandler {
	return &MinioHandler{minio: minio}
}

func (h *MinioHandler) GetUploadURL(c *gin.Context) {
	objectName := c.Query("file_name")

	url, err := h.minio.PresignedUpload(c.Request.Context(), objectName, time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"upload_url": url})
}

func (h *MinioHandler) GetDownloadURL(c *gin.Context) {
	fileName := c.Query("file_name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_name is required"})
		return
	}

	url, err := h.minio.PresignedDownload(c.Request.Context(), fileName, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate download url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"download_url": url,
	})
}
