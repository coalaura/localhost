package main

import (
	"fmt"
	"io/fs"
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
	dead    *bool
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

	artisan := filepath.Join(pwd, "artisan")

	if _, err := os.Stat(artisan); !os.IsNotExist(err) {
		pwd = filepath.Join(pwd, "public")

		if _, err := os.Stat(pwd); !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to find laravel public directory: %s", pwd)
		}
	} else {
		artisan = ""
	}

	cmd := exec.Command("php", "-S", "localhost:"+port, "-t", pwd)

	var (
		out *os.File
		err error
	)

	if options.Verbose {
		out, err = os.OpenFile("php.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}

		cmd.Stdout = out
		cmd.Stderr = out
	}

	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	var dead bool

	go func() {
		if err := cmd.Wait(); err != nil {
			ErrorF("PHP server exited with error: %s", err)
		}

		if out != nil {
			out.Close()
		}

		dead = true
	}()

	uri, _ := url.Parse("http://localhost:" + port)

	proxy := httputil.NewSingleHostReverseProxy(uri)

	return &PHPServer{
		dead:    &dead,
		process: cmd,
		proxy:   proxy,
		laravel: artisan != "",
	}, nil
}

// Handle forwards the request to the PHP server and returns the response
func (p *PHPServer) Handle(c *gin.Context) {
	if *p.dead {
		c.String(504, "PHP server is dead")

		return
	}

	p.proxy.ServeHTTP(c.Writer, c.Request)
}

func InitializePHP(r *gin.Engine, pwd string) error {
	if !HasPHPSupport() {
		InfoRed("PHP:  ", "unavailable")

		return nil
	}

	hasPHP, err := HasPHPFiles(pwd)
	if err == nil && !hasPHP {
		InfoRed("PHP:  ", "disabled")

		return nil
	}

	php, err = NewPHPServer(pwd)
	if err != nil {
		return err
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

	return nil
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

func HasPHPFiles(dir string) (bool, error) {
	var hasPHP bool

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if hasPHP {
			return filepath.SkipDir
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".php") {
			hasPHP = true

			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return hasPHP, nil
}
