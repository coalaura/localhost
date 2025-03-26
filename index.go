package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func EnsureIndex(pwd string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if IsPHPServer() || !strings.HasSuffix(path, "/") {
			c.Next()

			return
		}

		c.Request.URL.Path += "index.html"

		c.Next()
	}
}
