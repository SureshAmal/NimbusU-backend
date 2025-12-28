package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subjectRepository struct {
	db *pgxpool.Pool
}

func NewSubjectRepository(db *pgxpool.Pool) domain.SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(ctx context.Context, subject *domain.Subject) error {
	query := `
		INSERT INTO subjects (subject_id, subject_name, subject_code, department_id, credits, subject_type, description, syllabus, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`
	subject.SubjectID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		subject.SubjectID,
		subject.SubjectName,
		subject.SubjectCode,
		subject.DepartmentID,
		subject.Credits,
		subject.SubjectType,
		subject.Description,
		subject.Syllabus,
		subject.IsActive,
	).Scan(&subject.CreatedAt, &subject.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrSubjectCodeExists
		}
		return fmt.Errorf("failed to create subject: %w", err)
	}
	return nil
}

func (r *subjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subject, error) {
	query := `
		SELECT subject_id, subject_name, subject_code, department_id, credits, subject_type,
			   description, syllabus, is_active, created_at, updated_at
		FROM subjects
		WHERE subject_id = $1
	`
	var s domain.Subject
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.SubjectID, &s.SubjectName, &s.SubjectCode, &s.DepartmentID, &s.Credits, &s.SubjectType,
		&s.Description, &s.Syllabus, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrSubjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}
	return &s, nil
}

func (r *subjectRepository) GetByCode(ctx context.Context, code string) (*domain.Subject, error) {
	query := `
		SELECT subject_id, subject_name, subject_code, department_id, credits, subject_type,
			   description, syllabus, is_active, created_at, updated_at
		FROM subjects
		WHERE subject_code = $1
	`
	var s domain.Subject
	err := r.db.QueryRow(ctx, query, code).Scan(
		&s.SubjectID, &s.SubjectName, &s.SubjectCode, &s.DepartmentID, &s.Credits, &s.SubjectType,
		&s.Description, &s.Syllabus, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrSubjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subject by code: %w", err)
	}
	return &s, nil
}

func (r *subjectRepository) Update(ctx context.Context, subject *domain.Subject) error {
	query := `
		UPDATE subjects
		SET subject_name = $2, credits = $3, subject_type = $4, description = $5, syllabus = $6, is_active = $7, updated_at = now()
		WHERE subject_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		subject.SubjectID,
		subject.SubjectName,
		subject.Credits,
		subject.SubjectType,
		subject.Description,
		subject.Syllabus,
		subject.IsActive,
	).Scan(&subject.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrSubjectNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update subject: %w", err)
	}
	return nil
}

func (r *subjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE subjects SET is_active = false, updated_at = now() WHERE subject_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subject: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrSubjectNotFound
	}
	return nil
}

func (r *subjectRepository) List(ctx context.Context, filter domain.SubjectFilter, limit, offset int) ([]*domain.Subject, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.DepartmentID != nil {
		conditions = append(conditions, fmt.Sprintf("department_id = $%d", argNum))
		args = append(args, *filter.DepartmentID)
		argNum++
	}
	if filter.SubjectType != nil {
		conditions = append(conditions, fmt.Sprintf("subject_type = $%d", argNum))
		args = append(args, *filter.SubjectType)
		argNum++
	}
	if filter.Credits != nil {
		conditions = append(conditions, fmt.Sprintf("credits = $%d", argNum))
		args = append(args, *filter.Credits)
		argNum++
	}
	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf("(subject_name ILIKE $%d OR subject_code ILIKE $%d)", argNum, argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}
	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM subjects %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count subjects: %w", err)
	}

	// List query
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT subject_id, subject_name, subject_code, department_id, credits, subject_type,
			   description, syllabus, is_active, created_at, updated_at
		FROM subjects
		%s
		ORDER BY subject_code
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list subjects: %w", err)
	}
	defer rows.Close()

	var subjects []*domain.Subject
	for rows.Next() {
		var s domain.Subject
		if err := rows.Scan(
			&s.SubjectID, &s.SubjectName, &s.SubjectCode, &s.DepartmentID, &s.Credits, &s.SubjectType,
			&s.Description, &s.Syllabus, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan subject: %w", err)
		}
		subjects = append(subjects, &s)
	}

	return subjects, total, nil
}

func (r *subjectRepository) GetWithDetails(ctx context.Context, id uuid.UUID) (*domain.SubjectWithDetails, error) {
	query := `
		SELECT s.subject_id, s.subject_name, s.subject_code, s.department_id, s.credits, s.subject_type,
			   s.description, s.syllabus, s.is_active, s.created_at, s.updated_at,
			   d.department_id, d.department_name, d.department_code
		FROM subjects s
		JOIN departments d ON s.department_id = d.department_id
		WHERE s.subject_id = $1
	`

	var sd domain.SubjectWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sd.SubjectID, &sd.SubjectName, &sd.SubjectCode, &sd.DepartmentID, &sd.Credits, &sd.SubjectType,
		&sd.Description, &sd.Syllabus, &sd.IsActive, &sd.CreatedAt, &sd.UpdatedAt,
		&sd.Department.DepartmentID, &sd.Department.DepartmentName, &sd.Department.DepartmentCode,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrSubjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subject with details: %w", err)
	}

	// Get prerequisites
	prereqs, err := r.GetPrerequisites(ctx, id)
	if err != nil {
		return nil, err
	}
	sd.Prerequisites = prereqs

	// Get corequisites
	coreqs, err := r.GetCorequisites(ctx, id)
	if err != nil {
		return nil, err
	}
	sd.Corequisites = coreqs

	return &sd, nil
}

func (r *subjectRepository) AddPrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID, isMandatory bool) error {
	if subjectID == prerequisiteID {
		return domain.ErrSelfPrerequisite
	}

	query := `
		INSERT INTO subject_prerequisites (prerequisite_id, subject_id, prerequisite_subject_id, is_mandatory)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), subjectID, prerequisiteID, isMandatory)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("prerequisite already exists")
		}
		return fmt.Errorf("failed to add prerequisite: %w", err)
	}
	return nil
}

