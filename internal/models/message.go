package models

type Message struct {
	Role    string
	Content string
	IsUser  bool
}

func HistoryToString(messages []Message) string {
	var s string
	for _, msg := range messages {
		role := "Assistant"
		if msg.IsUser {
			role = "User"
		}
		s += role + ": " + msg.Content + "\n"
	}
	return s
}
