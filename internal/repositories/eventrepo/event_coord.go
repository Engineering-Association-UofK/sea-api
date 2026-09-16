package eventrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
)

func (r *EventRepository) CreateCoord(req *eventmodels.EventCoord) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, name, role)
	VALUES (:event_id, :name, :role)
	`, models.TableEventCoords)
	res, err := r.db.NamedExec(query, req)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) UpdateCoord(coord *eventmodels.EventCoord) error {
	query := fmt.Sprintf(`
	UPDATE %s SET
		name = :name,
		role = :role
	WHERE id = :id AND event_id = :event_id
	`, models.TableEventCoords)
	_, err := r.db.NamedExec(query, coord)
	return err
}

func (r *EventRepository) CreateBatchCoords(coords []eventmodels.EventCoord) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, name, role)
	VALUES (:event_id, :name, :role)
	`, models.TableEventCoords)
	_, err := r.db.NamedExec(query, coords)
	return err
}

func (r *EventRepository) GetCoordsByEventID(eventID int64) ([]eventmodels.EventCoord, error) {
	var list = []eventmodels.EventCoord{}
	query := fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ?`, models.TableEventCoords)
	err := r.db.Select(&list, query, eventID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) DeleteCoord(id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEventCoords)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EventRepository) DeleteCoordsByEventID(eventID int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE event_id = ?`, models.TableEventCoords)
	_, err := r.db.Exec(query, eventID)
	return err
}
