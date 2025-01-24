package request

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"websocket/server/connection"
)

type Headers map[string]string

type Request struct {
	version float32
	method  string
	path    string
	headers Headers
	body    []byte
	conn    *net.Conn
	ws      connection.Connection
}

type requestInfo struct {
	method  string
	path    string
	version float32
}

var (
	ERROR_INVALID_REQUEST = errors.New("Invalid http request")
)

// Comment
func getRequestInfo(http []string) (requestInfo, error) {
	httpSplit := strings.Split(http[0:1][0], " ")

	if len(httpSplit) != 3 {
		return requestInfo{}, ERROR_INVALID_REQUEST
	}

	httpInfo := strings.Split(httpSplit[2:3][0], "/")

	if strings.ToLower(httpInfo[0]) != "http" || len(httpInfo) != 2 {
		return requestInfo{}, ERROR_INVALID_REQUEST
	}

	version, err := strconv.ParseFloat(httpInfo[1], 1)

	if err != nil {
		return requestInfo{}, ERROR_INVALID_REQUEST
	}

	return requestInfo{
		method:  httpSplit[0:1][0],
		path:    httpSplit[1:2][0],
		version: float32(version),
	}, nil
}

func getHeaders(http []string) Headers {
	headers := make(Headers)

	for i := 0; i < len(http[1:]); i++ {
		header := strings.Split(http[1:][i], ":")

		if len(header) < 2 {
			headers[strings.Trim(strings.ToLower(header[0]), " ")] = ""
			continue
		}

		if len(header) > 2 {
			headers[strings.Trim(strings.ToLower(header[0]), " ")] = strings.Trim(strings.Join(header[1:], ":"), " ")
			continue
		}

		headers[strings.Trim(strings.ToLower(header[0]), " ")] = strings.Trim(header[1], " ")
	}

	return headers
}

// Comment - ParseHttp(http string) (Request, error)
func Create(conn net.Conn, http string) (Request, error) {
	httpArray := strings.Split(http, "\r\n")

	if len(httpArray) < 3 {
		return Request{}, ERROR_INVALID_REQUEST
	}

	info, err := getRequestInfo(httpArray)

	if err != nil {
		return Request{}, err
	}

	// connection := connection.Create(conn)

	return Request{
		version: info.version,
		method:  info.method,
		path:    info.path,
		headers: getHeaders(httpArray),
		conn:    &conn,
		ws:      connection.Create(conn),
	}, nil
}

// Comment
func (ctx *Request) Conn() *net.Conn {
	return ctx.conn
}

// Comment
func (ctx *Request) Ws() *connection.Connection {
	return &ctx.ws
}

// Comment
func (ctx *Request) Version() float32 {
	return ctx.version
}

// Comment
func (ctx *Request) Method() string {
	return ctx.method
}

// Comment
func (ctx *Request) Path() string {
	return ctx.path
}

// Comment
func (ctx *Request) Headers() Headers {
	return ctx.headers
}

// Comment
func (ctx *Request) Header(header string) string {
	header, _ = ctx.headers[header]

	return header
}

// Comment
func (ctx *Request) Body() []byte {
	return ctx.body
}
