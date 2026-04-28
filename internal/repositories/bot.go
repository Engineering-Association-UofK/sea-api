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

func (r *BotRepository) Transaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *BotRepository) UpsertNodes(nodes []models.BotNode, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, node_type, pos_x, pos_y, is_locked, is_start)
		VALUES (:id, :node_type, :pos_x, :pos_y, :is_locked, :is_start)
		ON DUPLICATE KEY UPDATE 
			node_type = VALUES(node_type), 
			pos_x = VALUES(pos_x),
			pos_y = VALUES(pos_y),
			is_start = VALUES(is_start)
	`, models.TableBotNodes)

	if tx != nil {
		_, err := tx.NamedExec(query, nodes)
		return err
	}

	_, err := r.db.NamedExec(query, nodes)
	return err
}

func (r *BotRepository) UpsertNodeTranslations(translations []models.NodeTranslation, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (node_id, language, content)
		VALUES (:node_id, :language, :content)
		ON DUPLICATE KEY UPDATE 
			content = VALUES(content)
	`, models.TableBotNodeTranslations)

	if tx != nil {
		_, err := tx.NamedExec(query, translations)
		return err
	}

	_, err := r.db.NamedExec(query, translations)
	return err
}

func (r *BotRepository) UpsertEdges(edges []models.BotEdge, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, from_node_id, to_node_id, keyword)
		VALUES (:id, :from_node_id, :to_node_id, :keyword)
		ON DUPLICATE KEY UPDATE 
			from_node_id = VALUES(from_node_id), 
			to_node_id = VALUES(to_node_id), 
			keyword = VALUES(keyword)
	`, models.TableBotEdges)

	if tx != nil {
		_, err := tx.NamedExec(query, edges)
		return err
	}

	_, err := r.db.NamedExec(query, edges)
	return err
}

func (r *BotRepository) UpsertEdgeTranslations(translations []models.EdgeTranslation, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (edge_id, language, label)
		VALUES (:edge_id, :language, :label)
		ON DUPLICATE KEY UPDATE 
			label = VALUES(label)
	`, models.TableBotEdgeTranslations)

	if tx != nil {
		_, err := tx.NamedExec(query, translations)
		return err
	}

	_, err := r.db.NamedExec(query, translations)
	return err
}

func (r *BotRepository) UpsertActions(actions []models.BotAction, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (node_id, action_type, action_text)
		VALUES (:node_id, :action_type, :action_text)
		ON DUPLICATE KEY UPDATE 
			action_type = VALUES(action_type), 
			action_text = VALUES(action_text)
	`, models.TableBotActions)

	if tx != nil {
		_, err := tx.NamedExec(query, actions)
		return err
	}

	_, err := r.db.NamedExec(query, actions)
	return err
}

func (r *BotRepository) UpsertSession(sessionID string, nodeID string, userID *int64) error {
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

func (r *BotRepository) GetAction(nodeID string) (*models.BotAction, error) {
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

func (r *BotRepository) GetNodeRow(nodeID string, lang models.Language) (*models.NodeRow, error) {
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

func (r *BotRepository) GetNextNodeRow(nodeID string, lang models.Language, keyword string) (*models.NodeRow, error) {
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

func (r *BotRepository) GetParentOrStartNodeRow(nodeID string, lang models.Language, keyword string) (*models.NodeRow, error) {
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

func (r *BotRepository) GetEdgesForNode(nodeID string, lang models.Language) ([]models.EdgeRow, error) {
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
//// GET ALL ////
/////////////////

func (r *BotRepository) GetGraphNodes() ([]models.BotGraphNodeRow, error) {
	var nodes []models.BotGraphNodeRow
	query := fmt.Sprintf(`
		SELECT 
			n.id, n.node_type, n.pos_x, n.pos_y, n.is_locked, n.is_start,
			t.language, t.content,
			a.action_type, a.action_text
		FROM %s n
		LEFT JOIN %s t ON n.id = t.node_id
		LEFT JOIN %s a ON n.id = a.node_id
	`, models.TableBotNodes, models.TableBotNodeTranslations, models.TableBotActions)

	err := r.db.Select(&nodes, query)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *BotRepository) GetGraphEdges() ([]models.BotGraphEdgeRow, error) {
	var edges []models.BotGraphEdgeRow
	query := fmt.Sprintf(`
		SELECT 
			e.id, e.from_node_id, e.to_node_id, e.keyword,
			t.language, t.label
		FROM %s e
		LEFT JOIN %s t ON e.id = t.edge_id
	`, models.TableBotEdges, models.TableBotEdgeTranslations)

	err := r.db.Select(&edges, query)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

// Get locked nodes
func (r *BotRepository) GetLockedNodes() ([]models.NodeActionRow, error) {
	var nodes []models.NodeActionRow
	query := fmt.Sprintf(`
		SELECT 
			n.id, n.node_type, n.is_start, n.is_locked, n.created_at,
			a.action_type, a.action_text
		FROM %s n
		LEFT JOIN %s a ON n.id = a.node_id
		WHERE n.is_locked = TRUE
	`, models.TableBotNodes, models.TableBotActions)

	err := r.db.Select(&nodes, query)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

/////////////////
////  TOOLS  ////
/////////////////

func (r *BotRepository) ClearBotDatabase(tx *sqlx.Tx) error {
	tables := []models.TableName{
		models.TableBotEdgeTranslations,
		models.TableBotNodeTranslations,
		models.TableBotActions,
		models.TableBotEdges,
		models.TableBotNodes,
		models.TableBotUserStates,
	}

	// Disable foreign key checks to allow for a clean wipe of all tables
	if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}

	// Truncate
	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
			return err
		}
	}

	// Re-enable foreign key checks
	if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		return err
	}

	return nil
}

func (r *BotRepository) ResetDefault() error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.ClearBotDatabase(tx)
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
			return err
		}
	}

	return tx.Commit()
}
