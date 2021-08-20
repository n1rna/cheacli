package term

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"golang.org/x/sys/unix"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/tigrawap/slit"
)

func RenderMarkdownFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cols, _, err := getwinsize()
	if err != nil {
		panic(err)
	}
	fmt.Println("cols: ", cols)
	_ = markdown.Render(string(source), 72, 6)

	// OutputStreamFromString(result)
}

func OutputStream(ch chan string) {
	var s *slit.Slit
	var err error

	s, err = slit.NewFromStream(ch)
	defer s.Shutdown()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s.SetFollow(false)
	s.SetKeepChars(0)

	s.Display()
}

func OutputStreamFromString(output []byte) {
	ch := make(chan string)

	var s *slit.Slit
	var err error

	s, err = slit.NewFromStream(ch)
	defer s.Shutdown()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	go func() {
		ch <- string(output)
		close(ch)
	}()

	s.SetFollow(false)
	s.SetKeepChars(0)

	s.Display()
}

func getwinsize() (int, int, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return 0, 0, err
	}
	ws, err := unix.IoctlGetWinsize(int(tty.Fd()), syscall.TIOCGWINSZ)
	err2 := tty.Close()
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Col), int(ws.Row), err2
}
