package types

import "fmt"

type Window struct {
	Handle    string `json:"handle"`
	Title     string `json:"title"`
	ProcessId string `json:"processid"`
}

func Hey() {
	fmt.Println("HEY")
}
