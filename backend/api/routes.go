package api

import (
	"to-mrz/controllers"
	"to-mrz/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Configure CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/login", controllers.Login)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/user", controllers.GetCurrentUser)

		// Department routes
		departments := protected.Group("/departments")
		{
			departments.GET("", controllers.GetDepartments)
			departments.GET("/:id", controllers.GetDepartment)
			departments.POST("", middleware.RoleMiddleware("admin", "academic", "department"), controllers.CreateDepartment)
			departments.PUT("/:id", middleware.RoleMiddleware("admin", "academic", "department"), controllers.UpdateDepartment)
			departments.DELETE("/:id", middleware.RoleMiddleware("admin", "academic"), controllers.DeleteDepartment)
		}

		// Course routes
		courses := protected.Group("/courses")
		{
			courses.GET("", controllers.GetCourses)
			courses.GET("/:id", controllers.GetCourse)
			courses.POST("", middleware.RoleMiddleware("admin", "academic", "department"), controllers.CreateCourse)
			courses.PUT("/:id", middleware.RoleMiddleware("admin", "academic", "department"), controllers.UpdateCourse)
			courses.DELETE("/:id", middleware.RoleMiddleware("admin", "academic"), controllers.DeleteCourse)
		}

		// Grade routes
		grades := protected.Group("/grades")
		{
			grades.GET("/student", controllers.GetStudentGrades)
			grades.GET("/course", middleware.RoleMiddleware("admin", "academic", "department", "teacher"), controllers.GetCourseGrades)
			grades.POST("", middleware.RoleMiddleware("admin", "academic", "department", "teacher"), controllers.CreateGrade)
			grades.PUT("/:id", middleware.RoleMiddleware("admin", "academic", "department", "teacher"), controllers.UpdateGrade)
		}

		// Future routes for majors, classes, courses, etc.
		// TODO: Implement these routes as we develop the controllers
	}

	return r
}
