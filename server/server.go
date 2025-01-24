package server

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strings"
	"websocket/router"
	"websocket/server/connection"
	"websocket/server/request"
	"websocket/server/response"
)

const (
	ESTABLISH_CONNECTION_PAYLOAD_SIZE = 2048
	SEC_WEB_SOCKET_ACCEPT_STATIC      = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

var ERROR_INVALID_REQUEST = errors.New("Invalid http request")

type ConnectionCallback func(conn *connection.Connection)

type MiddlewareCallback func(conn *connection.Connection) error

type RouteCallback func(request net.Conn)

type Server struct {
	Address   string
	Listener  net.Listener
	listeners []ConnectionCallback
	router    router.GroupRoutes
}

// Comment
func Connect(address string) (Server, error) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return Server{}, err
	}

	return Server{
		Address:  address,
		Listener: listener,
		router:   router.GroupRoutes{},
	}, nil
}

// Comment
func (ctx *Server) Connection(listener ConnectionCallback) {
	ctx.listeners = append(ctx.listeners, listener)
}

// Comment
func (ctx *Server) handshakeReply(req request.Request) error {
	res := response.Create(req.Ws())

	secWebsocketKey := req.Header("sec-websocket-key")

	if secWebsocketKey == "" {
		// Send not allowed response
		return ERROR_INVALID_REQUEST
	}

	alg := sha1.New()

	alg.Write([]byte(strings.Join([]string{secWebsocketKey, SEC_WEB_SOCKET_ACCEPT_STATIC}, "")))

	hashed := base64.StdEncoding.EncodeToString(alg.Sum(nil))

	res.SetProtocol("Switching Protocols")
	res.SetStatus(101)
	res.SetHeader("Upgrade", "websocket")
	res.SetHeader("Connection", "Upgrade")
	res.SetHeader("Sec-WebSocket-Accept", hashed)

	req.Ws().Key = secWebsocketKey // TODO Testing

	return res.Write([]byte(response.HttpBuilder(&res)))
}

// Comment
func (ctx *Server) handshake(conn net.Conn) (*request.Request, *response.Response, error) {
	read := make([]byte, ESTABLISH_CONNECTION_PAYLOAD_SIZE)

	_, err := conn.Read(read)

	if err != nil {
		return nil, nil, err
	}

	req, err := request.Create(conn, string(read))

	if err != nil {
		return nil, nil, err
	}

	err = ctx.handshakeReply(req)

	if err != nil {
		return nil, nil, err
	}

	res := response.Create(req.Ws())

	return &req, &res, nil
}

// Comment
func (ctx *Server) newConnection(conn net.Conn) {
	req, _, err := ctx.handshake(conn)

	if err != nil {
		conn.Close()
		return
	}

	route, err := ctx.router.WsRoute(req.Path())

	if err != nil {
		fmt.Println("Route", req.Path(), "not found")
		return
	}

	go func() {
		for i := 0; i < len(ctx.listeners); i++ {
			go func() {
				ctx.listeners[i](req.Ws())
			}()
		}
	}()

	route.Call(req, req.Ws())

	req.Ws().Emit(connection.EVENT_READY, []byte{})

	req.Ws().Listen()
}

// Comment
func (ctx *Server) Route() *router.Route {
	return &router.Route{
		Routes: &ctx.router,
	}
}

// Comment
func (ctx *Server) Listen() {
	for {
		conn, err := ctx.Listener.Accept()

		if err != nil {
			continue // must log error
		}

		go func() {
			ctx.newConnection(conn)
		}()
	}
}
