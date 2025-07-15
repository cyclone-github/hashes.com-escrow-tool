package main

import (
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
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
	// open log file
	file, err := os.OpenFile("websocket.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// send all logs to stderr (and to file)
	mw := io.MultiWriter(os.Stderr, file)
	log.SetOutput(mw)

	// trap SIGINT/SIGTERM to exit on Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigs)

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
			// raw json only
			os.Stdout.Write(append(message, '\n'))
		}
	}

	// reconnect loop (exit on signal)
	for {
		done := make(chan struct{})
		go func() {
			start()
			close(done)
		}()

		select {
		case <-sigs:
			log.Println("Interrupt received, stopping monitor")
			return nil
		case <-done:
			log.Println("Disconnected! Reconnecting.")
			time.Sleep(reconnectInterval)
		}
	}
}
