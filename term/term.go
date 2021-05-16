package term

import (
	"fmt"
	"io/ioutil"
	"os"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/tigrawap/slit"
	// "github.com/gomarkdown/markdown"
)

func ShowMDFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	result := markdown.Render(string(source), 80, 6)

	// fmt.Println(result)
	// extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	// parser := parser.NewWithExtensions(extensions)

	// // md := []byte("markdown text")
	// html := markdown.ToHTML(source, parser, nil)
	CheatOutput(result)
	// fmt.Println(string(result))
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

	// Probably should pass config to all slit constructors, with sane defaults
	s.SetFollow(false)
	s.SetKeepChars(0)

	s.Display()
}

func CheatOutput(output []byte) {
	// fmt.Println(string(output))

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

	// Probably should pass config to all slit constructors, with sane defaults
	s.SetFollow(false)
	s.SetKeepChars(0)

	s.Display()
}
