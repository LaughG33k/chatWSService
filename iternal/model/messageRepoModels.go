package model

type MessageForSave struct {
	SenderUuid   string
	ReceiverUuid string
	Text         string
	MessageId    string
	Time         int64
}

type MessageForEdit struct {
	Sender    string
	Recipient string
	MessageId string
	NewText   string
}

type MessageForDelete struct {
	Sender    string
	Receiver  string
	MessageId string
}

type MessageHistory struct {
	Received map[string]map[string]map[string]string `json:"received"`
	Sent     map[string]map[string]map[string]string `json:"sent"`
}
