package amivoice

import (
	"encoding/json"
	"errors"
	"io"
	"runtime"

	"github.com/gorilla/websocket"
)

type Conn struct {
	Conn     *websocket.Conn
	AppKey   string
	IsClosed bool
}

func (c *Conn) Close() error {
	c.IsClosed = true
	return c.Conn.Close()
}

type RecognitionConfig struct {
	AudioFormat      AudioFormat
	GrammarFileNames GrammarFile
	ProfileID        string
	ProfileWords     []ProfileWord
	Data             io.Reader
}

func (c *Conn) Transcribe(i *RecognitionConfig) (string, error) {
	err := c.Recognize(i)
	if err != nil {
		return "", err
	}
	res, err := c.CollectFinalResult()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *Conn) CollectResult(fixedResult chan<- *AEvent, progressResult chan<- *UEvent, notification chan<- string) error {
	for {
		err := c.CollectOneResult(fixedResult, progressResult, notification)
		if err == ErrEResponseReceived || err == ErrConnClosed {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

var (
	ErrConnClosed        = errors.New("read from closed connection")
	ErrEResponseReceived = errors.New("e response received")
)

func (c *Conn) CollectOneResult(fixedResult chan<- *AEvent, progressResult chan<- *UEvent, notification chan<- string) error {
	_, message, err := c.Conn.ReadMessage()
	if err != nil {
		if c.IsClosed {
			return ErrConnClosed
		}
		return err
	}

	if len(message) == 0 {
		return errors.New("invalid message")
	}
	switch message[0] {
	case 's':
		if notification != nil {
			notification <- string(message)
		}
		if len(message) > 1 {
			return errors.New(string(message))
		}
	case 'e':
		if notification != nil {
			notification <- string(message)
		}
		if len(message) > 1 {
			return errors.New(string(message))
		}
		return ErrEResponseReceived
	case 'p':
		return errors.New(string(message))
	case 'U':
		ret := UEvent{}
		err := json.Unmarshal(message[2:], &ret)
		if err != nil {
			return err
		}
		if progressResult != nil {
			progressResult <- &ret
		}
	case 'A':
		ret := AEvent{}
		err := json.Unmarshal(message[2:], &ret)
		if err != nil {
			return err
		}
		if fixedResult != nil {
			fixedResult <- &ret
		}
	case 'S':
		fallthrough
	case 'E':
		fallthrough
	case 'C':
		if notification != nil {
			notification <- string(message)
		}
	case 'G':
		// ignore
	default:
		return errors.New(string(message))
	}
	return nil
}

func (c *Conn) CollectFinalResult() (string, error) {
	final := make(chan *AEvent)
	var result string
	go func() {
		for f := range final {
			result += f.Text
		}
	}()

	err := c.CollectResult(final, nil, nil)
	if err != nil {
		return "", err
	}
	runtime.Gosched()
	return result, nil
}

func (c *Conn) Recognize(i *RecognitionConfig) error {
	s := &sCommand{
		AudioFormat:      i.AudioFormat,
		GrammarFileNames: i.GrammarFileNames,
		Authorization:    c.AppKey,
		ProfileID:        i.ProfileID,
		ProfileWords:     i.ProfileWords,
	}
	err := c.Conn.WriteMessage(websocket.TextMessage, s.Command())
	if err != nil {
		return err
	}
	for {
		w, err := c.Conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			return err
		}
		if _, err = w.Write([]byte("p")); err != nil {
			return err
		}
		_, err = io.CopyN(w, i.Data, 2048) // packet must be bigger than riff header?
		if err == io.EOF {
			e := &eCommand{}
			err = c.Conn.WriteMessage(websocket.TextMessage, e.Command())
			if err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}
		w.Close()
	}
	return nil
}

func NewConnection(appKey string, noLog bool) (*Conn, error) {
	url := wssLogURL
	if noLog {
		url = wssNoLogURL
	}
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: c, AppKey: appKey, IsClosed: false}, nil
}
