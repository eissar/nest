package progress

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

type ProgressMessage struct {
	Mutex sync.Mutex
}

type PrinterFunc func(msg any)

func Render0() {
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

func Render1(c context.Context, pollingCh chan bool, rFunc PrinterFunc) {
	t := time.NewTicker(300 * time.Millisecond)
	select {
	case <-t.C: // tick
		rFunc("asdf" + time.Now().Local().String())
	case <-c.Done():
		return
	}
}

// <https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html#cursor-navigation>
func Testrender1(w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	for i := range 20 {
		time.Sleep(1000 * time.Millisecond)
		fmt.Fprint(bw, "\u001b[1000D", i+1, "%")
	}
}

func TestRender2() {
	ticker := time.NewTicker(100 * time.Millisecond)
	// width, _, _ := term.GetSize(int(os.Stderr.Fd()))
	for range ticker.C {
		n := time.Now().Local().String()
		fmt.Fprint(os.Stdout, "\u001b[1000D", n)
	}
	// render1(os.Stdout)
}

// getTermWidth returns the width of the terminal.
func getTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback to a default width if there's an error
		return 80
	}
	return width
}

// TODO: how does bufio work?
// TODO: this is pretty bad.j

// func Render() {
// 	ticker := time.NewTicker(200 * time.Millisecond)
// 	defer ticker.Stop()
//
// 	w := os.Stdout
// 	linesRendered := 0
//
// 	// Hide cursor and ensure it's shown on exit
// 	fmt.Fprint(w, "\033[?25l")
// 	defer fmt.Fprint(w, "\033[?25h")
//
// 	for i := 0; i < 50; i++ {
// 		<-ticker.C
//
// 		// -- PREPARE THE NEW CONTENT --
// 		width := getTermWidth()
// 		// Building a long string to demonstrate wrapping
// 		msg := fmt.Sprintf("Update #%d: %s", i, strings.Repeat(time.Now().Format("15:04:05.000 "), 5))
//
// 		// -- RENDER LOGIC --
//
// 		// 1. Move cursor up to the start of the previous render if necessary
// 		if linesRendered > 1 {
// 			fmt.Fprintf(w, "\033[%dA", linesRendered)
// 		}
// 		// 2. Move cursor to the beginning of the line
// 		fmt.Fprint(w, "\r")
//
// 		// 3. Clear from cursor to the end of the screen
// 		fmt.Fprint(w, "\033[J")
//
// 		// 4. Print the new message
// 		fmt.Fprint(w, msg)
//
// 		// 5. Calculate how many lines the new message occupies
// 		linesRendered = (len(msg) + width - 1) / width
// 		if linesRendered == 0 {
// 			linesRendered = 1
// 		}
// 	}
// 	// Print a final newline to move past the animation
// 	fmt.Println()
// }

// TODO: clean up this garbage
//
// Render processes messages from a channel and displays them in the terminal,
// overwriting the previous output. It stops when the context is canceled
// or the message channel is closed.
func Render(ctx context.Context, wg *sync.WaitGroup, messageChannel <-chan string) {
	wg.Add(1)
	defer wg.Done()

	w := os.Stdout
	linesRendered := 0

	// Hide the cursor
	fmt.Fprint(w, "\033[?25l")
	defer fmt.Fprint(w, "\033[?25h")

	defer fmt.Println()

	for {
		select {
		case msg, ok := <-messageChannel:
			if !ok {
				return
			}

			if linesRendered > 1 {
				fmt.Fprintf(w, "\033[%dA", linesRendered-1)
			}

			fmt.Fprint(w, "\r")

			fmt.Fprint(w, "\033[J")

			fmt.Fprint(w, msg)

			width := getTermWidth()
			if width <= 0 {
				width = 80 // Fallback
			}

			linesRendered = 0
			for _, line := range strings.Split(msg, "\n") {
				linesForSegment := (len(line) + width - 1) / width
				if linesForSegment == 0 {
					// empty segment (e.g., "a\n\nb") still counts as one
					linesForSegment = 1
				}
				linesRendered += linesForSegment
			}

		case <-ctx.Done():
			return
		}
	}
}
