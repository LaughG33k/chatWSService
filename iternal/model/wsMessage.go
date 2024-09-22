package model

type WsMessage struct {
	Type int            `json:"type"`
	Body map[string]any `json:"body"`
}

type Body101 struct {
	Text      string `json:"text"`
	Receiver  string `json:"recipient"`
	MessageId string `json:"message_id"`
}

type Body103 struct {
	WithWhom string `json:"with_whom"`
}

type Body104 struct {
	MessageId     string `json:"message_id"`
	WithWhom      string `json:"with_whom"`
	FlagDelForEvr string `json:"flag_del_for_evr"`
}
