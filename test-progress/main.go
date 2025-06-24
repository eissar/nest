package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	// "golang.org/x/term"
)

type ProgressMessage struct {
	Mutex sync.Mutex
}

func render() {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	fmt.Println("test")

	// eliminate flickering on terminals that support synchronized output
	fmt.Fprint(w, "\033[?2026h")
	defer fmt.Fprint(w, "\033[?2026l")

	fmt.Fprint(w, "\033[?25l")
	defer fmt.Fprint(w, "\033[?25h")

	fmt.Fprint(w, "testing", time.Now().Local().String(), "\033[K")

	// move the cursor back to the beginning
	fmt.Fprint(w, "\033[A")
	fmt.Fprint(w, "\033[1G")
}

// <https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#cursor-navigation>
func render1(w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	for i := range 20 {
		time.Sleep(200 * time.Millisecond)
		fmt.Fprint(bw, "\u001b[1000D", i+1, "%")
	}
}

func main() {
	ticker := time.NewTicker(100 * time.Millisecond)
	// width, _, _ := term.GetSize(int(os.Stderr.Fd()))
	for range ticker.C {
		n := time.Now().Local().String()
		fmt.Fprint(os.Stdout, "\u001b[1000D", n)
	}
	// render1(os.Stdout)
}
