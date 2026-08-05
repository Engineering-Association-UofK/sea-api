package eventrepo

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

func (r *EventRepository) CreateComponent(component *models.EventComponentModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`, models.TableEventComponents)
	res, err := r.db.NamedExec(query, &component)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateComponent(components []models.EventComponentModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`, models.TableEventComponents)
	if tx != nil {
		query = tx.Rebind(query)
		_, err := tx.NamedExec(query, components)
		return err
	}

	query, args, err := sqlx.Named(query, components)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) UpdateComponent(component *models.EventComponentModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET name = :name, description = :description, max_score = :max_score
	WHERE id = :id
	`, models.TableEventComponents)
	_, err := r.db.NamedExec(query, &component)
	return err
}

func (r *EventRepository) MassUpdateComponent(components []models.EventComponentModel) error {
	if len(components) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
	UPDATE %s
	SET name = :name,
	    description = :description,
	    max_score = :max_score
	WHERE id = :id
	`, models.TableEventComponents)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, c := range components {
		if _, err := stmt.Exec(c); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *EventRepository) GetComponentByID(id int64) (*models.EventComponentModel, error) {
	var component models.EventComponentModel
	err := r.db.Get(&component, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableEventComponents), id)
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *EventRepository) GetComponentsByEventID(eventID int64) ([]models.EventComponentModel, error) {
	var components []models.EventComponentModel
	err := r.db.Select(&components, fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ?`, models.TableEventComponents), eventID)
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *EventRepository) DeleteComponent(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEventComponents), id)
	return err
}

func (r *EventRepository) MassDeleteComponent(ids []int64) error {
	query, args, err := sqlx.In(fmt.Sprintf(`DELETE FROM %s WHERE id IN (?)`, models.TableEventComponents), ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}
