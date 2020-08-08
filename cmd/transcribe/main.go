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

	c, err := amivoice.NewConnection(setting.AppKey, setting.NoLog)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = c.Transcribe(setting.GenerateRecognitionConfig(f), os.Stdout)
	if err != nil {
		panic(err)
	}
}
