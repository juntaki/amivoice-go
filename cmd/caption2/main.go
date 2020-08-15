package main

import (
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/juntaki/amivoice-go"
	"github.com/juntaki/amivoice-go/cmd/lib"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

type Caption struct {
	f  *gtk.Fixed
	bg []*gtk.Label
}

func NewCaption() *Caption {
	f, err := gtk.FixedNew()
	if err != nil {
		log.Fatal("Unable to create fixed:", err)
	}

	bg := make([]*gtk.Label, 9)
	for i := range bg {
		bg[i], err = gtk.LabelNew("")
		if err != nil {
			log.Fatal("Unable to create label:", err)
		}
	}
	for _, b := range bg {
		b.SetWidthChars(200)
		b.SetEllipsize(pango.ELLIPSIZE_START)
		b.SetXAlign(0)
		b.SetYAlign(0)
	}
	line := 3
	m := [][]int{
		{-1, 1},
		{1, 1},
		{1, -1},
		{-1, -1},
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}
	for i := range m {
		f.Put(bg[i+1], m[i][0]*line, m[i][1]*line)
	}
	f.Put(bg[0], 0, 0)
	return &Caption{f: f, bg: bg}
}

func (c *Caption) setMessage(message string) {
	format := `<b><span foreground="%s" size="xx-large" lang="ja">%s</span></b>`
	c.bg[0].SetMarkup(fmt.Sprintf(format, "white", message))
	bg := fmt.Sprintf(format, "black", message)
	for _, b := range c.bg[1:] {
		b.SetMarkup(bg)
	}
}

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
	c, err := amivoice.NewConnection(setting.AppKey, setting.NoLog)
	if err != nil {
		return
	}
	defer c.Close()

	final := make(chan *amivoice.AEvent)
	progress := make(chan *amivoice.UEvent)

	go c.CollectResult(final, progress, nil)

	go c.Recognize(setting.GenerateRecognitionConfig(pr))

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetTitle("Transcribe")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.Connect("button-press-event", func(w *gtk.Window, e *gdk.Event) {
		ev := gdk.EventButton{Event: e}
		if ev.Button() != 1 {
			return
		}
		win.BeginMoveDrag(int(ev.Button()), int(ev.XRoot()), int(ev.YRoot()), uint32(time.Now().Unix()))
	})
	win.SetDecorated(false)
	win.SetKeepAbove(true)
	win.SetResizable(true)
	win.SetAppPaintable(true)

	finalText := "Caption->"
	currentText := "Caption->"

	cap := NewCaption()
	tick := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case val := <-final:
				finalText += val.Text
				currentText = finalText
			case val := <-progress:
				currentText = finalText + val.Text
			case <-tick.C:
				cap.setMessage(currentText)
			}
		}
	}()

	win.Add(cap.f)
	sc := win.GetScreen()
	v, err := sc.GetRGBAVisual()
	win.SetVisual(v)
	win.SetDefaultSize(0, 0)
	win.ShowAll()
	gtk.Main()
}
