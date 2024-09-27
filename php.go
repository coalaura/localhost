package main

import (
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type PHPServer struct {
	laravel bool
	process *exec.Cmd
	proxy   *httputil.ReverseProxy
}

var (
	php *PHPServer
)

// HasPHPSupport checks if the current system has a "php" executable
func HasPHPSupport() bool {
	php, err := exec.LookPath("php")

	return err == nil && php != ""
}

func IsLaravelServer() bool {
	return php != nil && php.laravel
}

func IsPHPServer() bool {
	return php != nil && !php.laravel
}

func NewPHPServer(pwd string) (*PHPServer, error) {
	var cmd *exec.Cmd

	artisan := filepath.Join(pwd, "artisan")

	if _, err := os.Stat(artisan); os.IsNotExist(err) {
		cmd = exec.Command("php", "-S", "localhost:8989", "-t", pwd)
	} else {
		cmd = exec.Command("php", artisan, "serve", "--port=8989")
	}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	uri, _ := url.Parse("http://localhost:8989")

	proxy := httputil.NewSingleHostReverseProxy(uri)

	return &PHPServer{
		process: cmd,
		proxy:   proxy,
		laravel: artisan != "",
	}, nil
}

// Handle forwards the request to the PHP server and returns the response
func (p *PHPServer) Handle(c *gin.Context) {
	p.proxy.ServeHTTP(c.Writer, c.Request)
}

func InitializePHP(r *gin.Engine) {
	if !HasPHPSupport() {
		InfoRed("PHP:  ", "disabled")

		return
	}

	r.Use(func(c *gin.Context) {
		php.Handle(c)

		c.Abort()
	})

	if php.laravel {
		InfoGreen("PHP:  ", "enabled (laravel)")
	} else {
		InfoGreen("PHP:  ", "enabled (plain)")
	}
}
