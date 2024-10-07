package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/coalaura/arguments"
	"github.com/coalaura/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log = logger.New()
)

func main() {
	help()

	var (
		err  error
		port int
		host string
	)

	dir := arguments.String("d", "directory", "")
	cert := arguments.String("c", "cert", "")
	key := arguments.String("k", "key", "")

	if cert != "" && key != "" {
		port = arguments.IntN("p", "port", 443)

		host = "https://localhost"

		if port != 443 {
			host += fmt.Sprintf(":%d", port)
		}
	} else {
		port = arguments.IntN("p", "port", 80)

		host = "http://localhost"

		if port != 80 {
			host += fmt.Sprintf(":%d", port)
		}
	}

	if dir == "" {
		dir = "."
	}

	dir, err = filepath.Abs(dir)
	must(err)

	php, err = NewPHPServer(dir)
	must(err)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(log.Middleware())
	r.Use(cors.Default())
	r.Use(EnsureIndex(dir))

	// Handle php files
	InitializePHP(r)

	// Handle static files
	r.Use(HandleBasic(dir))

	InfoPlain("Host: ", host)
	InfoPlain("Root: ", dir)

	if cert != "" && key != "" {
		InfoPlain("TLS:  ", "enabled")
	} else {
		InfoPlain("TLS:  ", "disabled")
	}

	fmt.Println()

	if arguments.Bool("o", "open", false) {
		exec.Command("rundll32", "url.dll,FileProtocolHandler", host).Start()
	}

	if cert != "" && key != "" {
		// Redirect http to https
		if arguments.Bool("r", "redirect", false) {
			go func() {
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					target := "https://" + r.Host + r.URL.RequestURI()

					http.Redirect(w, r, target, http.StatusTemporaryRedirect)
				})

				must(http.ListenAndServe(":80", nil))
			}()
		}

		must(r.RunTLS(fmt.Sprintf(":%d", port), cert, key))
	} else {
		must(r.Run(fmt.Sprintf(":%d", port)))
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
