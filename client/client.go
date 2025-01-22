package client

type WebSocket struct {
	Host string
}

type EventCallback func(data []byte)

// Comment
func Connect(url string) WebSocket {

	return WebSocket{}
}

// Comment
func (ctx *WebSocket) OnMessage(event EventCallback) {

}

// Comment
func (ctx *WebSocket) OnError(event EventCallback) {

}

// Comment
func (ctx *WebSocket) OnClose(event EventCallback) {

}

// Comment
func (ctx *WebSocket) OnOpen(event EventCallback) {

}

// Comment
func (ctx *WebSocket) Listen() {

}
