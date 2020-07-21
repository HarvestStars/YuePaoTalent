package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8080", "http service address")
var hub *server.Hub

func main() {
	flag.Parse()
	hub = newHub()
	go hub.run()
	r := gin.Default()
	r.GET("/ws", WebSocket)
	r.Run(":8080")
	// r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })
	// err := http.ListenAndServe(*addr, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}

func WebSocket(c *gin.Context) {
	serveWs(hub, c.Writer, c.Request)
}
