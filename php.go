package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	port := FindFreePort()

	var cmd *exec.Cmd

	artisan := filepath.Join(pwd, "artisan")

	if _, err := os.Stat(artisan); os.IsNotExist(err) {
		cmd = exec.Command("php", "-S", "localhost:"+port, "-t", pwd)
	} else {
		cmd = exec.Command("php", artisan, "serve", "--port="+port)
	}

	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	reader, writer := io.Pipe()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(writer)

	go io.Copy(multiWriter, stdout)
	go io.Copy(multiWriter, stderr)

	go func() {
		defer writer.Close()
		defer reader.Close()

		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			text := strings.TrimRight(scanner.Text(), "\r\n")

			PHPOut(text)
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	uri, _ := url.Parse("http://localhost:" + port)

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

func FindFreePort() string {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	return fmt.Sprint(port)
}
