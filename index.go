package main

import (
	"strings"

	"github.com/coalaura/arguments"
	"github.com/gin-gonic/gin"
)

func EnsureIndex(pwd string) gin.HandlerFunc {
	var def string

	if IsPHPServer() {
		def = "index.php"
	} else if !IsLaravelServer() {
		def = "index.html"
	}

	index := arguments.String("i", "index", def)

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if !strings.HasSuffix(path, "/") {
			c.Next()

			return
		}

		if index != "" {
			c.Request.URL.Path += index
		}

		c.Next()
	}
}
