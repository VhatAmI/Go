package main

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"net/http"
	"log"
	"trace"
)

type room struct {
	forward chan *message
	
	join chan *client
	
	leave chan *client
	
	clients map[*client]bool
	
	tracer trace.Tracer
}

func (r *room) run() {

	for {
	select {
	case client := <-r.join:
		r.clients[client] = true
		r.tracer.Trace("New client joined")
	case client := <-r.leave:
		delete(r.clients, client)
		close(client.send)
		r.tracer.Trace("Client left")
	case msg := <-r.forward:
		r.tracer.Trace("Message recieved : ", msg.Message)
		for client := range r.clients {
		select {
		case client.send <- msg:
			//msg sends
			r.tracer.Trace("--sent a message")
		default:
			//failed
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("--message failed to send, removed client")
				}
										}
			}
		}
}

const (
 socketBufferSize = 1024
 messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
socket, err := upgrader.Upgrade(w, req, nil)
if err != nil {
	log.Fatal("ServeHTTP : ", err)
	return
}

authCookie, err := req.Cookie("auth")
if err != nil {
	log.Fatal("Failed to get required authentication cookie: ", err)
	return
}

client := &client{
	socket: socket,
	send: make(chan *message, messageBufferSize),
	room: r,
	userData: objx.MustFromBase64(authCookie.Value),
}
r.join <- client
defer func() { r.leave <-client }()
go client.write()
client.read()
}

func newRoom() *room {
	return &room{
	forward: make(chan *message),
	join:    make(chan *client),
	leave:   make(chan *client),
	clients: make(map[*client]bool),
	}
}