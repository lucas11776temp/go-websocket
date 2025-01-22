package main

import (
	"fmt"
	"log"
	"time"
	"websocket/server"
	"websocket/server/connection"
)

// chats/12
// http://localhost/chats/12

/*

route := websocket.Route{}


route.WsGroup('chats', func (route route.Ws) {
	route.Ws('/', func(req server.Request, socket server.Connection) {
		socket.OnReady(func (data []byte) {
			conn.OnMessage(func(data []byte) {
				fmt.Println(string(data))
			})
		})
	})
})

route.Ws('/', func(req server.Request, socket server.Connection) {
	socket.OnReady(func (data []byte) {
		conn.OnMessage(func(data []byte) {
			fmt.Println(string(data))
		})
	})
})

*/

func main() {
	machine, err := server.Connect("127.0.0.1:4567")

	// fmt.Println("0x03: ", int(0x03))
	// fmt.Println("0x05: ", int(0x05))
	// fmt.Println("0x09: ", byte(0x09))
	// fmt.Println("0x7F: ", byte(0x7F))

	// fmt.Println("127 ^ b", 127>>0x02)

	if err != nil {
		log.Fatal(err)
	}

	machine.Connection(func(conn *connection.Connection) {
		fmt.Println("New Connection:", conn.Address)

		conn.OnReady(func(data []byte) {
			count := 0

			conn.OnMessage(func(data []byte) {
				count++
				// fmt.Println("Message: ", count)
			})

			go func() {
				for {
					time.Sleep(time.Second * 2)

					if !conn.Alive {
						break
					}

					// conn.Send([]byte("Hello World"))

					str := `[{"latitude":170.99115,"longitude":26.6284564943},`
					str += `{"latitude":170.9911527235034,"longitude":26.6242529},`
					str += `{"latitude":170.99111111}]`

					conn.Send([]byte(str))

				}
			}()

		})
	})

	fmt.Println("Listening:", machine.Address)

	machine.Listen()
}
