package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type BotRepository struct {
	db *sqlx.DB
}

func NewBotRepository(db *sqlx.DB) *BotRepository {
	return &BotRepository{db: db}
}

////////////////
//// CREATE ////
////////////////

// Set/Update User State
func (r *BotRepository) UpsertSession(sessionID string, nodeID int64, userID *int64) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (session_id, current_node_id, user_id)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE current_node_id = VALUES(current_node_id), user_id = VALUES(user_id)
	`, models.TableBotUserStates)

	_, err := r.db.Exec(query, sessionID, nodeID, userID)
	return err
}

/////////////////
//// GET ONE ////
/////////////////

func (r *BotRepository) GetSession(sessionID string) (*models.UserState, error) {
	var state models.UserState
	query := fmt.Sprintf(`SELECT * FROM %s WHERE session_id = ?`, models.TableBotUserStates)
	err := r.db.Get(&state, query, sessionID)
	if err == sql.ErrNoRows {
		return nil, nil // No session exists yet
	}
	return &state, err
}

func (r *BotRepository) GetAction(nodeID int64) (*models.BotAction, error) {
	var action models.BotAction
	query := fmt.Sprintf(`SELECT * FROM %s WHERE node_id = ?`, models.TableBotActions)
	err := r.db.Get(&action, query, nodeID)
	return &action, err
}

func (r *BotRepository) GetStartNode(lang *models.Language) (*models.NodeRow, error) {
	var node models.NodeRow
	query := fmt.Sprintf(`
		SELECT n.id, n.node_type, n.is_start, n.is_locked, n.created_at, t.content
		FROM %s n
		LEFT JOIN %s t ON n.id = t.node_id AND t.language = ?
		WHERE is_start = TRUE
	`, models.TableBotNodes, models.TableBotNodeTranslations)

	err := r.db.Get(&node, query, lang)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *BotRepository) GetNodeRow(nodeID int64, lang models.Language) (*models.NodeRow, error) {
	var node models.NodeRow
	query := fmt.Sprintf(`
		SELECT n.id, n.node_type, n.is_start, n.is_locked, n.created_at, t.content
		FROM %s n
        LEFT JOIN %s t ON n.id = t.node_id AND t.language = ?
		WHERE n.id = ?`, models.TableBotNodes, models.TableBotNodeTranslations)

	err := r.db.Get(&node, query, lang, nodeID)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *BotRepository) GetNextNodeRow(nodeID int64, lang models.Language, keyword string) (*models.NodeRow, error) {
	var node models.NodeRow
	query := fmt.Sprintf(`
		SELECT n.id, n.node_type, n.is_start, n.is_locked, n.created_at, t.content
		FROM %s n
		JOIN %s e ON n.id = e.to_node_id
		LEFT JOIN %s t ON n.id = t.node_id AND t.language = ?
		WHERE e.from_node_id = ? AND e.keyword = ?`,
		models.TableBotNodes, models.TableBotEdges, models.TableBotNodeTranslations)

	err := r.db.Get(&node, query, lang, nodeID, keyword)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *BotRepository) GetParentOrStartNodeRow(nodeID int64, lang models.Language, keyword string) (*models.NodeRow, error) {
	var node models.NodeRow
	query := fmt.Sprintf(`
		SELECT n.id, n.node_type, n.is_start, n.is_locked, n.created_at, t.content
		FROM %s n
		JOIN %s e ON n.id = e.from_node_id
		LEFT JOIN %s t ON n.id = t.node_id AND t.language = ?
		WHERE e.to_node_id = ? AND e.keyword = ?`,
		models.TableBotNodes, models.TableBotEdges, models.TableBotNodeTranslations)

	err := r.db.Get(&node, query, lang, nodeID, keyword)
	if err != nil {
		if err == sql.ErrNoRows {
			return r.GetStartNode(&lang)
		}
		return nil, err
	}
	return &node, nil
}

func (r *BotRepository) GetEdgeRow(edgeID int64, lang models.Language) (*models.EdgeRow, error) {
	var edge models.EdgeRow
	query := fmt.Sprintf(`
		SELECT e.id, e.from_node_id, e.to_node_id, e.keyword, t.label
		FROM %s e
		LEFT JOIN %s t ON e.id = t.edge_id AND t.language = ?
		WHERE e.id = ?`, models.TableBotEdges, models.TableBotEdgeTranslations)

	err := r.db.Get(&edge, query, lang, edgeID)
	if err != nil {
		return nil, err
	}
	return &edge, nil
}

func (r *BotRepository) GetEdgesForNode(nodeID int64, lang models.Language) ([]models.EdgeRow, error) {
	var edges []models.EdgeRow
	query := fmt.Sprintf(`
		SELECT e.id, e.from_node_id, e.to_node_id, e.keyword, t.label
		FROM %s e
		LEFT JOIN %s t ON e.id = t.edge_id AND t.language = ?
		WHERE e.from_node_id = ?`, models.TableBotEdges, models.TableBotEdgeTranslations)

	err := r.db.Select(&edges, query, lang, nodeID)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

/////////////////
////  TOOLS  ////
/////////////////

func (r *BotRepository) ResetDefault() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	resourcesPath := config.App.ResourcesDir
	var sqlContent []byte
	sqlContent, err = os.ReadFile(resourcesPath + "/bot_default.sql")
	if err != nil {
		return err
	}

	queries := strings.Split(string(sqlContent), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
