package model

type WsMessage struct {
	Type int            `json:"type"`
	Body map[string]any `json:"body"`
}

type Body101 struct {
	Text      string
	Receiver  string
	MessageId string
}
