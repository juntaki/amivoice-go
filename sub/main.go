package main

import (
	"fmt"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
	"github.com/juntaki/amivoice-go"
	"github.com/webview/webview"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int16, 64)

	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	stream.Start()

	// Output file.
	out, err := os.Create("output.wav")
	if err != nil {
		log.Fatal(err)
	}
	out.Truncate(0)

	defer out.Close()
	e := wav.NewEncoder(out, 44100, 16, 1, 1)

	for {
		stream.Read()
		for _, i := range in {
			e.WriteFrame(i)
		}

		select {
		case <-sig:
			e.Close()
			stream.Stop()

			token := os.Getenv("ACP_TOKEN")
			c, err := amivoice.NewConnection(token)
			if err != nil {
				return
			}
			defer c.Close()

			out.Sync()
			buf, _ := ioutil.ReadFile("./output.wav")
			s := &amivoice.InterpretRequest{
				AudioFormat:      amivoice.CodecMuLaw,
				GrammarFileNames: amivoice.GammarFileGeneral,
				Data:             buf,
			}
			result, err := c.Interpret(s)
			if err != nil {
				panic(err)
			}
			fmt.Println(result)

			fmt.Println("aaa")
			return
		default:
		}
	}
}

func ui() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")
	w.Run()
}
