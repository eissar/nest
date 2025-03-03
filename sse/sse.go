package sse

import (
	"net/http"
)

type Client struct {
	Conn http.ResponseWriter
	Ch   chan string
}
