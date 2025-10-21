package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"video-service/pb"
	"video-service/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type VideoHandler struct {
	storage        *storage.MinIOStorage
	metadataClient pb.MetadataServiceClient
}

func NewVideoHandler(storage *storage.MinIOStorage, metadataClient pb.MetadataServiceClient) *VideoHandler {
	return &VideoHandler{
		storage:        storage,
		metadataClient: metadataClient,
	}
}

// Загрузка видео
func (h *VideoHandler) UploadVideo(c echo.Context) error {
	file, err := c.FormFile("video")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No video file"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Cannot open file"})
	}
	defer src.Close()

	videoID := uuid.New().String()
	objectName := fmt.Sprintf("videos/%s.mp4", videoID)

	// Сохраняем в MinIO
	if err := h.storage.UploadFile(context.Background(), objectName, src, file.Size); err != nil {
		log.Printf("Upload failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Upload failed"})
	}

	// Создаём запись в Metadata Service
	title := c.FormValue("title")
	ownerID := c.FormValue("owner_id")
	if ownerID == "" {
		ownerID = "anonymous"
	}
	if title == "" {
		title = filepath.Base(file.Filename)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.metadataClient.CreateVideo(ctx, &pb.VideoRequest{
		Title:    title,
		OwnerId:  ownerID,
		MinioKey: objectName,
		Status:   "uploading",
	})
	if err != nil {
		log.Printf("Metadata create failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Metadata error"})
	}

	// Обновляем статус на "ready"
	_, err = h.metadataClient.UpdateVideoStatus(ctx, &pb.UpdateStatusRequest{
		VideoId: videoID,
		Status:  "ready",
	})
	if err != nil {
		log.Printf("Status update failed: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"video_id": videoID,
		"message":  "Upload successful",
	})
}

// Стриминг видео с поддержкой Range
func (h *VideoHandler) StreamVideo(c echo.Context) error {
	videoID := c.Param("id")
	if videoID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing video ID"})
	}

	ctx := context.Background()
	videoResp, err := h.metadataClient.GetVideo(ctx, &pb.VideoIdRequest{VideoId: videoID})
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Video not found"})
	}

	objectName := videoResp.MinioKey
	size, err := h.storage.GetObjectSize(ctx, objectName)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found in storage"})
	}

	file, err := h.storage.GetObject(ctx, objectName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Storage error"})
	}
	defer file.Close()

	// Поддержка Range-запросов
	rangeHeader := c.Request().Header.Get("Range")
	if rangeHeader == "" {
		c.Response().Header().Set("Content-Type", "video/mp4")
		c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", size))
		c.Response().WriteHeader(http.StatusOK)
		io.Copy(c.Response(), file)
		return nil
	}

	// Парсим Range: bytes=0-1023
	var start, end int64 = 0, size - 1
	if strings.HasPrefix(rangeHeader, "bytes=") {
		parts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
		fmt.Sscanf(parts[0], "%d", &start)
		if parts[1] != "" {
			fmt.Sscanf(parts[1], "%d", &end)
		}
	}

	if start >= size {
		return c.JSON(http.StatusRequestedRangeNotSatisfiable, nil)
	}
	if end >= size {
		end = size - 1
	}

	length := end - start + 1

	c.Response().Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
	c.Response().Header().Set("Accept-Ranges", "bytes")
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", length))
	c.Response().Header().Set("Content-Type", "video/mp4")
	c.Response().WriteHeader(http.StatusPartialContent)

	// Пропускаем начало
	io.CopyN(io.Discard, file, start)
	io.CopyN(c.Response(), file, length)

	return nil
}
