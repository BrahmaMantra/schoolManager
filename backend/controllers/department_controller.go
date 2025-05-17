package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"to-mrz/db"
	"to-mrz/models"

	"github.com/gin-gonic/gin"
)

// DepartmentRequest contains the department request data
type DepartmentRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// GetDepartments returns a list of all departments
func GetDepartments(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, code FROM departments")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}
	defer rows.Close()

	departments := []models.Department{}
	for rows.Next() {
		var department models.Department
		if err := rows.Scan(&department.ID, &department.Name, &department.Code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan department data"})
			return
		}
		departments = append(departments, department)
	}

	c.JSON(http.StatusOK, departments)
}

// GetDepartment returns a specific department by ID
func GetDepartment(c *gin.Context) {
	id := c.Param("id")
	departmentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	var department models.Department
	err = db.DB.QueryRow("SELECT id, name, code FROM departments WHERE id = ?", departmentID).
		Scan(&department.ID, &department.Name, &department.Code)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve department"})
		return
	}

	c.JSON(http.StatusOK, department)
}

// CreateDepartment creates a new department
func CreateDepartment(c *gin.Context) {
	var request DepartmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if code already exists
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM departments WHERE code = ?", request.Code).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check department code"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department code already exists"})
		return
	}

	// Insert new department
	result, err := db.DB.Exec(
		"INSERT INTO departments (name, code) VALUES (?, ?)",
		request.Name, request.Code,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get new department ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   id,
		"name": request.Name,
		"code": request.Code,
	})
}

// UpdateDepartment updates an existing department
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	departmentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	var request DepartmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if department exists
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM departments WHERE id = ?", departmentID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check department"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	// Check if code is already used by another department
	err = db.DB.QueryRow("SELECT COUNT(*) FROM departments WHERE code = ? AND id != ?", request.Code, departmentID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check department code"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department code already exists"})
		return
	}

	// Update department
	_, err = db.DB.Exec(
		"UPDATE departments SET name = ?, code = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		request.Name, request.Code, departmentID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   departmentID,
		"name": request.Name,
		"code": request.Code,
	})
}

// DeleteDepartment deletes a department
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	departmentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	// Check if department exists
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM departments WHERE id = ?", departmentID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check department"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	// Check if department has related majors
	err = db.DB.QueryRow("SELECT COUNT(*) FROM majors WHERE department_id = ?", departmentID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check related majors"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete department with related majors"})
		return
	}

	// Delete department
	_, err = db.DB.Exec("DELETE FROM departments WHERE id = ?", departmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete department"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})
}
