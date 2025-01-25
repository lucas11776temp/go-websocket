package main

import (
	"fmt"
	"log"
	"time"
	"websocket/router"
	"websocket/server"
	"websocket/server/connection"
	"websocket/server/request"
	"websocket/server/response"
)

// Comment
func Auth(req *request.Request, res *response.Response, next router.Next) *response.Response {
	if req.Header("upgrade") == "websocket" {
		fmt.Println("Middleware Websocket")
	}

	return next()
}

func main() {
	machine, err := server.Connect("127.0.0.1:4567")

	if err != nil {
		log.Fatal(err)
	}

	machine.Route().Group("/", func(route *router.Route) {
		route.Ws("chats/{id}", Moving).Middleware(Auth) // TODO Fix: Does not add Middleware...
	})

	fmt.Println("Listening:", machine.Address)

	machine.Listen()
}

// Can move to controller...
func Moving(req *request.Request, ws *connection.Connection) {
	fmt.Println("Connected: ", ws.Address, req.Parameter("id"))

	ws.OnReady(func(data []byte) {
		ws.OnMessage(func(data []byte) {
			fmt.Println("Message: ", string(data))
		})

		go func() {
			for {
				time.Sleep(time.Second * 2)

				if !ws.Alive {
					break
				}

				// ws.Send([]byte(`{"latitude":170.99115,"longitude":26.6284564943}`))
			}
		}()
	})
}
