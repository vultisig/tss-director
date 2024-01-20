package model

// Message is a struct that represents a message sent from one user to another.
type Message struct {
	SessionID string   `json:"session_id,omitempty"`
	From      string   `json:"from,omitempty"`
	To        []string `json:"to,omitempty"`
	Body      string   `json:"body,omitempty"`
}

type Session struct {
	SessionID    string   `json:"session_id,omitempty"`
	Participants []string `json:"participants,omitempty"`
}
