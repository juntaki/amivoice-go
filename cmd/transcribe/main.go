package main

import (
	"flag"
	"fmt"
	"github.com/juntaki/amivoice-go/cmd/lib"
	"os"

	"github.com/juntaki/amivoice-go"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("invalid args")
		return
	}

	setting, err := lib.ReadSetting()
	if err != nil {
		panic(err)
	}

	c, err := amivoice.NewConnection(setting.AppKey, true)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	result, err := c.Transcribe(setting.GenerateRecognitionConfig(f))
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
