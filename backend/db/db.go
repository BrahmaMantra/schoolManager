package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"to-mrz/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database connection and creates the necessary tables
func InitDB() error {
	// Make sure the data directory exists
	dbDir := "./data"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	dbPath := filepath.Join(dbDir, "university.db")
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database successfully")

	// Create tables
	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// createTables creates all the necessary tables in the database
func createTables() error {
	// Users table
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		role TEXT NOT NULL,
		email TEXT,
		phone TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Departments table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS departments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		code TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Majors table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS majors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		code TEXT UNIQUE NOT NULL,
		department_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (department_id) REFERENCES departments(id)
	)`)
	if err != nil {
		return err
	}

	// Classes table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS classes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		code TEXT UNIQUE NOT NULL,
		major_id INTEGER NOT NULL,
		year INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (major_id) REFERENCES majors(id)
	)`)
	if err != nil {
		return err
	}

	// Teachers table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS teachers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER UNIQUE NOT NULL,
		department_id INTEGER NOT NULL,
		title TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (department_id) REFERENCES departments(id)
	)`)
	if err != nil {
		return err
	}

	// Students table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER UNIQUE NOT NULL,
		student_id TEXT UNIQUE NOT NULL,
		class_id INTEGER NOT NULL,
		enroll_year INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (class_id) REFERENCES classes(id)
	)`)
	if err != nil {
		return err
	}

	// Courses table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS courses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		code TEXT UNIQUE NOT NULL,
		credits REAL NOT NULL,
		hours INTEGER NOT NULL,
		type TEXT NOT NULL,
		department_id INTEGER NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (department_id) REFERENCES departments(id)
	)`)
	if err != nil {
		return err
	}

	// Course Prerequisites table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS course_prerequisites (
		course_id INTEGER NOT NULL,
		prerequisite_id INTEGER NOT NULL,
		PRIMARY KEY (course_id, prerequisite_id),
		FOREIGN KEY (course_id) REFERENCES courses(id),
		FOREIGN KEY (prerequisite_id) REFERENCES courses(id)
	)`)
	if err != nil {
		return err
	}

	// Semesters table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS semesters (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		current BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Course Offerings table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS course_offerings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		course_id INTEGER NOT NULL,
		semester_id INTEGER NOT NULL,
		teacher_id INTEGER NOT NULL,
		capacity INTEGER NOT NULL,
		location TEXT,
		schedule TEXT,
		status TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_id) REFERENCES courses(id),
		FOREIGN KEY (semester_id) REFERENCES semesters(id),
		FOREIGN KEY (teacher_id) REFERENCES teachers(id)
	)`)
	if err != nil {
		return err
	}

	// Enrollments table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS enrollments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER NOT NULL,
		course_offering_id INTEGER NOT NULL,
		grade REAL,
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (student_id) REFERENCES students(id),
		FOREIGN KEY (course_offering_id) REFERENCES course_offerings(id),
		UNIQUE(student_id, course_offering_id)
	)`)
	if err != nil {
		return err
	}

	// Grade Components table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS grade_components (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		course_offering_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		weight REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_offering_id) REFERENCES course_offerings(id)
	)`)
	if err != nil {
		return err
	}

	// Grade Details table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS grade_details (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		enrollment_id INTEGER NOT NULL,
		grade_component_id INTEGER NOT NULL,
		score REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (enrollment_id) REFERENCES enrollments(id),
		FOREIGN KEY (grade_component_id) REFERENCES grade_components(id),
		UNIQUE(enrollment_id, grade_component_id)
	)`)
	if err != nil {
		return err
	}

	// Evaluations table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS evaluations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		course_offering_id INTEGER NOT NULL,
		student_id INTEGER NOT NULL,
		rating INTEGER NOT NULL,
		comment TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_offering_id) REFERENCES course_offerings(id),
		FOREIGN KEY (student_id) REFERENCES students(id)
	)`)
	if err != nil {
		return err
	}

	// Textbooks table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS textbooks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT,
		publisher TEXT,
		isbn TEXT UNIQUE,
		price REAL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Textbook Applications table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS textbook_applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		course_offering_id INTEGER NOT NULL,
		textbook_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (course_offering_id) REFERENCES course_offerings(id),
		FOREIGN KEY (textbook_id) REFERENCES textbooks(id)
	)`)
	if err != nil {
		return err
	}

	// Internships table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS internships (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER NOT NULL,
		company TEXT NOT NULL,
		position TEXT NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		supervisor TEXT,
		teacher_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		report TEXT,
		grade REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (student_id) REFERENCES students(id),
		FOREIGN KEY (teacher_id) REFERENCES teachers(id)
	)`)
	if err != nil {
		return err
	}

	// Theses table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS theses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER NOT NULL,
		teacher_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		abstract TEXT,
		status TEXT NOT NULL,
		grade REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (student_id) REFERENCES students(id),
		FOREIGN KEY (teacher_id) REFERENCES teachers(id)
	)`)
	if err != nil {
		return err
	}

	return nil
}

// SeedDB seeds the database with sample data
func SeedDB() error {
	// Check if admin exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	// If admin doesn't exist, create default admin
	if count == 0 {
		_, err = DB.Exec(`
		INSERT INTO users (username, password, name, role, email) 
		VALUES ('admin', '$2a$10$XQiGvhANVHWBgfA.q3uuSO7XNt0XC/K/F2RIl8Yi8Ow/kH4pQJEFm', 'System Admin', 'admin', 'admin@university.edu')
		`)
		if err != nil {
			return err
		}
		log.Println("Default admin created with username 'admin' and password 'admin123'")
	}

	// Check if departments exist
	err = DB.QueryRow("SELECT COUNT(*) FROM departments").Scan(&count)
	if err != nil {
		return err
	}

	// If no departments exist, create sample departments
	if count == 0 {
		_, err = DB.Exec(`
		INSERT INTO departments (name, code) VALUES 
		('计算机学院', 'CS'),
		('数学学院', 'MATH'),
		('物理学院', 'PHY'),
		('经济管理学院', 'ECON')
		`)
		if err != nil {
			return err
		}
		log.Println("Sample departments created")
	}

	// Check if courses exist
	err = DB.QueryRow("SELECT COUNT(*) FROM courses").Scan(&count)
	if err != nil {
		return err
	}

	// If no courses exist, create sample courses
	if count == 0 {
		_, err = DB.Exec(`
		INSERT INTO courses (name, code, credits, hours, type, department_id, description) VALUES 
		('计算机导论', 'CS101', 3, 48, '必修课', 1, '计算机基础入门课程'),
		('数据结构', 'CS201', 4, 64, '必修课', 1, '数据结构基础课程'),
		('高等数学', 'MA101', 5, 80, '必修课', 2, '微积分、线性代数等基础数学知识'),
		('Web前端开发', 'CS301', 2.5, 40, '选修课', 1, 'HTML, CSS, JavaScript基础教学'),
		('数据库原理', 'CS202', 4, 64, '必修课', 1, '数据库设计与应用'),
		('计算机网络', 'CS203', 3.5, 56, '必修课', 1, '网络协议与架构'),
		('物理实验', 'PHY101', 2, 32, '实践课', 3, '基础物理实验'),
		('经济学原理', 'ECON101', 3, 48, '必修课', 4, '微观经济学与宏观经济学基础')
		`)
		if err != nil {
			return err
		}
		log.Println("Sample courses created")
	}

	return nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := DB.QueryRow(
		"SELECT id, username, password, name, role, email, phone FROM users WHERE username = ?",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.Role, &user.Email, &user.Phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}
