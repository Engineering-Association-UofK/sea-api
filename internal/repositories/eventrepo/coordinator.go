package eventrepo

import (
	"fmt"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/models/tables"

	"github.com/jmoiron/sqlx"
)

func (r *EventRepository) CreateCoord(req *eventmodels.Coordinator) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, name, role)
	VALUES (:event_id, :name, :role)
	`, tables.EventCoords)
	res, err := r.db.NamedExec(query, req)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) UpdateCoord(coord *eventmodels.Coordinator) error {
	query := fmt.Sprintf(`
	UPDATE %s SET
		name = :name,
		role = :role
	WHERE id = :id AND event_id = :event_id
	`, tables.EventCoords)
	_, err := r.db.NamedExec(query, coord)
	return err
}

func (r *EventRepository) CreateBatchCoords(coords []eventmodels.Coordinator) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, name, role)
	VALUES (:event_id, :name, :role)
	`, tables.EventCoords)
	_, err := r.db.NamedExec(query, coords)
	return err
}

func (r *EventRepository) GetCoordsByEventID(eventID int64) ([]eventmodels.Coordinator, error) {
	var list = []eventmodels.Coordinator{}
	query := fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ?`, tables.EventCoords)
	err := r.db.Select(&list, query, eventID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) GetCoordsByEventIDs(eventIDs []int64) ([]eventmodels.Coordinator, error) {
	if len(eventIDs) == 0 {
		return []eventmodels.Coordinator{}, nil
	}

	query, args, err := sqlx.In(
		fmt.Sprintf(`SELECT * FROM %s WHERE event_id IN (?) ORDER BY id`, tables.EventCoords),
		eventIDs,
	)
	if err != nil {
		return nil, err
	}

	var list []eventmodels.Coordinator
	if err := r.db.Select(&list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) DeleteCoord(id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, tables.EventCoords)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EventRepository) DeleteCoordsByEventID(eventID int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE event_id = ?`, tables.EventCoords)
	_, err := r.db.Exec(query, eventID)
	return err
}
