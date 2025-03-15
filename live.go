package main

import (
	_ "embed"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	//go:embed live.html
	liveHtml []byte

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	changeId uint64
	watcher  *fsnotify.Watcher

	connectionId uint64
	connections  sync.Map
)

func InitializeLive(pwd string) error {
	var err error

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(pwd)
	if err != nil {
		return err
	}

	err = filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	debounce := time.NewTimer(250 * time.Millisecond)
	debounce.Stop()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok || !event.Has(fsnotify.Write) {
					return
				}

				debounce.Reset(250 * time.Millisecond)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Warningf("Watcher error: %v\n", err)
			case <-debounce.C:
				atomic.AddUint64(&changeId, 1)

				BroadcastLiveReload()
			}
		}
	}()

	return nil
}

func InjectLive(content []byte) []byte {
	return append(content, liveHtml...)
}

func HandleLive(c *gin.Context) {
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Warningf("Failed to upgrade socket: %v\n", err)

		HandleStatus(c, 400)
	}

	id := atomic.AddUint64(&connectionId, 1)

	connections.Store(id, connection)

	defer func() {
		connections.Delete(id)
		connection.Close()
	}()

	connection.WriteMessage(websocket.TextMessage, []byte(strconv.FormatUint(atomic.LoadUint64(&changeId), 16)))

	for {
		_, _, err := connection.ReadMessage()
		if err != nil {
			break
		}
	}
}

func BroadcastLiveReload() {
	version := []byte(strconv.FormatUint(atomic.LoadUint64(&changeId), 16))

	connections.Range(func(key, value interface{}) bool {
		client, ok := value.(*websocket.Conn)
		if !ok {
			return true
		}

		err := client.WriteMessage(websocket.TextMessage, version)
		if err != nil {
			log.Warningf("Failed to send reload message to client: %v\n", err)

			return true
		}

		return true
	})
}
