package amivoice

const (
	httpLogURL   = "https://acp-api.amivoice.com/v1/recognize"
	httpNoLogURL = "https://acp-api.amivoice.com/v1/nolog/recognize"
	WssLogURL    = "wss://acp-api.amivoice.com/v1/"
	wssNoLogURL  = "wss://acp-api.amivoice.com/v1/nolog/"
)

type GrammarFile string

const (
	GammarFileGeneral      GrammarFile = "-a-general"
	GammarFileMedgeneral               = "-a-medgeneral"
	GammarFileBizmrreport              = "-a-bizmrreport"
	GammarFileBizfinance               = "-a-bizfinance"
	GammarFileMedcare                  = "-a-medcare"
	GammarFileMedkarte                 = "-a-medkarte"
	GammarFileMedpharmacy              = "-a-medpharmacy"
	GammarFileBizinsurance             = "-a-bizinsurance"
	GammarFileMeeting                  = "-a-meeting"
)

type AudioFormat string

const (
	AudioFormatLSB8k  AudioFormat = "lsb8k"
	AudioFormatMSB8k              = "msb8k"
	AudioFormatLSB11k             = "lsb11k"
	AudioFormatMSB11k             = "msb11k"
	AudioFormatLSB16k             = "lsb16k"
	AudioFormatMSB16k             = "msb16k"
	AudioFormatLSB22k             = "lsb22k"
	AudioFormatMSB22k             = "msb22k"
	AudioFormatMuLaw              = "mulaw"
	AudioFormatALaw               = "alaw"
	AudioFormat8k                 = "8k"
	AudioFormat16k                = "16k"
)
