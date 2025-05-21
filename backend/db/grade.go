package db

import (
	"fmt"
	"time"

	"to-mrz/models"
)

// GetStudentGrades retrieves all grades for a specific student
func GetStudentGrades(studentID uint) ([]models.Enrollment, error) {
	enrollments := []models.Enrollment{}

	rows, err := DB.Query(`
		SELECT 
			e.id, e.student_id, e.course_offering_id, e.grade, e.status,
			e.created_at, e.updated_at,
			co.id as co_id, co.semester_id,
			c.id as c_id, c.name as course_name, c.code as course_code, c.credits,
			t.id as t_id, u.id as u_id, u.name as teacher_name,
			s.id as s_id, s.name as semester_name
		FROM enrollments e
		JOIN course_offerings co ON e.course_offering_id = co.id
		JOIN courses c ON co.course_id = c.id
		JOIN teachers t ON co.teacher_id = t.id
		JOIN users u ON t.user_id = u.id
		JOIN semesters s ON co.semester_id = s.id
		WHERE e.student_id = ?
	`, studentID)

	if err != nil {
		return nil, fmt.Errorf("failed to query student grades: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Enrollment
		var courseOffering models.CourseOffering
		var course models.Course
		var teacher models.Teacher
		var teacherUser models.User
		var semester models.Semester

		err := rows.Scan(
			&e.ID, &e.StudentID, &e.CourseOfferingID, &e.Grade, &e.Status,
			&e.CreatedAt, &e.UpdatedAt,
			&courseOffering.ID, &courseOffering.SemesterID,
			&course.ID, &course.Name, &course.Code, &course.Credits,
			&teacher.ID, &teacherUser.ID, &teacherUser.Name,
			&semester.ID, &semester.Name,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan student grade: %w", err)
		}

		teacher.User = teacherUser
		courseOffering.Course = course
		courseOffering.Teacher = teacher
		courseOffering.Semester = semester
		e.CourseOffering = courseOffering

		enrollments = append(enrollments, e)
	}

	return enrollments, nil
}

// GetCourseGrades retrieves all student grades for a specific course offering
func GetCourseGrades(courseOfferingID uint) ([]models.Enrollment, error) {
	enrollments := []models.Enrollment{}

	rows, err := DB.Query(`
		SELECT 
			e.id, e.student_id, e.course_offering_id, e.grade, e.status,
			e.created_at, e.updated_at,
			s.id as s_id, s.student_id as student_number,
			u.id as u_id, u.name as student_name,
			c.id as c_id, c.name as class_name
		FROM enrollments e
		JOIN students s ON e.student_id = s.id
		JOIN users u ON s.user_id = u.id
		JOIN classes c ON s.class_id = c.id
		WHERE e.course_offering_id = ?
	`, courseOfferingID)

	if err != nil {
		return nil, fmt.Errorf("failed to query course grades: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Enrollment
		var student models.Student
		var studentUser models.User
		var class models.Class

		err := rows.Scan(
			&e.ID, &e.StudentID, &e.CourseOfferingID, &e.Grade, &e.Status,
			&e.CreatedAt, &e.UpdatedAt,
			&student.ID, &student.StudentID,
			&studentUser.ID, &studentUser.Name,
			&class.ID, &class.Name,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan course grade: %w", err)
		}

		student.User = studentUser
		student.Class = class
		e.Student = student

		enrollments = append(enrollments, e)
	}

	return enrollments, nil
}

// UpdateGrade updates a student's grade for a specific enrollment
func UpdateGrade(id uint, grade float64) error {
	_, err := DB.Exec(`
		UPDATE enrollments 
		SET grade = ?, status = ?, updated_at = ?
		WHERE id = ?
	`, grade, getStatusFromGrade(grade), time.Now(), id)

	if err != nil {
		return fmt.Errorf("failed to update grade: %w", err)
	}

	return nil
}

// CreateGrade creates a new enrollment record with a grade
func CreateGrade(enrollment models.Enrollment) (uint, error) {
	result, err := DB.Exec(`
		INSERT INTO enrollments (student_id, course_offering_id, grade, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, enrollment.StudentID, enrollment.CourseOfferingID, enrollment.Grade, getStatusFromGrade(enrollment.Grade), time.Now(), time.Now())

	if err != nil {
		return 0, fmt.Errorf("failed to create grade: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get created grade ID: %w", err)
	}

	return uint(id), nil
}

// Helper function to determine status based on grade
func getStatusFromGrade(grade float64) string {
	if grade >= 60 {
		return "已完成"
	} else {
		return "未通过"
	}
}
