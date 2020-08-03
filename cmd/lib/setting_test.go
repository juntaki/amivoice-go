package lib

import (
	"github.com/juntaki/amivoice-go"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	r := RecognitionSetting{
		AppKey:           "key",
		AudioFormat:      amivoice.AudioFormat16k,
		GrammarFileNames: amivoice.GammarFileGeneral,
		ProfileID:        "profile",
		ProfileWords:     []ProfileWord{
			{
				Notation: "aaa",
				Sound:    "あああ",
			},
			{
				Notation: "bbbb",
				Sound:    "いいい",
			},
		},
	}
	e := yaml.NewEncoder(os.Stdout)
	e.Encode(r)
}