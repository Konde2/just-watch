package main

import (
	"log"
	"os"

	"video-service/handlers"
	pb "video-service/pb"
	"video-service/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// MinIO
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	minioStorage, err := storage.NewMinIOStorage(endpoint, accessKey, secretKey, "videos")
	if err != nil {
		log.Fatal("MinIO init failed:", err)
	}

	// gRPC клиент к Metadata Service
	metadataAddr := os.Getenv("METADATA_SERVICE_ADDR")
	if metadataAddr == "" {
		metadataAddr = "localhost:50051"
	}
	conn, err := grpc.Dial(metadataAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("gRPC connect failed:", err)
	}
	defer conn.Close()
	metadataClient := pb.NewMetadataServiceClient(conn)

	// Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := handlers.NewVideoHandler(minioStorage, metadataClient)

	e.POST("/upload", handler.UploadVideo)
	e.GET("/stream/:id", handler.StreamVideo)

	// Метрики Prometheus
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	port := os.Getenv("VIDEO_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("Video Service listening on :" + port)
	log.Fatal(e.Start(":" + port))
}
