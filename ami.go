package amivoice

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
)

type Conn struct {
	conn  *websocket.Conn
	token string
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

type InterpretRequest struct {
	AudioFormat      CodecType
	GrammarFileNames GrammarFile
	ProfileID        string
	ProfileWords     []ProfileWord
	Data             []byte
}

func (c *Conn) Interpret(i *InterpretRequest) (string, error) {
	err := c.SendCommandSet(i)
	if err != nil {
		return "", err
	}
	res, err := c.WaitCommandSetEnd()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *Conn) WaitCommandSetEnd() (string, error) {
	var result string
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return "", err
		}

		switch message[0] {
		case 's':
			if len(message) > 1 {
				return "", errors.New(string(message))
			}
		case 'p':
			return "", errors.New(string(message))
		case 'e':
			if len(message) > 1 {
				return "", errors.New(string(message))
			}
			return result, nil
		case 'S':
			// ignore
		case 'E':
			// ignore
		case 'C':
			// ignore
		case 'U':
			// ignore
		case 'A':
			ret := AEvent{}
			err := json.Unmarshal(message[2:len(message)], &ret)
			if err != nil {
				panic(err)
			}
			result += ret.Text
		case 'G':
			// ignore
		default:
			return "", errors.New(string(message))
		}
	}
}

func (c *Conn) SendCommandSet(i *InterpretRequest) error {
	s := &sCommand{
		AudioFormat:      i.AudioFormat,
		GrammarFileNames: i.GrammarFileNames,
		Authorization:    c.token,
		ProfileID:        i.ProfileID,
		ProfileWords:     i.ProfileWords,
	}
	p := &pCommand{Data: i.Data}
	e := &eCommand{}
	err := c.conn.WriteMessage(websocket.TextMessage, s.Command())
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.BinaryMessage, p.Command())
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, e.Command())
	if err != nil {
		return err
	}
	return nil
}

func NewConnection(token string) (*Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(WssLogURL, nil)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: c, token: token}, nil
}
