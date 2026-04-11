package repositories

import (
	"database/sql"
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type FormRepository struct {
	db *sqlx.DB
}

func NewFormRepository(db *sqlx.DB) *FormRepository {
	return &FormRepository{db: db}
}

// ======== FORM ANALYSIS ========

func (r *FormRepository) GetFormAnalysisData(formID int64) ([]models.FormAnalysisRow, error) {
	var rows []models.FormAnalysisRow
	query := fmt.Sprintf(`
    SELECT 
        q.id AS question_id,
        q.question_text,
        q.type,
        a.answer_value,
        COUNT(a.id) as answer_count
    FROM %s q
    JOIN %s p ON q.form_page_id = p.id
    LEFT JOIN %s a ON q.id = a.question_id
    LEFT JOIN %s r ON a.response_id = r.id
    WHERE p.form_id = ? AND r.status = 'SUBMITTED'
    GROUP BY q.id, a.answer_value
    ORDER BY p.page_num, q.display_order
    `, models.TableFormQuestions, models.TableFormPages, models.TableFormAnswers, models.TableFormResponses)

	err := r.db.Select(&rows, query, formID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ======== SPECIAL ========

func (r *FormRepository) GetFormWithQuestions(formID int64) ([]models.FormRow, error) {
	var rows []models.FormRow
	query := fmt.Sprintf(`
	SELECT 
		f.id AS form_id, f.title, f.description, f.allow_multiple, f.is_active, f.header_image_id,
		p.id AS page_id, p.page_num,
		q.id AS question_id, q.question_text, q.type, q.options, q.is_required, q.display_order
	FROM %s f
	LEFT JOIN %s p ON f.id = p.form_id
	LEFT JOIN %s q ON p.id = q.form_page_id
	WHERE f.id = ? 
	ORDER BY p.page_num ASC, q.display_order ASC;
`, models.TableForms, models.TableFormPages, models.TableFormQuestions)

	err := r.db.Select(&rows, query, formID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []models.FormRow{}, sql.ErrNoRows
	}
	return rows, nil
}

func (r *FormRepository) UpsertAnswer(answer *models.FormAnswerModel) error {
	query := fmt.Sprintf(`
    INSERT INTO %s (response_id, question_id, answer_value)
    VALUES (:response_id, :question_id, :answer_value)
    ON DUPLICATE KEY UPDATE answer_value = VALUES(answer_value)
    `, models.TableFormAnswers)
	_, err := r.db.NamedExec(query, answer)
	return err
}

func (r *FormRepository) CreateAnswersBatch(answers []models.FormAnswerModel) error {
	if len(answers) == 0 {
		return nil
	}
	query := fmt.Sprintf(`
    INSERT INTO %s (response_id, question_id, answer_value)
    VALUES (:response_id, :question_id, :answer_value)
    `, models.TableFormAnswers)
	_, err := r.db.NamedExec(query, answers)
	return err
}

// ======== CREATE ========

func (r *FormRepository) CreateForm(form *models.FormModel) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO %s (title, description, header_image_id, allow_multiple, is_active, created_by, created_at)
	VALUES (:title, :description, :header_image_id, :allow_multiple, :is_active, :created_by, :created_at)
	`, models.TableForms)
	res, err := r.db.NamedExec(query, form)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FormRepository) CreatePage(page *models.FormPageModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (form_id, page_num)
	VALUES (:form_id, :page_num)
	`, models.TableFormPages)
	res, err := r.db.NamedExec(query, page)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FormRepository) CreateQuestion(question *models.FormQuestionModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (form_page_id, question_text, type, options, is_required, display_order)
	VALUES (:form_page_id, :question_text, :type, :options, :is_required, :display_order)
	`, models.TableFormQuestions)
	res, err := r.db.NamedExec(query, question)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FormRepository) CreateQuestionsBatch(questions []models.FormQuestionModel) error {
	if len(questions) == 0 {
		return nil
	}
	query := fmt.Sprintf(`
	INSERT INTO %s (form_page_id, question_text, type, options, is_required, display_order)
	VALUES (:form_page_id, :question_text, :type, :options, :is_required, :display_order)
	`, models.TableFormQuestions)
	_, err := r.db.NamedExec(query, questions)
	return err
}

func (r *FormRepository) CreateResponse(response *models.FormResponseModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (form_id, user_id, status, submitted_at)
	VALUES (:form_id, :user_id, :status, :submitted_at)
	`, models.TableFormResponses)
	res, err := r.db.NamedExec(query, response)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FormRepository) CreateAnswer(answer *models.FormAnswerModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (response_id, question_id, answer_value)
	VALUES (:response_id, :question_id, :answer_value)
	`, models.TableFormAnswers)
	res, err := r.db.NamedExec(query, answer)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ======== GET ONE ========

func (r *FormRepository) GetFormByID(id int64) (*models.FormModel, error) {
	var form models.FormModel
	err := r.db.Get(&form, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableForms), id)
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *FormRepository) GetPageByID(id int64) (*models.FormPageModel, error) {
	var page models.FormPageModel
	err := r.db.Get(&page, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableFormPages), id)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *FormRepository) GetPageByFormIdAndPageNumber(formID int64, pageNumber int) (*models.FormPageModel, error) {
	var page models.FormPageModel
	err := r.db.Get(&page, fmt.Sprintf(`SELECT * FROM %s WHERE form_id = ? AND page_num = ?`, models.TableFormPages), formID, pageNumber)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *FormRepository) GetQuestionByID(id int64) (*models.FormQuestionModel, error) {
	var question models.FormQuestionModel
	err := r.db.Get(&question, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableFormQuestions), id)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *FormRepository) GetResponseByID(id int64) (*models.FormResponseModel, error) {
	var response models.FormResponseModel
	err := r.db.Get(&response, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableFormResponses), id)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ======== GET MANY ========

func (r *FormRepository) GetAllForms() ([]models.FormModel, error) {
	var forms []models.FormModel
	err := r.db.Select(&forms, fmt.Sprintf(`SELECT * FROM %s`, models.TableForms))
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *FormRepository) GetPagesByFormID(formID int64) ([]models.FormPageModel, error) {
	var pages []models.FormPageModel
	err := r.db.Select(&pages, fmt.Sprintf(`SELECT * FROM %s WHERE form_id = ? ORDER BY page_num ASC`, models.TableFormPages), formID)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *FormRepository) GetQuestionsByPageID(pageID int64) ([]models.FormQuestionModel, error) {
	var questions []models.FormQuestionModel
	err := r.db.Select(&questions, fmt.Sprintf(`SELECT * FROM %s WHERE form_page_id = ? ORDER BY display_order ASC`, models.TableFormQuestions), pageID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *FormRepository) GetResponsesByFormID(formID int64) ([]models.FormResponseModel, error) {
	var responses []models.FormResponseModel
	err := r.db.Select(&responses, fmt.Sprintf(`SELECT * FROM %s WHERE form_id = ?`, models.TableFormResponses), formID)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *FormRepository) GetNumberOfResponsesByFormID(formID int64) (int, error) {
	var count int
	err := r.db.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE form_id = ?`, models.TableFormResponses), formID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FormRepository) GetAnswersByResponseID(responseID int64) ([]models.FormAnswerModel, error) {
	var answers []models.FormAnswerModel
	err := r.db.Select(&answers, fmt.Sprintf(`SELECT * FROM %s WHERE response_id = ?`, models.TableFormAnswers), responseID)
	if err != nil {
		return nil, err
	}
	return answers, nil
}

