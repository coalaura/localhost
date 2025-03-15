package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/coalaura/arguments"
	"github.com/coalaura/logger"
	adapter "github.com/coalaura/logger/gin"
	"github.com/gin-gonic/gin"
)

var (
	log     = logger.New()
	options = Options{
		Directory: ".",
	}
)

func main() {
	arguments.Register("cert", 'c', &options.Certificate).WithHelp("Path to ssl certificate")
	arguments.Register("directory", 'd', &options.Directory).WithHelp("Document root")
	arguments.Register("index", 'i', &options.Index).WithHelp("Index file")
	arguments.Register("key", 'k', &options.Key).WithHelp("Path to ssl key")
	arguments.Register("live", 'l', &options.LiveReload).WithHelp("Automatically reload when files change (not with PHP)")
	arguments.Register("open", 'o', &options.Open).WithHelp("Open in default browser once started")
	arguments.Register("port", 'p', &options.Port).WithHelp("Web server port")
	arguments.Register("redirect", 'r', &options.Redirect).WithHelp("Redirect http to https")
	arguments.Register("verbose", 'v', &options.Verbose).WithHelp("Log php server output to php.log")

	arguments.RegisterHelp(true, "Show this help page.")

	arguments.MustParse()

	dir, err := filepath.Abs(options.Directory)
	log.MustPanic(err)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(adapter.GinMiddleware(log))
	r.Use(cors())
	r.Use(EnsureIndex(dir))

	// Handle php files
	InitializePHP(r, dir)

	// Handle static files
	r.Use(HandleBasic(dir))

	// Handle live socket
	if options.LiveReload {
		err = InitializeLive(dir)
		log.MustPanic(err)

		InfoGreen("Live: ", "enabled")
	} else {
		InfoRed("Live: ", "disabled")
	}

	InfoPlain("Host: ", options.GetHost())
	InfoPlain("Root: ", dir)
	InfoPlain("CORS: ", "allow-all")

	if options.HasSSL() {
		InfoGreen("TLS:  ", "enabled")
	} else {
		InfoRed("TLS:  ", "disabled")
	}

	fmt.Println()

	if options.Open {
		exec.Command("rundll32", "url.dll,FileProtocolHandler", options.GetHost()).Start()
	}

	// Redirect http to https
	if options.HasSSL() {
		if options.Redirect {
			go func() {
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					target := "https://" + r.Host + r.URL.RequestURI()

					w.Header().Set("Access-Control-Allow-Origin", "*")

					http.Redirect(w, r, target, http.StatusTemporaryRedirect)
				})

				log.MustPanic(http.ListenAndServe(":80", nil))
			}()
		}

		log.MustPanic(r.RunTLS(fmt.Sprintf(":%d", options.GetPort()), options.Certificate, options.Key))
	} else {
		log.MustPanic(r.Run(fmt.Sprintf(":%d", options.GetPort())))
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Next()
	}
}
