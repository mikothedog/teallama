package models

// Message is chat history. Capital M = exported (public).
type Message struct {
	Role    string
	Content string
	IsUser  bool
}

// HistoryToString converts messages to Ollama context.
// Capital letter = can be used by other packages.
func HistoryToString(messages []Message) string {
	var s string
	for _, m := range messages {
		role := "Assistant"
		if m.IsUser {
			role = "User"
		}
		s += role + ": " + m.Content + "\n"
	}
	return s
}
