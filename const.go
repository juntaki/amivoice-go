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

type CodecType string

const (
	CodecLSB8k  CodecType = "lsb8k"
	CodecMSB8k            = "msb8k"
	CodecLSB11k           = "lsb11k"
	CodecMSB11k           = "msb11k"
	CodecLSB16k           = "lsb16k"
	CodecMSB16k           = "msb16k"
	CodecLSB22k           = "lsb22k"
	CodecMSB22k           = "msb22k"
	CodecMuLaw            = "mulaw"
	CodecALaw             = "alaw"
	Codec8k               = "8k"
	Codec16k              = "16k"
)

























