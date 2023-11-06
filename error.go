package main

import "github.com/gin-gonic/gin"

func HandleStatus(c *gin.Context, status int) {
	switch status {
	case 404:
		c.String(404, "404 Not Found")
		c.Abort()

		return
	case 500:
		c.String(500, "500 Internal Server Error")
		c.Abort()

		return
	}

	c.String(status, "Unknown Error")
}
