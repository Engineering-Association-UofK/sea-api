package models

import (
	"time"
)

type NodeMapPath string
type NodeMap map[*NodeMapPath][]NodeRow

const (
	NodeNext NodeMapPath = "next"
	NodeCurr NodeMapPath = "curr"
	NodePrev NodeMapPath = "prev"
)

//////////////////
///   MODELS   ///
//////////////////

type NodeType string

const (
	NodeMessage NodeType = "message"
	NodeInput   NodeType = "input"
	NodeAction  NodeType = "action"
)

type FeedbackType string
type BotActionType string

const (
	FeedbackOrg       FeedbackType = "organization_feedback"
	FeedbackTechnical FeedbackType = "technical_feedback"
	FeedbackGeneral   FeedbackType = "general_feedback"
	FeedbackBug       FeedbackType = "bug_report"
	FeedbackOther     FeedbackType = "other_feedback"

	BotActionRedirect BotActionType = "redirect"
	BotActionFeedback BotActionType = "feedback"
)

var AllowedFeedbackTypes = map[FeedbackType]bool{
	FeedbackOrg:       true,
	FeedbackTechnical: true,
	FeedbackGeneral:   true,
	FeedbackBug:       true,
	FeedbackOther:     true,
}

// BotNode represents a state in the conversation graph
type BotNode struct {
	ID        int       `db:"id" json:"id"`
	Slug      string    `db:"slug" json:"slug"`
	Type      NodeType  `db:"node_type" json:"node_type"`
	IsLocked  bool      `db:"is_locked" json:"is_locked"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// BotEdge defines the transition between nodes
type BotEdge struct {
	ID         int    `db:"id" json:"id"`
	FromNodeID int    `db:"from_node_id" json:"from_node_id"`
	ToNodeID   int    `db:"to_node_id" json:"to_node_id"`
	Keyword    string `db:"keyword" json:"keyword"` // Internal logic identifier
}

// Translation structs for multi-language support
type NodeTranslation struct {
	NodeID   int    `db:"node_id"`
	Language string `db:"language"`
	Content  string `db:"content"`
}

type EdgeTranslation struct {
	EdgeID   int    `db:"edge_id"`
	Language string `db:"language"`
	Label    string `db:"label"`
}

// UserState handles the "Memory"
type UserState struct {
	SessionID     string    `db:"session_id"`
	CurrentNodeID int64     `db:"current_node_id"`
	UserID        *int64    `db:"user_id"` // Pointer for nullable field
	UpdatedAt     time.Time `db:"updated_at"`
}

type BotAction struct {
	ID         int           `db:"id"`
	NodeID     int           `db:"node_id"`
	ActionType BotActionType `db:"action_type"`
	ActionText string        `db:"action_text"`
}

//////////////////
///    ROWS    ///
//////////////////

type NodeRow struct {
	ID        int64     `db:"id"`
	Type      NodeType  `db:"node_type"`
	IsStart   bool      `db:"is_start"`
	IsLocked  bool      `db:"is_locked"`
	CreatedAt time.Time `db:"created_at"`
	Content   string    `db:"content"` // Joined from translations
}

type EdgeRow struct {
	ID         int    `db:"id"`
	FromNodeID int    `db:"from_node_id"`
	ToNodeID   int    `db:"to_node_id"`
	Keyword    string `db:"keyword"`
	Label      string `db:"label"`
}

//////////////////
///    DTOS    ///
//////////////////

type BotRequest struct {
	SessionID string   `json:"session_id" binding:"required"`
	Keyword   string   `json:"keyword" binding:"required"`
	Input     string   `json:"input"`
	Language  Language `json:"language" binding:"required"`
}

type BotOptionView struct {
	Label   string `json:"label"`
	Keyword string `json:"keyword"`
}

type BotResponse struct {
	NodeType NodeType        `json:"node_type"`
	Content  string          `json:"content"`
	Options  []BotOptionView `json:"options"`
	Metadata interface{}     `json:"metadata,omitempty"`
}

//////////////////
///   OTHERS   ///
//////////////////

type BotRedirect struct {
	ActionText string `json:"action_text"`
	IsInternal bool   `json:"is_internal"`
}
