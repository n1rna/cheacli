package term

import (
	"cheat/utils"

	"fmt"
	"io/ioutil"
	"os"

	markdown "github.com/n1rna/go-term-markdown"
	"github.com/tigrawap/slit"
)

func RenderMarkdownFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cols, _, err := utils.GetWinSize()
	if err != nil {
		panic(err)
	}

	columnsNum := 2
	result := markdown.Render(string(source), cols/columnsNum, columnsNum, 2)

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
