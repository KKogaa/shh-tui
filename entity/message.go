package entity

type Message struct {
	Msg  string `json:"msg"`
	User string `json:"user"`
	// Timestamp time.Time `json:"timestamp"`
}
