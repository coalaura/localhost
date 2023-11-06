package main

import (
	"net/http/httputil"
	"net/url"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type PHPServer struct {
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

func NewPHPServer() (*PHPServer, error) {
	cmd := exec.Command("php", "-S", "localhost:8989")

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	uri, _ := url.Parse("http://localhost:8989") // Change this URL to your destination

	proxy := httputil.NewSingleHostReverseProxy(uri)

	return &PHPServer{
		process: cmd,
		proxy:   proxy,
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

	s, err := NewPHPServer()
	must(err)

	php = s

	r.Use(func(c *gin.Context) {
		mime := GetMimeType(c)

		if mime == "php" {
			php.Handle(c)
			c.Abort()

			return
		}

		c.Next()
	})

	InfoGreen("PHP:  ", "enabled")
}
