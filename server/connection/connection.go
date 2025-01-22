package connection

import (
	"net"
	"websocket/server/frame"
)

const (
	// PAYLOAD_SIZE = 80
	PAYLOAD_SIZE = 2048
)

type Event string

const (
	EVENT_READY   Event = "ready"
	EVENT_MESSAGE Event = "message"
	EVENT_PING    Event = "ping"
	EVENT_PONG    Event = "pong"
	EVENT_ERROR   Event = "error"
	EVENT_CLOSE   Event = "close"
)

type EventCallback func(data []byte)

type EventCallbacks map[Event][]EventCallback

type Connection struct {
	Address   string
	Alive     bool
	Listener  net.Conn
	listeners EventCallbacks
	Key       string
}

// Comment
func Create(conn net.Conn) Connection {
	return Connection{
		Address:   conn.RemoteAddr().String(),
		Alive:     true,
		Listener:  conn,
		listeners: make(EventCallbacks),
	}
}

// Comment
func (ctx *Connection) Emit(event Event, data []byte) {
	for e := range ctx.listeners[event] {
		go func() {
			ctx.listeners[event][e](data)
		}()
	}
}

// Comment
func (ctx *Connection) OnReady(callback EventCallback) {
	ctx.listeners[EVENT_READY] = append(ctx.listeners[EVENT_READY], callback)
}

// Comment
func (ctx *Connection) OnMessage(callback EventCallback) {
	ctx.listeners[EVENT_MESSAGE] = append(ctx.listeners[EVENT_MESSAGE], callback)
}

// Comment
func (ctx *Connection) OnError(callback EventCallback) {
	ctx.listeners[EVENT_ERROR] = append(ctx.listeners[EVENT_ERROR], callback)
}

// Comment
func (ctx *Connection) OnPing(callback EventCallback) {
	ctx.listeners[EVENT_PING] = append(ctx.listeners[EVENT_PING], callback)
}

// Comment
func (ctx *Connection) OnClose(callback EventCallback) {
	ctx.listeners[EVENT_CLOSE] = append(ctx.listeners[EVENT_CLOSE], callback)
}

// Comment
func (ctx *Connection) Send(data []byte) error {
	frame := frame.Encode(data)

	_, err := ctx.Listener.Write(frame.Payload())

	return err
}

// Comment
func (ctx *Connection) Ping(data []byte) error {

	return nil
}

// Comment
func (ctx *Connection) Pong(data []byte) error {

	return nil
}

// Comment
func (ctx *Connection) Close() error {
	return ctx.Listener.Close()
}

// Comment
func (ctx *Connection) Listen() {
	for {
		payload := make([]byte, PAYLOAD_SIZE)

		_, err := ctx.Listener.Read(payload)

		if err != nil {
			ctx.Alive = false

			continue
		}

		go func() {
			frame := frame.Decode(payload)

			if frame.IsClose() {
				ctx.Close()
				return
			}

			if frame.IsPing() {

			}

			for i := 0; i < len(ctx.listeners[EVENT_MESSAGE]); i++ {
				go func() {
					ctx.listeners[EVENT_MESSAGE][i](frame.Data())
				}()
			}
		}()
	}
}

// func b([]byte)  {

// }
