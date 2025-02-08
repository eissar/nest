package types

import (
	_ "fmt"
	"net/http"
)

type Client struct {
	Conn http.ResponseWriter
	Ch   chan string
}

type Window struct {
	Handle    string `json:"handle"`
	Title     string `json:"title"`
	ProcessId string `json:"processid"`
}

//func Hey() {
//	fmt.Println("HEY")
//}
