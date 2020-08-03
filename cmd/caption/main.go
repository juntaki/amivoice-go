package main

import (
	"encoding/binary"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/gordonklaus/portaudio"
	"github.com/juntaki/amivoice-go"
	"github.com/juntaki/amivoice-go/cmd/lib"
	"io"
	"os"
	"os/signal"
)

func main() {
	setting, err := lib.ReadSetting()
	if err != nil {
		panic(err)
	}

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
	c, err := amivoice.NewConnection(setting.AppKey, true)
	if err != nil {
		return
	}
	defer c.Close()

	final := make(chan *amivoice.AEvent)
	progress := make(chan *amivoice.UEvent)

	go c.CollectResult(final, progress, nil)

	go c.Recognize(setting.GenerateRecognitionConfig(pr))

	finalText := ""
	currentText := ""

	app := app.New()
	app.Settings().SetTheme(myTheme{theme.DarkTheme()})

	w := app.NewWindow("Caption")
	w.Resize(fyne.Size{Width: 600, Height: 70})
	wd := widget.NewLabel(currentText)
	wd.Wrapping = fyne.TextWrapBreak
	sc := widget.NewVScrollContainer(wd)

	// Read result loop
	go func() {
		for {
			select {
			case val := <-final:
				finalText += val.Text
				currentText = finalText
			case val := <-progress:
				currentText = finalText + val.Text
			}
			wd.Text = currentText
			sc.Offset = fyne.NewPos(0, wd.Size().Height)
			wd.Refresh()
			sc.Refresh()
		}
	}()
	w.SetContent(sc)
	w.ShowAndRun()
}
