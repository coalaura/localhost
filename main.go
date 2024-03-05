package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/coalaura/logger"
	"github.com/gin-gonic/gin"
)

var (
	log = logger.New()
)

func main() {
	pwd, err := os.Getwd()
	must(err)

	port := parsePort()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(log.Middleware())
	r.Use(EnsureIndex(pwd))

	// Handle php files
	InitializePHP(r)

	// Handle static files
	r.Use(HandleBasic(pwd))

	InfoPlain("Host: ", "http://localhost"+port)
	InfoPlain("Root: ", pwd)
	fmt.Println()

	exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost"+port).Start()

	must(r.Run(port))
}

func parsePort() string {
	if len(os.Args) > 1 {
		num, err := strconv.ParseInt(os.Args[1], 10, 64)
		must(err)

		return fmt.Sprintf(":%d", num)
	}

	return ":80"
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
