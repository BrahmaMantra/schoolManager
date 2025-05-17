package db

import (
	"database/sql"
	"time"
	"to-mrz/models"
)

// GetAllCourses retrieves all courses with their departments
func GetAllCourses() ([]*models.Course, error) {
	rows, err := DB.Query(`
		SELECT c.id, c.name, c.code, c.credits, c.hours, c.type, c.department_id, c.description,
		       c.created_at, c.updated_at, d.id, d.name, d.code, d.created_at, d.updated_at
		FROM courses c
		LEFT JOIN departments d ON c.department_id = d.id
		ORDER BY c.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*models.Course
	for rows.Next() {
		var course models.Course
		var department models.Department
		var createdAt, updatedAt, deptCreatedAt, deptUpdatedAt string

		err := rows.Scan(
			&course.ID, &course.Name, &course.Code, &course.Credits, &course.Hours,
			&course.Type, &course.DepartmentID, &course.Description, &createdAt, &updatedAt,
			&department.ID, &department.Name, &department.Code, &deptCreatedAt, &deptUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		course.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		course.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		department.CreatedAt, _ = time.Parse(time.RFC3339, deptCreatedAt)
		department.UpdatedAt, _ = time.Parse(time.RFC3339, deptUpdatedAt)
		course.Department = department

		courses = append(courses, &course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// GetCourseByID retrieves a course by ID
func GetCourseByID(id uint) (*models.Course, error) {
	var course models.Course
	var department models.Department
	var createdAt, updatedAt, deptCreatedAt, deptUpdatedAt string

	err := DB.QueryRow(`
		SELECT c.id, c.name, c.code, c.credits, c.hours, c.type, c.department_id, c.description,
		       c.created_at, c.updated_at, d.id, d.name, d.code, d.created_at, d.updated_at
		FROM courses c
		LEFT JOIN departments d ON c.department_id = d.id
		WHERE c.id = ?
	`, id).Scan(
		&course.ID, &course.Name, &course.Code, &course.Credits, &course.Hours,
		&course.Type, &course.DepartmentID, &course.Description, &createdAt, &updatedAt,
		&department.ID, &department.Name, &department.Code, &deptCreatedAt, &deptUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	course.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	course.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	department.CreatedAt, _ = time.Parse(time.RFC3339, deptCreatedAt)
	department.UpdatedAt, _ = time.Parse(time.RFC3339, deptUpdatedAt)
	course.Department = department

	return &course, nil
}

// CreateCourse creates a new course
func CreateCourse(course *models.Course) (uint, error) {
	now := time.Now().Format(time.RFC3339)

	result, err := DB.Exec(`
		INSERT INTO courses (name, code, credits, hours, type, department_id, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, course.Name, course.Code, course.Credits, course.Hours, course.Type, course.DepartmentID, course.Description, now, now)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// UpdateCourse updates an existing course
func UpdateCourse(course *models.Course) error {
	now := time.Now().Format(time.RFC3339)

	_, err := DB.Exec(`
		UPDATE courses
		SET name = ?, code = ?, credits = ?, hours = ?, type = ?, department_id = ?, description = ?, updated_at = ?
		WHERE id = ?
	`, course.Name, course.Code, course.Credits, course.Hours, course.Type, course.DepartmentID, course.Description, now, course.ID)

	return err
}

// DeleteCourse deletes a course
func DeleteCourse(id uint) error {
	_, err := DB.Exec("DELETE FROM courses WHERE id = ?", id)
	return err
}
