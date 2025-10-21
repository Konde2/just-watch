package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Получаем адреса сервисов из переменных окружения
	videoAddr := os.Getenv("VIDEO_SERVICE_ADDR")
	if videoAddr == "" {
		videoAddr = "http://video-service:8081"
	}

	// Проксируем запросы
	e.POST("/upload", proxyTo(videoAddr+"/upload"))
	e.GET("/stream/:id", proxyTo(videoAddr+"/stream/:id"))

	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("API Gateway listening on :" + port)
	log.Fatal(e.Start(":" + port))
}

// Вспомогательная функция для проксирования
func proxyTo(target string) echo.HandlerFunc {
	u, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(u)
	return func(c echo.Context) error {
		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
