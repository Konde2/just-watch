package handlers

import (
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
)

func ProxyTo(target string) echo.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return func(c echo.Context) error {
		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
