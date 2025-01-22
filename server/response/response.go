package response

import (
	"strconv"
	"strings"
	"websocket/server/connection"
)

type Headers map[string]string

type Response struct {
	status   int16
	protocol string
	headers  Headers
	body     []byte
	*connection.Connection
}

// Comment
func Create(connection *connection.Connection) Response {
	return Response{
		protocol:   "",
		status:     200,
		headers:    make(Headers),
		Connection: connection,
		body:       []byte{},
	}
}

// Comment
func (ctx *Response) SetStatus(code int16) *Response {
	ctx.status = code

	return ctx
}

// Comment
func (ctx *Response) SetProtocol(protocol string) *Response {
	ctx.protocol = protocol

	return ctx
}

// Comment
func (ctx *Response) SetHeader(header string, value string) *Response {
	ctx.headers[header] = value

	return ctx
}

// Comment
func (ctx *Response) SetBody(data []byte) *Response {
	ctx.body = data

	return ctx
}

// Comment
func (ctx *Response) Write(data []byte) error {
	_, err := ctx.Listener.Write(data)

	return err
}

// Comment
func HttpBuilder(response *Response) string {
	builder := []string{
		strings.Trim(
			strings.Join([]string{"HTTP/1.1", strconv.Itoa(int(response.status)), response.protocol}, " "), " ",
		),
	}

	for header := range response.headers {
		builder = append(
			builder,
			strings.Join([]string{strings.Trim(header, " "), strings.Trim(response.headers[header], " ")}, ": "),
		)
	}

	return strings.Join(builder, "\r\n") + "\r\n\r\n" // TODO Must Add data field
}

// Comment
func (ctx *Response) Send(data []byte) error {
	err := ctx.Write([]byte(HttpBuilder(ctx)))

	if err != nil {
		return err
	}

	return nil
}
