package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/HarvestStars/YuePaoTalent/protocol"
	"github.com/gorilla/websocket"
)

var user string
var ServerHost string

func main() {
	for k, v := range os.Args {
		if k == 0 {
			continue
		}
		ServerHost = v
		break
	}
	//var addr = flag.String("addr", ServerHost, "http service address")
	// 登录 or 注册
	var req protocol.UserReq
	req.UserName, req.PassWord = typeInScreen()
	user = req.UserName
	b, _ := json.Marshal(&req)
	resp, err := http.Post("http://"+ServerHost+"/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	var res protocol.UserResp
	json.Unmarshal(bytes, &res)
	fmt.Print(res.Data + "\n")
	if res.Code == 400 {
		return
	}

	// 开始聊天
	flag.Parse()
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: ServerHost, Path: "/chat"}
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

func typeInScreen() (string, string) {
	fmt.Printf("Please enter your USER NAME:")
	inputReaderUser := bufio.NewReader(os.Stdin)
	user, _ = inputReaderUser.ReadString('\r')
	user = strings.Replace(user, "\r", "", -1)

	fmt.Printf("Please enter your PASS WORD:")
	inputReaderPWD := bufio.NewReader(os.Stdin)
	pwd, _ := inputReaderPWD.ReadString('\r')
	pwd = strings.Replace(pwd, "\r", "", -1)
	return user, pwd
}
