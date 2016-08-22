package client

// Command is a standard message format for clients sending commands to the
// server
type Command struct {
	Command string      `json:"cmd"`
	Value   interface{} `json:"val"`
}

// Message is a standard message format for server sending messages to the
// client
type Message struct {
	Type   string `json:"type"`
	Weight int    `json:"weight"`
}
