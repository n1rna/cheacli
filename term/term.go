package term

import (
	"fmt"
	"io/ioutil"
	"os"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/tigrawap/slit"
)

func RenderMarkdownFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	result := markdown.Render(string(source), 80, 6)

	OutputStreamFromString(result)
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
