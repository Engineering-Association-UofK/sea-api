package repositories

import (
	"database/sql"
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
	query := `
    SELECT 
        q.id AS question_id,
        q.question_text,
        q.type,
        a.answer_value,
        COUNT(a.id) as answer_count
    FROM form_questions q
    JOIN form_pages p ON q.form_page_id = p.id
    LEFT JOIN form_answers a ON q.id = a.question_id
    LEFT JOIN form_responses r ON a.response_id = r.id
    WHERE p.form_id = ? AND r.status = 'SUBMITTED'
    GROUP BY q.id, a.answer_value
    ORDER BY p.page_num, q.display_order
    `

	err := r.db.Select(&rows, query, formID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ======== SPECIAL ========

func (r *FormRepository) GetFormWithQuestions(formID int64) ([]models.FormRow, error) {
	var rows []models.FormRow
	query := `
	SELECT 
		f.id AS form_id, f.title, f.description, f.allow_multiple, f.is_active, f.header_image_id,
		p.id AS page_id, p.page_num,
		q.id AS question_id, q.question_text, q.type, q.options, q.is_required, q.display_order
	FROM forms f
	LEFT JOIN form_pages p ON f.id = p.form_id
	LEFT JOIN form_questions q ON p.id = q.form_page_id
	WHERE f.id = ? 
	ORDER BY p.page_num ASC, q.display_order ASC;
`

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
	query := `
    INSERT INTO form_answers (response_id, question_id, answer_value)
    VALUES (:response_id, :question_id, :answer_value)
    ON DUPLICATE KEY UPDATE answer_value = VALUES(answer_value)
    `
	_, err := r.db.NamedExec(query, answer)
	return err
}

func (r *FormRepository) CreateAnswersBatch(answers []models.FormAnswerModel) error {
	if len(answers) == 0 {
		return nil
	}
	query := `
    INSERT INTO form_answers (response_id, question_id, answer_value)
    VALUES (:response_id, :question_id, :answer_value)
    `
	_, err := r.db.NamedExec(query, answers)
	return err
}

// ======== CREATE ========

func (r *FormRepository) CreateForm(form *models.FormModel) (int64, error) {
	query := `INSERT INTO forms (title, description, header_image_id, allow_multiple, is_active, created_by, created_at)
	VALUES (:title, :description, :header_image_id, :allow_multiple, :is_active, :created_by, :created_at)
	`
	res, err := r.db.NamedExec(query, form)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FormRepository) CreatePage(page *models.FormPageModel) (int64, error) {
	query := `
	INSERT INTO form_pages (form_id, page_num)
	VALUES (:form_id, :page_num)
	`
	res, err := r.db.NamedExec(query, page)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FormRepository) CreateQuestion(question *models.FormQuestionModel) (int64, error) {
	query := `
	INSERT INTO form_questions (form_page_id, question_text, type, options, is_required, display_order)
	VALUES (:form_page_id, :question_text, :type, :options, :is_required, :display_order)
	`
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
	query := `
	INSERT INTO form_questions (form_page_id, question_text, type, options, is_required, display_order)
	VALUES (:form_page_id, :question_text, :type, :options, :is_required, :display_order)
	`
	_, err := r.db.NamedExec(query, questions)
	return err
}

func (r *FormRepository) CreateResponse(response *models.FormResponseModel) (int64, error) {
	query := `
	INSERT INTO form_responses (form_id, user_id, status, submitted_at)
	VALUES (:form_id, :user_id, :status, :submitted_at)
	`
	res, err := r.db.NamedExec(query, response)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FormRepository) CreateAnswer(answer *models.FormAnswerModel) (int64, error) {
	query := `
	INSERT INTO form_answers (response_id, question_id, answer_value)
	VALUES (:response_id, :question_id, :answer_value)
	`
	res, err := r.db.NamedExec(query, answer)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ======== GET ONE ========

func (r *FormRepository) GetFormByID(id int64) (*models.FormModel, error) {
	var form models.FormModel
	err := r.db.Get(&form, `SELECT * FROM forms WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *FormRepository) GetPageByID(id int64) (*models.FormPageModel, error) {
	var page models.FormPageModel
	err := r.db.Get(&page, `SELECT * FROM form_pages WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *FormRepository) GetPageByFormIdAndPageNumber(formID int64, pageNumber int) (*models.FormPageModel, error) {
	var page models.FormPageModel
	err := r.db.Get(&page, `SELECT * FROM form_pages WHERE form_id = ? AND page_num = ?`, formID, pageNumber)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *FormRepository) GetQuestionByID(id int64) (*models.FormQuestionModel, error) {
	var question models.FormQuestionModel
	err := r.db.Get(&question, `SELECT * FROM form_questions WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *FormRepository) GetResponseByID(id int64) (*models.FormResponseModel, error) {
	var response models.FormResponseModel
	err := r.db.Get(&response, `SELECT * FROM form_responses WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ======== GET MANY ========

func (r *FormRepository) GetAllForms() ([]models.FormModel, error) {
	var forms []models.FormModel
	err := r.db.Select(&forms, `SELECT * FROM forms`)
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *FormRepository) GetPagesByFormID(formID int64) ([]models.FormPageModel, error) {
	var pages []models.FormPageModel
	err := r.db.Select(&pages, `SELECT * FROM form_pages WHERE form_id = ? ORDER BY page_num ASC`, formID)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *FormRepository) GetQuestionsByPageID(pageID int64) ([]models.FormQuestionModel, error) {
	var questions []models.FormQuestionModel
	err := r.db.Select(&questions, `SELECT * FROM form_questions WHERE form_page_id = ? ORDER BY display_order ASC`, pageID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *FormRepository) GetResponsesByFormID(formID int64) ([]models.FormResponseModel, error) {
	var responses []models.FormResponseModel
	err := r.db.Select(&responses, `SELECT * FROM form_responses WHERE form_id = ?`, formID)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *FormRepository) GetNumberOfResponsesByFormID(formID int64) (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM form_responses WHERE form_id = ?`, formID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FormRepository) GetAnswersByResponseID(responseID int64) ([]models.FormAnswerModel, error) {
	var answers []models.FormAnswerModel
	err := r.db.Select(&answers, `SELECT * FROM form_answers WHERE response_id = ?`, responseID)
	if err != nil {
		return nil, err
	}
	return answers, nil
}

// Get all required questions for form by formID
func (r *FormRepository) GetRequiredQuestionsByFormID(formID int64) ([]models.FormQuestionModel, error) {
	var questions []models.FormQuestionModel
	query := `
	SELECT q.* FROM form_questions q
	JOIN form_pages p ON q.form_page_id = p.id
	WHERE p.form_id = ? AND q.is_required = TRUE
	`
	err := r.db.Select(&questions, query, formID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// Get All form questions
func (r *FormRepository) GetQuestionsByFormID(formID int64) ([]models.FormQuestionModel, error) {
	var questions []models.FormQuestionModel
	query := `
	SELECT q.* FROM form_questions q
	JOIN form_pages p ON q.form_page_id = p.id
	WHERE p.form_id = ?
	ORDER BY p.page_num ASC, q.display_order ASC
	`
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
	query, args, err := sqlx.In(`SELECT * FROM form_answers WHERE response_id IN (?)`, responseIDs)
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
	err := r.db.Select(&responses, `SELECT * FROM form_responses WHERE user_id = ? AND form_id = ?`, userID, formID)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

// ======== UPDATE =========

func (r *FormRepository) UpdateForm(form *models.FormModel) error {
	query := `
	UPDATE forms
	SET title = :title, description = :description, allow_multiple = :allow_multiple, header_image_id = :header_image_id, is_active = :is_active
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, form)
	return err
}

func (r *FormRepository) UpdatePage(page *models.FormPageModel) error {
	query := `
	UPDATE form_pages
	SET page_num = :page_num
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, page)
	return err
}

func (r *FormRepository) UpdateQuestion(question *models.FormQuestionModel) error {
	query := `
	UPDATE form_questions
	SET question_text = :question_text, type = :type, options = :options, is_required = :is_required, display_order = :display_order
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, question)
	return err
}

func (r *FormRepository) UpdateQuestions(questions []models.FormQuestionModel) error {
	if len(questions) == 0 {
		return nil
	}

	query := `
	UPDATE form_questions
	SET question_text = :question_text,
	    type = :type,
	    options = :options,
	    is_required = :is_required,
	    display_order = :display_order
	WHERE id = :id
	`

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
	query := `
	UPDATE form_answers
	SET answer_value = :answer_value
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, answer)
	return err
}

func (r *FormRepository) UpdateResponseStatus(id int64, status models.ResponseStatus) error {
	_, err := r.db.Exec(`UPDATE form_responses SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

// ======== DELETE =========

func (r *FormRepository) DeleteForm(id int64) error {
	_, err := r.db.Exec(`DELETE FROM forms WHERE id = ?`, id)
	return err
}

func (r *FormRepository) DeletePage(id int64) error {
	_, err := r.db.Exec(`DELETE FROM form_pages WHERE id = ?`, id)
	return err
}

func (r *FormRepository) DeleteQuestion(id int64) error {
	_, err := r.db.Exec(`DELETE FROM form_questions WHERE id = ?`, id)
	return err
}

func (r *FormRepository) DeleteResponse(id int64) error {
	_, err := r.db.Exec(`DELETE FROM form_responses WHERE id = ?`, id)
	return err
}

func (r *FormRepository) DeleteAnswer(id int64) error {
	_, err := r.db.Exec(`DELETE FROM form_answers WHERE id = ?`, id)
	return err
}
