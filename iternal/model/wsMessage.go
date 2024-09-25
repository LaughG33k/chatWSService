package model

type WsMessage struct {
	Type int            `json:"type"`
	Body map[string]any `json:"body"`
}

type Body101 struct {
	Text      string `json:"text"`
	Receiver  string `json:"recipient"`
	MessageId string `json:"messageId"`
}

type Body103 struct {
	WithWhom string `json:"withWhom"`
}

type Body104 struct {
	MessageId     string `json:"messageId"`
	WithWhom      string `json:"withWhom"`
	FlagDelForEvr bool   `json:"flagDelForEvr"`
}

type Body105 struct {
	MessageId   string `json:"messageId"`
	WithWhom    string `json:"withWhom"`
	UpdatedText string `json:"updatedText"`
}
