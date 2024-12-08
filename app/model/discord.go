package model

// A content of a Discord message
type DiscordContent struct {
	Title        string
	CurrentPrice uint64
	LowestPrice  uint64
}

// A body of a Discord message
type DiscordMessageBody struct {
	Content string `json:"content"`
}
