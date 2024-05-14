package entity

type Message struct {
	Msg      string `json:"msg"`
	User     string `json:"user"`
	Chatroom string `json:"chatroom"`
	// Timestamp time.Time `json:"timestamp"`
}