func (r *subjectRepository) RemovePrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID) error {
	query := `DELETE FROM subject_prerequisites WHERE subject_id = $1 AND prerequisite_subject_id = $2`
	result, err := r.db.Exec(ctx, query, subjectID, prerequisiteID)
	if err != nil {
		return fmt.Errorf("failed to remove prerequisite: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("prerequisite not found")
	}
	return nil
}

func (r *subjectRepository) GetPrerequisites(ctx context.Context, subjectID uuid.UUID) ([]domain.SubjectPrerequisite, error) {
	query := `
		SELECT sp.prerequisite_id, sp.subject_id, sp.prerequisite_subject_id, sp.is_mandatory, sp.created_at,
			   s.subject_code, s.subject_name
		FROM subject_prerequisites sp
		JOIN subjects s ON sp.prerequisite_subject_id = s.subject_id
		WHERE sp.subject_id = $1
	`
	rows, err := r.db.Query(ctx, query, subjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prerequisites: %w", err)
	}
	defer rows.Close()

	var prereqs []domain.SubjectPrerequisite
	for rows.Next() {
		var p domain.SubjectPrerequisite
		if err := rows.Scan(
			&p.PrerequisiteID, &p.SubjectID, &p.PrerequisiteSubjectID, &p.IsMandatory, &p.CreatedAt,
			&p.SubjectCode, &p.SubjectName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan prerequisite: %w", err)
		}
		prereqs = append(prereqs, p)
	}
	return prereqs, nil
}

func (r *subjectRepository) AddCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error {
	if subjectID == corequisiteID {
		return fmt.Errorf("subject cannot be its own corequisite")
	}

	query := `
		INSERT INTO subject_corequisites (corequisite_id, subject_id, corequisite_subject_id)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, uuid.New(), subjectID, corequisiteID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("corequisite already exists")
		}
		return fmt.Errorf("failed to add corequisite: %w", err)
	}
	return nil
}

func (r *subjectRepository) RemoveCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error {
	query := `DELETE FROM subject_corequisites WHERE subject_id = $1 AND corequisite_subject_id = $2`
	result, err := r.db.Exec(ctx, query, subjectID, corequisiteID)
	if err != nil {
		return fmt.Errorf("failed to remove corequisite: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("corequisite not found")
	}
	return nil
}

func (r *subjectRepository) GetCorequisites(ctx context.Context, subjectID uuid.UUID) ([]domain.SubjectBasic, error) {
	query := `
		SELECT s.subject_id, s.subject_code, s.subject_name, s.credits, s.subject_type
		FROM subject_corequisites sc
		JOIN subjects s ON sc.corequisite_subject_id = s.subject_id
		WHERE sc.subject_id = $1
	`
	rows, err := r.db.Query(ctx, query, subjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get corequisites: %w", err)
	}
	defer rows.Close()

	var coreqs []domain.SubjectBasic
	for rows.Next() {
		var c domain.SubjectBasic
		if err := rows.Scan(&c.SubjectID, &c.SubjectCode, &c.SubjectName, &c.Credits, &c.SubjectType); err != nil {
			return nil, fmt.Errorf("failed to scan corequisite: %w", err)
		}
		coreqs = append(coreqs, c)
	}
	return coreqs, nil
}
