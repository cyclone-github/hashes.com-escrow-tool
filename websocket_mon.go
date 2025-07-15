package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	address           = "hashes.com:443"        // hashes.com URL
	path              = "/en/api/jobs_wss/"     // websocket api
	scheme            = "wss"                   // scheme
	pingInterval      = 60 * time.Second        // ping every nth seconds
	reconnectInterval = 1000 * time.Millisecond // reconnect interval
)

// connect to hashes.com websocket using API key
func connect(apiKey string) (*websocket.Conn, error) {
	u := url.URL{Scheme: scheme, Host: address, Path: path}
	q := u.Query()
	q.Add("key", apiKey)
	u.RawQuery = q.Encode()

	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
		return nil, err
	}
	return c, nil
}

// start websocket feed
func monitorWebsocket(apiKey string) error {
	start := func() {
		conn, err := connect(apiKey)
		if err != nil {
			log.Println("Error connecting:", err)
			return
		}
		defer conn.Close()

		pingTicker := time.NewTicker(pingInterval)
		defer pingTicker.Stop()

		// send regular pings
		go func() {
			for range pingTicker.C {
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("ping error:", err)
					return
				}
			}
		}()

		// read loop: raw json to stdout and errors to stderr via log
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			os.Stdout.Write(append(message, '\n'))
		}
	}

	// reconnect loop
	for {
		start()
		log.Println("Disconnected! Reconnecting.")
		time.Sleep(reconnectInterval)
	}
}
