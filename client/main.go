package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "47.101.72.50:8080", "http service address")
var user string

func main() {
	fmt.Printf("Please enter your full name: ")
	inputReader := bufio.NewReader(os.Stdin)
	user, _ = inputReader.ReadString('\r')
	user = strings.Replace(user, "\r", "", -1)

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	msgChannel := make(chan string)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("%s", message)
		}
	}()

	go func() {
		defer close(done)
		for {
			inputReader := bufio.NewReader(os.Stdin)
			input, err := inputReader.ReadString('\r')
			input = strings.Replace(input, "\r", "", -1)
			if err != nil {
				continue
			}
			msgStr := user + ":" + input
			msgChannel <- msgStr
		}

	}()

	for {
		select {
		case <-done:
			return
		case msgString := <-msgChannel:
			err := c.WriteMessage(websocket.TextMessage, []byte(msgString))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
