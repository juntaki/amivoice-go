package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/juntaki/amivoice-go"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("invalid args")
		return
	}
	token := os.Getenv("ACP_TOKEN")
	if token == "" {
		fmt.Println("set ACP_TOKEN")
		return
	}
	c, err := amivoice.NewConnection(token)
	defer c.Close()

	f, err := os.Open("./test.wav")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := &amivoice.RecognitionConfig{
		AudioFormat:      amivoice.AudioFormat16k,
		GrammarFileNames: amivoice.GammarFileGeneral,
		Data:             f,
	}
	result, err := c.Transcribe(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
