package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func EnsureIndex(pwd string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if !strings.HasSuffix(path, "/") {
			c.Next()

			return
		}

		_, ok := GetLocalFile(pwd, "index.html")
		if ok {
			c.Request.URL.Path = "/index.html"

			c.Next()

			return
		}

		_, ok = GetLocalFile(pwd, "index.htm")
		if ok {
			c.Request.URL.Path = "/index.htm"

			c.Next()

			return
		}

		_, ok = GetLocalFile(pwd, "index.php")
		if ok {
			c.Request.URL.Path = "/index.php"

			c.Next()

			return
		}

		HandleStatus(c, 404)
	}
}
