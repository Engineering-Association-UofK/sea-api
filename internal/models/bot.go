package models

import (
	"database/sql"
	"time"
)

type NodeMapPath string
type NodeMap map[*NodeMapPath][]NodeRow

const (
	NodeNext NodeMapPath = "next"
	NodeCurr NodeMapPath = "curr"
	NodePrev NodeMapPath = "prev"
)

type NodeType string

const (
	NodeMessage NodeType = "message"
	NodeInput   NodeType = "input"
	NodeAction  NodeType = "action"
	NodeEnd     NodeType = "end"
)

var AllowedNodeTypes = map[NodeType]bool{
	NodeMessage: true,
	NodeInput:   true,
	NodeAction:  true,
	NodeEnd:     true,
}

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

var AllowedBotActionTypes = map[BotActionType]bool{
	BotActionRedirect: true,
}

var AllowedFeedbackTypes = map[FeedbackType]bool{
	FeedbackOrg:       true,
	FeedbackTechnical: true,
	FeedbackGeneral:   true,
	FeedbackBug:       true,
	FeedbackOther:     true,
}

//////////////////
///   MODELS   ///
//////////////////

// BotNode represents a state in the conversation graph
type BotNode struct {
	ID        string    `db:"id" json:"id"`
	Type      NodeType  `db:"node_type" json:"node_type"`
	PosX      int       `db:"pos_x" json:"pos_x"`
	PosY      int       `db:"pos_y" json:"pos_y"`
	IsLocked  bool      `db:"is_locked" json:"is_locked"`
	IsStart   bool      `db:"is_start" json:"is_start"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// BotEdge defines the transition between nodes
type BotEdge struct {
	ID         string `db:"id" json:"id"`
	FromNodeID string `db:"from_node_id" json:"from_node_id"`
	ToNodeID   string `db:"to_node_id" json:"to_node_id"`
	Keyword    string `db:"keyword" json:"keyword"` // Internal logic identifier
}

// Translation structs for multi-language support
type NodeTranslation struct {
	NodeID   string `db:"node_id"`
	Language string `db:"language"`
	Content  string `db:"content"`
}

type EdgeTranslation struct {
	EdgeID   string `db:"edge_id"`
	Language string `db:"language"`
	Label    string `db:"label"`
}

// UserState handles the "Memory"
type UserState struct {
	SessionID     string    `db:"session_id"`
	CurrentNodeID string    `db:"current_node_id"`
	UserID        *int64    `db:"user_id"` // Pointer for nullable field
	UpdatedAt     time.Time `db:"updated_at"`
}

type BotAction struct {
	NodeID     string        `db:"node_id"`
	ActionType BotActionType `db:"action_type"`
	ActionText string        `db:"action_text"`
}

//////////////////
///    ROWS    ///
//////////////////

type NodeRow struct {
	ID        string    `db:"id"`
	Type      NodeType  `db:"node_type"`
	IsStart   bool      `db:"is_start"`
	IsLocked  bool      `db:"is_locked"`
	CreatedAt time.Time `db:"created_at"`
	Content   string    `db:"content"` // Joined from translations
}

// node with actions row
type NodeActionRow struct {
	NodeRow
	ActionType sql.NullString `db:"action_type"`
	ActionText sql.NullString `db:"action_text"`
}

type EdgeRow struct {
	ID         string `db:"id"`
	FromNodeID string `db:"from_node_id"`
	ToNodeID   string `db:"to_node_id"`
	Keyword    string `db:"keyword"`
	Label      string `db:"label"`
}

type BotGraphNodeRow struct {
	ID       string   `db:"id"`
	Type     NodeType `db:"node_type"`
	PosX     int      `db:"pos_x"`
	PosY     int      `db:"pos_y"`
	IsLocked bool     `db:"is_locked"`
	IsStart  bool     `db:"is_start"`

	// Joined Translation Data
	Language sql.NullString `db:"language"`
	Content  sql.NullString `db:"content"`

	// Joined Action Data
	ActionType sql.NullString `db:"action_type"`
	ActionText sql.NullString `db:"action_text"`
}

type BotGraphEdgeRow struct {
	ID         string         `db:"id"`
	FromNodeID string         `db:"from_node_id"`
	ToNodeID   string         `db:"to_node_id"`
	Keyword    string         `db:"keyword"`
	Language   sql.NullString `db:"language"`
	Label      sql.NullString `db:"label"`
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

type BotGraphData struct {
	Nodes []BotGraphNode `json:"nodes"`
	Edges []BotGraphEdge `json:"edges"`
}

type BotGraphNode struct {
	ID       string               `json:"id"`
	Type     NodeType             `json:"type"`
	Position BotGraphNodePosition `json:"position"`
	Data     BotGraphNodeData     `json:"data"`
}

type BotGraphNodePosition struct {
	X int `json:"pos_x"`
	Y int `json:"pos_y"`
}

type BotGraphNodeData struct {
	Translations map[Language]string `json:"translations"`
	IsStart      bool                `json:"is_start"`
	IsLocked     bool                `json:"is_locked"`
	Action       *BotGraphAction     `json:"action,omitempty"`
}

type BotGraphAction struct {
	Type BotActionType `json:"type"`
	Text string        `json:"text"`
}

type BotGraphEdge struct {
	ID      string           `json:"id"`
	Source  string           `json:"source"`
	Target  string           `json:"target"`
	Keyword string           `json:"keyword"`
	Data    BotGraphEdgeData `json:"data"`
}

type BotGraphEdgeData struct {
	Translations map[Language]string `json:"translations"`
}

//////////////////
///   OTHERS   ///
//////////////////

type BotRedirect struct {
	ActionText string `json:"action_text"`
	IsInternal bool   `json:"is_internal"`
}

// {
//   "nodes": [
//     {
//       "id": "node_1",
//       "type": "message",
//       "position": { "x": 250, "y": 5 },
//       "data": {
//         "translations": {
//           "en": "Hello!",
//           "ar": "مرحباً!"
//         }
//       }
//     }
//   ],
//   "edges": [
//     {
//       "id": "edge_1_2",
//       "source": "node_1",
//       "target": "node_2",
//       "keyword": "start",
//       "data": {
//         "translations": {
//           "en": "Start",
//           "ar": "ابدأ"
//         }
//       }
//     }
//   ]
// }
