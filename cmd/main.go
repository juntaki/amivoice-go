package main

import (
	"fmt"
	"github.com/juntaki/amivoice-go"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	token := os.Getenv("ACP_TOKEN")
	c, err := amivoice.NewConnection(token)
	defer c.Close()

	b, err := ioutil.ReadFile("./test.mp3")
	if err != nil {
		log.Println("open:", err)
		return
	}

	s := &amivoice.InterpretRequest{
		AudioFormat:      amivoice.Codec16k,
		GrammarFileNames: amivoice.GammarFileGeneral,
		Data:             b,
	}
	result := c.Interpret(s)
	fmt.Println(result)
}
