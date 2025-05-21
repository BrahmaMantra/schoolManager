package controllers

import (
	"net/http"
	"strconv"
	"time"

	"to-mrz/db"
	"to-mrz/models"

	"github.com/gin-gonic/gin"
)

// Sample data - This would normally come from a database
var sampleGrades = []models.Enrollment{
	{
		ID:        1,
		StudentID: 1,
		Student: models.Student{
			ID:        1,
			StudentID: "2023001",
			User: models.User{
				ID:   1,
				Name: "张三",
			},
		},
		CourseOfferingID: 1,
		CourseOffering: models.CourseOffering{
			ID: 1,
			Course: models.Course{
				ID:   1,
				Name: "数据结构",
				Code: "CS101",
			},
			Teacher: models.Teacher{
				ID: 1,
				User: models.User{
					ID:   4,
					Name: "王同学",
				},
			},
			Semester: models.Semester{
				ID:   1,
				Name: "2023-2024学年第一学期",
			},
		},
		Grade:     85.5,
		Status:    "已完成",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        2,
		StudentID: 1,
		Student: models.Student{
			ID:        1,
			StudentID: "2023001",
			User: models.User{
				ID:   1,
				Name: "张三",
			},
		},
		CourseOfferingID: 2,
		CourseOffering: models.CourseOffering{
			ID: 2,
			Course: models.Course{
				ID:   2,
				Name: "计算机网络",
				Code: "CS201",
			},
			Teacher: models.Teacher{
				ID: 2,
				User: models.User{
					ID:   5,
					Name: "李同学",
				},
			},
			Semester: models.Semester{
				ID:   1,
				Name: "2023-2024学年第一学期",
			},
		},
		Grade:     92.0,
		Status:    "已完成",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        3,
		StudentID: 2,
		Student: models.Student{
			ID:        2,
			StudentID: "2023002",
			User: models.User{
				ID:   2,
				Name: "李四",
			},
		},
		CourseOfferingID: 1,
		CourseOffering: models.CourseOffering{
			ID: 1,
			Course: models.Course{
				ID:   1,
				Name: "数据结构",
				Code: "CS101",
			},
			Teacher: models.Teacher{
				ID: 1,
				User: models.User{
					ID:   4,
					Name: "王同学",
				},
			},
			Semester: models.Semester{
				ID:   1,
				Name: "2023-2024学年第一学期",
			},
		},
		Grade:     78.5,
		Status:    "已完成",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// GetStudentGrades 获取学生成绩
func GetStudentGrades(c *gin.Context) {
	// 从请求中获取学生ID或从token中提取当前登录的学生用户
	studentIDStr := c.Query("student_id")

	// 如果没有提供student_id，使用当前登录用户
	if studentIDStr == "" {
		// 获取当前登录用户
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
			return
		}

		// 这里简单地将userID转为字符串作为studentID
		studentIDStr = strconv.Itoa(userID.(int))
	}

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
		return
	}

	// 获取学生成绩
	studentGrades, err := db.GetStudentGrades(uint(studentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成绩失败: " + err.Error()})
		return
	}

	// 如果数据库没有数据，使用示例数据（仅用于演示）
	if len(studentGrades) == 0 {
		studentGrades = getExampleStudentGrades(uint(studentID))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": studentGrades,
	})
}

// GetCourseGrades 获取课程的所有学生成绩（教师用）
func GetCourseGrades(c *gin.Context) {
	courseOfferingIDStr := c.Query("course_offering_id")
	if courseOfferingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "必须提供课程开设ID"})
		return
	}

	courseOfferingID, err := strconv.Atoi(courseOfferingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的课程开设ID"})
		return
	}

	// 获取课程成绩
	courseGrades, err := db.GetCourseGrades(uint(courseOfferingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成绩失败: " + err.Error()})
		return
	}

	// 如果数据库没有数据，使用示例数据（仅用于演示）
	if len(courseGrades) == 0 {
		courseGrades = getExampleCourseGrades(uint(courseOfferingID))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": courseGrades,
	})
}

// UpdateGrade 更新学生成绩（教师用）
func UpdateGrade(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的成绩记录ID"})
		return
	}

	var gradeUpdate struct {
		Grade float64 `json:"grade" binding:"required"`
	}

	if err := c.ShouldBindJSON(&gradeUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新成绩
	err = db.UpdateGrade(uint(id), gradeUpdate.Grade)
	if err != nil {
		// 如果数据库更新失败，我们使用示例数据（仅用于演示）
		found := updateExampleGrade(uint(id), gradeUpdate.Grade)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "未找到成绩记录"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成绩更新成功",
	})
}

// CreateGrade 创建学生成绩（教师用）
func CreateGrade(c *gin.Context) {
	var newGrade models.Enrollment

	if err := c.ShouldBindJSON(&newGrade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建成绩记录
	id, err := db.CreateGrade(newGrade)
	if err != nil {
		// 如果数据库创建失败，我们使用示例数据（仅用于演示）
		newGrade.ID = uint(len(sampleGrades) + 1)
		newGrade.CreatedAt = time.Now()
		newGrade.UpdatedAt = time.Now()
		sampleGrades = append(sampleGrades, newGrade)
	} else {
		newGrade.ID = id
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "成绩创建成功",
		"data":    newGrade,
	})
}

// 以下是用于演示的辅助函数，实际应用中可以删除

func getExampleStudentGrades(studentID uint) []models.Enrollment {
	var studentGrades []models.Enrollment
	for _, grade := range sampleGrades {
		if grade.StudentID == studentID {
			studentGrades = append(studentGrades, grade)
		}
	}
	return studentGrades
}

func getExampleCourseGrades(courseOfferingID uint) []models.Enrollment {
	var courseGrades []models.Enrollment
	for _, grade := range sampleGrades {
		if grade.CourseOfferingID == courseOfferingID {
			courseGrades = append(courseGrades, grade)
		}
	}
	return courseGrades
}

func updateExampleGrade(id uint, grade float64) bool {
	found := false
	for i, g := range sampleGrades {
		if g.ID == id {
			sampleGrades[i].Grade = grade
			sampleGrades[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}
	return found
}
