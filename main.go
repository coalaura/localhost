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
	arguments.Register("cert", 'c', &options.Certificate).WithHelp("Path to the SSL certificate file for HTTPS (enables SSL)")
	arguments.Register("directory", 'd', &options.Directory).WithHelp("Document root directory for serving files")
	arguments.Register("index", 'i', &options.Index).WithHelp("Name of the index file to serve (e.g., index.html)")
	arguments.Register("key", 'k', &options.Key).WithHelp("Path to the SSL key file for HTTPS (required with --cert)")
	arguments.Register("live", 'l', &options.LiveReload).WithHelp("Enable live reload on file changes (refreshes browser)")
	arguments.Register("open", 'o', &options.Open).WithHelp("Automatically open the web server URL in the default browser")
	arguments.Register("port", 'p', &options.Port).WithHelp("Port number for the web server to listen on (default: 80)")
	arguments.Register("redirect", 'r', &options.Redirect).WithHelp("Redirect HTTP traffic to HTTPS (requires --cert and --key)")
	arguments.Register("verbose", 'v', &options.Verbose).WithHelp("Enable verbose logging of PHP server output to php.log")

	arguments.RegisterHelp(true, "Display this help message and exit")

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
