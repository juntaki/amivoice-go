package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/gordonklaus/portaudio"
	"github.com/juntaki/amivoice-go"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	pr, pw := io.Pipe()

	// PortAudio input with buffer
	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int16, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	stream.Start()
	go func() {
		for {
			for {
				f, err := stream.AvailableToRead()
				if err != nil {
					pw.CloseWithError(err)
				} else if f > 64 {
					break
				}
			}
			stream.Read()
			err = binary.Write(pw, binary.LittleEndian, in)
			if err != nil {
				pw.CloseWithError(err)
			}
			select {
			case <-sig:
				pw.Close()
				break
			default:
			}
		}
	}()

	// Initialize amivoice
	token := os.Getenv("ACP_TOKEN")
	c, err := amivoice.NewConnection(token)
	if err != nil {
		return
	}
	defer c.Close()

	final := make(chan *amivoice.AEvent)
	progress := make(chan *amivoice.UEvent)

	go c.CollectResult(final, progress, nil)

	// Read result loop
	go func() {
		for {
			select {
			case val := <-final:
				fmt.Println(val.Text)
			case val := <-progress:
				fmt.Println(val.Text)
			}
		}
	}()

	err = c.Recognize(&amivoice.RecognitionConfig{
		AudioFormat:      amivoice.AudioFormatLSB16k,
		GrammarFileNames: amivoice.GammarFileGeneral,
		Data:             pr,
	})
	if err != nil {
		panic(err)
	}
}
