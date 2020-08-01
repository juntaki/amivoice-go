package amivoice

type UEvent struct {
	Results []struct {
		Tokens []struct {
			Written string `json:"written"`
		} `json:"tokens"`
		Text string `json:"text"`
	} `json:"results"`
	Text string `json:"text"`
}

type AEvent struct {
	Results []struct {
		Tokens []struct {
			Written    string  `json:"written"`
			Confidence float64 `json:"confidence"`
			StartTime  int     `json:"starttime"`
			EndTime    int     `json:"endtime"`
			Spoken     string  `json:"spoken"`
		} `json:"tokens"`
		Confidence float64       `json:"confidence"`
		StartTime  int           `json:"starttime"`
		EndTime    int           `json:"endtime"`
		Tags       []interface{} `json:"tags"`
		RuleName   string        `json:"rulename"`
		Text       string        `json:"text"`
	} `json:"results"`
	UtteranceID string `json:"utteranceid"`
	Text        string `json:"text"`
	Code        string `json:"code"`
	Message     string `json:"message"`
}