// Get all required questions for form by formID
func (r *FormRepository) GetRequiredQuestionsByFormID(formID int64) ([]models.FormQuestionModel, error) {
	var questions []models.FormQuestionModel
	query := fmt.Sprintf(`
	SELECT q.* FROM %s q
	JOIN %s p ON q.form_page_id = p.id
	WHERE p.form_id = ? AND q.is_required = TRUE
	`, models.TableFormQuestions, models.TableFormPages)
	err := r.db.Select(&questions, query, formID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// Get All form questions
func (r *FormRepository) GetQuestionsByFormID(formID int64) ([]models.FormQuestionModel, error) {
	var questions []models.FormQuestionModel
	query := fmt.Sprintf(`
	SELECT q.* FROM %s q
	JOIN %s p ON q.form_page_id = p.id
	WHERE p.form_id = ?
	ORDER BY p.page_num ASC, q.display_order ASC
	`, models.TableFormQuestions, models.TableFormPages)
	err := r.db.Select(&questions, query, formID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// Get answers By responses IDs
func (r *FormRepository) GetAnswersByResponseIDs(responseIDs []int64) ([]models.FormAnswerModel, error) {
	if len(responseIDs) == 0 {
		return []models.FormAnswerModel{}, nil
	}
	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE response_id IN (?)`, models.TableFormAnswers), responseIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var answers []models.FormAnswerModel
	err = r.db.Select(&answers, query, args...)
	return answers, err
}

func (r *FormRepository) GetUserResponsesForForm(userID, formID int64) ([]models.FormResponseModel, error) {
	var responses []models.FormResponseModel
	err := r.db.Select(&responses, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ? AND form_id = ?`, models.TableFormResponses), userID, formID)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

// ======== UPDATE =========

func (r *FormRepository) UpdateForm(form *models.FormModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET title = :title, description = :description, allow_multiple = :allow_multiple, header_image_id = :header_image_id, is_active = :is_active
	WHERE id = :id
	`, models.TableForms)
	_, err := r.db.NamedExec(query, form)
	return err
}

func (r *FormRepository) UpdatePage(page *models.FormPageModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET page_num = :page_num
	WHERE id = :id
	`, models.TableFormPages)
	_, err := r.db.NamedExec(query, page)
	return err
}

func (r *FormRepository) UpdateQuestion(question *models.FormQuestionModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET question_text = :question_text, type = :type, options = :options, is_required = :is_required, display_order = :display_order
	WHERE id = :id
	`, models.TableFormQuestions)
	_, err := r.db.NamedExec(query, question)
	return err
}

func (r *FormRepository) UpdateQuestions(questions []models.FormQuestionModel) error {
	if len(questions) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
	UPDATE %s
	SET question_text = :question_text,
	    type = :type,
	    options = :options,
	    is_required = :is_required,
	    display_order = :display_order
	WHERE id = :id
	`, models.TableFormQuestions)

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

	for _, q := range questions {
		if _, err := stmt.Exec(q); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *FormRepository) UpdateAnswer(answer *models.FormAnswerModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET answer_value = :answer_value
	WHERE id = :id
	`, models.TableFormAnswers)
	_, err := r.db.NamedExec(query, answer)
	return err
}

func (r *FormRepository) UpdateResponseStatus(id int64, status models.ResponseStatus) error {
	_, err := r.db.Exec(fmt.Sprintf(`UPDATE %s SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, models.TableFormResponses), status, id)
	return err
}

// ======== DELETE =========

func (r *FormRepository) DeleteForm(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableForms), id)
	return err
}

func (r *FormRepository) DeletePage(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableFormPages), id)
	return err
}

func (r *FormRepository) DeleteQuestion(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableFormQuestions), id)
	return err
}

func (r *FormRepository) DeleteResponse(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableFormResponses), id)
	return err
}

func (r *FormRepository) DeleteAnswer(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableFormAnswers), id)
	return err
}
