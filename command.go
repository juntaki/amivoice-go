package amivoice

import (
	"fmt"
	"strings"
)

type pCommand struct {
	Data []byte
}

func (p *pCommand) Command() []byte {
	ret := make([]byte, 0, len(p.Data)+1)
	ret = append(ret, byte('p'))
	ret = append(ret, p.Data...)
	return ret
}

type eCommand struct{}

func (e *eCommand) Command() []byte {
	return []byte("e")
}

type sCommand struct {
	AudioFormat      AudioFormat
	GrammarFileNames GrammarFile
	Authorization    string
	ProfileID        string
	ProfileWords     []ProfileWord
}

func (s *sCommand) Command() []byte {
	return []byte(s.String())
}

// s <audio_format> <grammar_file_names> <key>=<value> ...
func (s *sCommand) String() string {
	ret := fmt.Sprintf("s %s %s", s.AudioFormat, s.GrammarFileNames)
	ret += fmt.Sprintf(" authorization=%s", s.Authorization)
	if s.ProfileID != "" {
		ret += fmt.Sprintf(" profileId=:%s", s.ProfileID)
	}
	if len(s.ProfileWords) != 0 {
		pw := make([]string, len(s.ProfileWords))
		for i, p := range s.ProfileWords {
			pw[i] = p.String()
		}
		ret += fmt.Sprintf(`profileWords="%s"`, strings.Join(pw, "|"))
	}
	return ret
}

type ProfileWord struct {
	Notation string
	Sound    string
}

func (p *ProfileWord) String() string {
	return p.Notation + " " + p.Sound
}
