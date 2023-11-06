package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	logger_v2 "gitlab.com/milan44/logger-v2"
)

var (
	log = logger_v2.NewColored()
)

func main() {
	pwd, err := os.Getwd()
	must(err)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(log.Middleware())
	r.Use(EnsureIndex(pwd))

	// Handle php files
	InitializePHP(r)

	// Handle static files
	r.Use(HandleBasic(pwd))

	InfoPlain("Host: ", "http://localhost")
	InfoPlain("Root: ", pwd)
	fmt.Println()

	exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost").Start()

	must(r.Run(":80"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
