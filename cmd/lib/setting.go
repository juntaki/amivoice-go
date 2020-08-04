package lib

import (
	"github.com/juntaki/amivoice-go"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type RecognitionSetting struct {
	AppKey           string               `yaml:"app_key"`
	AudioFormat      amivoice.AudioFormat `yaml:"audio_format"`
	GrammarFileNames amivoice.GrammarFile `yaml:"grammar_file"`
	ProfileID        string               `yaml:"profile_id,omitempty"`
	ProfileWords     []ProfileWord        `yaml:"profile_words,omitempty"`
	NoLog            bool                 `yaml:"no_log,omitempty"`
}

type ProfileWord struct {
	Notation string `yaml:"notation"`
	Sound    string `yaml:"sound"`
}

func (r *RecognitionSetting) GenerateRecognitionConfig(data io.Reader) *amivoice.RecognitionConfig {
	profileWords := make([]amivoice.ProfileWord, len(r.ProfileWords))
	for i, p := range r.ProfileWords {
		profileWords[i] = amivoice.ProfileWord{
			Notation: p.Notation,
			Sound:    p.Sound,
		}
	}
	return &amivoice.RecognitionConfig{
		AudioFormat:      r.AudioFormat,
		GrammarFileNames: r.GrammarFileNames,
		ProfileID:        r.ProfileID,
		ProfileWords:     profileWords,
		Data:             data,
	}
}

func ReadSetting() (*RecognitionSetting, error) {
	f, err := os.Open("setting.yaml")
	if err != nil {
		return nil, err
	}
	d := yaml.NewDecoder(f)
	var ret RecognitionSetting

	err = d.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
