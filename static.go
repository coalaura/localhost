package main

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HandleBasic(pwd string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		file, ok := GetLocalFile(pwd, path)
		if !ok {
			HandleStatus(c, 404)

			return
		}

		content, status := GetFileData(file)
		if status != 200 {
			HandleStatus(c, status)

			return
		}

		c.Data(200, GetMimeType(c), content)
	}
}

func GetLocalFile(pwd, path string) (string, bool) {
	file := filepath.Join(pwd, path)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", false
	}

	return file, true
}

func GetFileData(path string) ([]byte, int) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, 500
	}

	return content, 200
}

func GetMimeType(c *gin.Context) string {
	ext := filepath.Ext(c.Request.URL.Path)

	switch ext {
	case ".css":
		return "text/css"
	case ".js":
		return "text/javascript"
	case ".html", ".htm":
		return "text/html"
	case ".png":
		return "image/png"
	case ".jpg", "jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".pdf":
		return "application/pdf"
	case ".ttf":
		return "font/ttf"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".eot":
		return "application/vnd.ms-fontobject"

	// Special cases
	case ".php":
		return "php"
	}

	return "text/plain"
}
