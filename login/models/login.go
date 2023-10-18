package models

type LoginPayload struct {
	Error   int64   `json:"error"`
	Payload Payload `json:"payload"`
	Type    string  `json:"type"`
}

type Payload struct {
	Sig string `json:"sig"`
	Ts  int64  `json:"ts"`
}
