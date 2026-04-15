package models

// Bot Commands
type BotCommand struct {
	ID          int    `json:"id" db:"id"`
	Keyword     string `json:"keyword" db:"keyword"`
	Description string `json:"description" db:"description"`
}

type BotCommandTranslation struct {
	CommandID int    `json:"command_id" db:"command_id"`
	Text      string `json:"text" db:"text"`
	Language  string `json:"language" db:"language"`
}

type BotCommandTrigger struct {
	CommandID   int    `json:"command_id" db:"command_id"`
	TriggerText string `json:"trigger_text" db:"trigger_text"`
	Language    string `json:"language" db:"language"`
}

type BotCommandOption struct {
	CommandID   int    `json:"command_id" db:"command_id"`
	NextKeyword string `json:"next_keyword" db:"next_keyword"`
}
