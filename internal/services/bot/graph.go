package bot

import (
	"fmt"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
)

func (s *BotService) GetBotGraph() (*models.BotGraphData, error) {
	nodeRows, err := s.repo.GetGraphNodes()
	if err != nil {
		slog.Debug("Failed to get graph nodes")
		return nil, err
	}

	edgeRows, err := s.repo.GetGraphEdges()
	if err != nil {
		slog.Debug("Failed to get graph edges")
		return nil, err
	}

	// Process Nodes
	nodesMap := make(map[string]*models.BotGraphNode)
	for _, row := range nodeRows {
		node, exists := nodesMap[row.ID]
		if !exists {
			node = &models.BotGraphNode{
				ID:   fmt.Sprintf("%s", row.ID),
				Type: row.Type,
				Position: models.BotGraphNodePosition{
					X: row.PosX,
					Y: row.PosY,
				},
				Data: models.BotGraphNodeData{
					Translations: make(map[models.Language]string),
					IsLocked:     row.IsLocked,
					IsStart:      row.IsStart,
				},
			}
			if row.ActionType.Valid {
				node.Data.Action = &models.BotGraphAction{
					Type: models.BotActionType(row.ActionType.String),
					Text: row.ActionText.String,
				}
			}
			nodesMap[row.ID] = node
		}
		if row.Language.Valid && row.Content.Valid {
			node.Data.Translations[models.Language(row.Language.String)] = row.Content.String
		}
	}

	// Process Edges
	edgesMap := make(map[string]*models.BotGraphEdge)
	for _, row := range edgeRows {
		edge, exists := edgesMap[row.ID]
		if !exists {
			edge = &models.BotGraphEdge{
				ID:      fmt.Sprintf("%s", row.ID),
				Source:  fmt.Sprintf("%s", row.FromNodeID),
				Target:  fmt.Sprintf("%s", row.ToNodeID),
				Keyword: row.Keyword,
				Data: models.BotGraphEdgeData{
					Translations: make(map[models.Language]string),
				},
			}
			edgesMap[row.ID] = edge
		}
		if row.Language.Valid {
			edge.Data.Translations[models.Language(row.Language.String)] = row.Label.String
		}
	}

	// Convert maps to slices
	nodes := make([]models.BotGraphNode, 0, len(nodesMap))
	for _, node := range nodesMap {
		nodes = append(nodes, *node)
	}

	edges := make([]models.BotGraphEdge, 0, len(edgesMap))
	for _, edge := range edgesMap {
		edges = append(edges, *edge)
	}

	return &models.BotGraphData{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

func (s *BotService) UpdateBotGraph(data *models.BotGraphData) error {
	err := s.ValidateGraph(data)
	if err != nil {
		slog.Debug("Validation failed")
		return err
	}

	// 1. Prepare data for batch upserts

	nodes := make([]models.BotNode, 0, len(data.Nodes))
	nodeTranslations := make([]models.NodeTranslation, 0, len(data.Nodes))
	actions := make([]models.BotAction, 0, len(data.Nodes))
	for _, node := range data.Nodes {
		nodes = append(nodes, models.BotNode{
			ID:       node.ID,
			Type:     node.Type,
			PosX:     node.Position.X,
			PosY:     node.Position.Y,
			IsLocked: node.Data.IsLocked,
			IsStart:  node.Data.IsStart,
		})

		for lang, content := range node.Data.Translations {
			nodeTranslations = append(nodeTranslations, models.NodeTranslation{
				NodeID:   node.ID,
				Language: string(lang),
				Content:  content,
			})
		}

		if node.Data.Action != nil {
			actions = append(actions, models.BotAction{
				NodeID:     node.ID,
				ActionType: node.Data.Action.Type,
				ActionText: node.Data.Action.Text,
			})
		}
		// slog.Debug("Node created", "id", id, "type", node.Type)
	}

	edges := make([]models.BotEdge, 0, len(data.Edges))
	edgeTranslations := make([]models.EdgeTranslation, 0, len(data.Edges))
	for _, edge := range data.Edges {

		edges = append(edges, models.BotEdge{
			ID:         edge.ID,
			FromNodeID: edge.Source,
			ToNodeID:   edge.Target,
			Keyword:    edge.Keyword,
		})

		for lang, label := range edge.Data.Translations {
			edgeTranslations = append(edgeTranslations, models.EdgeTranslation{
				EdgeID:   edge.ID,
				Language: string(lang),
				Label:    label,
			})
		}
		// slog.Debug("Edge created", "id", id, "fromID", fromID, "toID", toID)
	}

	// 2. Upsert everything inside a transaction
	tx, err := s.repo.Transaction()
	if err != nil {
		slog.Debug("Failed to start transaction")
		return err
	}
	defer tx.Rollback()

	err = s.repo.ResetDefault()
	if err != nil {
		slog.Debug("Failed to reset default")
		return err
	}

	err = s.repo.UpsertNodes(nodes, tx)
	if err != nil {
		slog.Debug("Failed to upsert nodes")
		return err
	}

	err = s.repo.UpsertNodeTranslations(nodeTranslations, tx)
	if err != nil {
		slog.Debug("Failed to upsert node translations")
		return err
	}

	err = s.repo.UpsertEdges(edges, tx)
	if err != nil {
		slog.Debug("Failed to upsert edges")
		return err
	}

	err = s.repo.UpsertEdgeTranslations(edgeTranslations, tx)
	if err != nil {
		slog.Debug("Failed to upsert edge translations")
		return err
	}

	err = s.repo.UpsertActions(actions, tx)
	if err != nil {
		slog.Debug("Failed to upsert actions")
		return err
	}

	return tx.Commit()
}

func (s *BotService) ValidateGraph(data *models.BotGraphData) error {
	err := s.CheckLocked(data)
	if err != nil {
		return err
	}
	msgs := map[string]string{}

	rootID := ""

	edgeIdsMap := make(map[string]bool)
	keysMap := make(map[string]bool)
	sourceKeyMap := make(map[string]string)

	edgeSourceMap := make(map[string]bool)
	edgeTargetMap := make(map[string]bool)

	// Loop through edges
	for i, edge := range data.Edges {
		edgeSourceMap[edge.Source] = true
		edgeTargetMap[edge.Target] = true

		// Check for duplicate edge IDs
		if edgeIdsMap[edge.ID] {
			msgs[fmt.Sprintf("edge_%s", edge.ID)] = "duplicate ID"
		}
		edgeIdsMap[edge.ID] = true

		// Check for dublicate keys
		if sourceKeyMap[edge.Keyword] == edge.Source {
			msgs[fmt.Sprintf("edge_%s", edge.ID)] = "duplicate keyword"
		}
		keysMap[edge.Keyword] = true
		sourceKeyMap[edge.Keyword] = edge.Source

		// Check for missing fields
		if edge.Source == "" || edge.Target == "" {
			msgs[fmt.Sprintf("edge_%s", edge.ID)] = "must have both source and target"
		}
		if edge.Keyword == "" {
			msgs[fmt.Sprintf("edge_%d", i)] = "missing a keyword"
		}
		if len(edge.Data.Translations) == 0 {
			msgs[fmt.Sprintf("edge_%s", edge.ID)] = "must provide all languages translations"
		} else {
			count := 0
			for lang, label := range edge.Data.Translations {
				if !models.AllowedLanguages[lang] {
					msgs[fmt.Sprintf("edge %s translation %s", edge.ID, lang)] = "invalid language"
				} else {
					count++
				}
				if label == "" {
					msgs[fmt.Sprintf("edge %s translation %s", edge.ID, lang)] = "label cannot be empty"
				}
			}
			if count != len(models.AllowedLanguages) {
				msgs[fmt.Sprintf("edge %s", edge.ID)] = "must provide all languages translations"
			}
		}
	}

	// Loop through all the nodes
	nodeIdMap := make(map[string]bool)
	for _, node := range data.Nodes {

		if node.Data.IsStart {
			rootID = node.ID
		}

		if !models.AllowedNodeTypes[node.Type] {
			msgs[fmt.Sprintf("node_%s", node.ID)] = "invalid node type"
		}

		if !edgeTargetMap[node.ID] {
			if !node.Data.IsStart {
				msgs[fmt.Sprintf("node_%s", node.ID)] = "node is unreachable (no incoming edges)"
			}
		}

		if !edgeSourceMap[node.ID] {
			if node.Type != models.NodeEnd && node.Type != models.NodeAction {
				msgs[fmt.Sprintf("node_%s", node.ID)] = "node is a dead end (no outgoing edges)"
			}
		}

		// Check for dublicate node IDs
		if nodeIdMap[node.ID] {
			msgs[fmt.Sprintf("node_%s", node.ID)] = "duplicate ID"
		}
		nodeIdMap[node.ID] = true

		// Check for missing fields
		if node.ID == "" {
			msgs[fmt.Sprintf("node_%s", node.ID)] = "missing an ID"
		}
		if len(node.Data.Translations) == 0 {
			msgs[fmt.Sprintf("node_%s", node.ID)] = "must provide all languages translations"
		} else {
			count := 0

			// Check for invalid languages and empty contents
			for lang, content := range node.Data.Translations {
				if !models.AllowedLanguages[lang] {
					msgs[fmt.Sprintf("node %s translation %s", node.ID, lang)] = "invalid language"
				} else {
					count++
				}
				if content == "" {
					msgs[fmt.Sprintf("node %s translation %s", node.ID, lang)] = "content cannot be empty"
				}
			}

			// Check for missing translations
			if count != len(models.AllowedLanguages) {
				msgs[fmt.Sprintf("node %s", node.ID)] = "must provide all languages translations"
			}
		}

		if node.Data.Action != nil {
			if node.Data.Action.Type == "" {
				msgs[fmt.Sprintf("node_%s_action", node.ID)] = "action type cannot be empty"
			}
			if !models.AllowedBotActionTypes[node.Data.Action.Type] {
				if !node.Data.IsLocked {
					msgs[fmt.Sprintf("node_%s_action", node.ID)] = "invalid action type"
				}
			}
			if node.Data.Action.Text == "" {
				msgs[fmt.Sprintf("node_%s_action", node.ID)] = "action text cannot be empty"
			}
		}

	}

	if len(msgs) > 0 {
		slog.Debug("Validation failed", "msgs", msgs)
		return errs.New(errs.MultiBadRequest, "validation errors", msgs)
	}

	return BFS(data, rootID)
}

func (s *BotService) CheckLocked(data *models.BotGraphData) error {
	msgs := map[string]string{}

	lockedNodes, err := s.repo.GetLockedNodes()
	lockedNodesMap := make(map[string]*models.NodeActionRow)
	if err != nil {
		slog.Debug("Failed to get locked nodes")
		return err
	}
	for _, node := range lockedNodes {
		lockedNodesMap[node.ID] = &node
	}
	edgesMap := make(map[string]*models.BotGraphEdge)
	for _, edge := range data.Edges {
		edgesMap[edge.Target] = &edge
	}

	count := 0
	for _, node := range data.Nodes {
		if locked, exists := lockedNodesMap[node.ID]; exists {
			count++
			if node.Type != locked.Type {
				msgs[fmt.Sprintf("node_%s", node.ID)] = "locked and its type cannot be changed"
			}
			if node.Data.Action != nil {
				if locked.ActionType.Valid && locked.ActionText.Valid {
					if node.Data.Action.Text != locked.ActionText.String {
						msgs[fmt.Sprintf("node_%s", node.ID)] = "locked and its action text cannot be changed"
					}
					if node.Data.Action.Type != models.BotActionType(locked.ActionType.String) {
						msgs[fmt.Sprintf("node_%s", node.ID)] = "locked and its action type cannot be changed"
					}
				} else {
					msgs[fmt.Sprintf("node_%s", node.ID)] = "locked and its action cannot be changed"
				}
			}

			if node.Data.IsLocked != locked.IsLocked {
				msgs[fmt.Sprintf("node_%s", node.ID)] = "locked and its locked status cannot be changed"
			}

			if _, exists := edgesMap[locked.ID]; !exists {
				msgs[fmt.Sprintf("node_%s", node.ID)] = "locked and cannot be removed or disconnected"
			}
		}
	}

	if count < len(lockedNodes) {
		msgs["locked_nodes_missing"] = "one or more locked nodes are missing from the graph"
	}

	if len(msgs) > 0 {
		slog.Debug("Validation failed", "msgs", msgs)
		return errs.New(errs.MultiBadRequest, "validation errors", msgs)
	}

	return nil
}

// Breadth-First Search to go down all the route in the bot tree
func BFS(data *models.BotGraphData, rootNodeID string) error {
	if rootNodeID == "" {
		return errs.New(errs.BadRequest, "No start node defined", nil)
	}

	// 1. Build an Adjacency List (Map of Source -> [Targets])
	// This fixes type/formatting mismatches and speeds up lookup
	adj := make(map[string][]string)
	for _, edge := range data.Edges {
		adj[edge.Source] = append(adj[edge.Source], edge.Target)
	}

	visited := make(map[string]bool)
	queue := []string{rootNodeID}
	visited[rootNodeID] = true

	// 2. Standard BFS
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Get only the neighbors for the current node
		for _, target := range adj[current] {
			if !visited[target] {
				visited[target] = true
				queue = append(queue, target)
			}
		}
	}

	// 3. Debugging Log
	if len(visited) != len(data.Nodes) {
		slog.Warn("Graph integrity check failed",
			"visited_count", len(visited),
			"total_nodes", len(data.Nodes),
			"root_node", rootNodeID)

		return errs.New(errs.BadRequest,
			fmt.Sprintf("Some nodes are unreachable. Visited %d out of %d nodes.",
				len(visited), len(data.Nodes)), nil)
	}

	return nil
}
