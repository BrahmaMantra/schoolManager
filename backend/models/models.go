package models

import "time"

// 用户角色
type Role string

const (
	RoleAdmin      Role = "admin"      // 系统管理员
	RoleAcademic   Role = "academic"   // 教务处人员
	RoleDepartment Role = "department" // 院系管理员
	RoleTeacher    Role = "teacher"    // 教师
	RoleStudent    Role = "student"    // 学生
	RoleSupervisor Role = "supervisor" // 教学督导
	RoleFinance    Role = "finance"    // 财务人员
)

// User 用户信息
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"-"` // 密码不返回
	Name      string    `json:"name"`
	Role      Role      `json:"role"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Department 院系信息
type Department struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique"`
	Code      string    `json:"code" gorm:"unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Major 专业信息
type Major struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name"`
	Code         string     `json:"code" gorm:"unique"`
	DepartmentID uint       `json:"department_id"`
	Department   Department `json:"department" gorm:"foreignKey:DepartmentID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Class 班级信息
type Class struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Code      string    `json:"code" gorm:"unique"`
	MajorID   uint      `json:"major_id"`
	Major     Major     `json:"major" gorm:"foreignKey:MajorID"`
	Year      int       `json:"year"` // 入学年份
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Teacher 教师信息
type Teacher struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	UserID       uint       `json:"user_id"`
	User         User       `json:"user" gorm:"foreignKey:UserID"`
	DepartmentID uint       `json:"department_id"`
	Department   Department `json:"department" gorm:"foreignKey:DepartmentID"`
	Title        string     `json:"title"` // 职称
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Student 学生信息
type Student struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	StudentID  string    `json:"student_id" gorm:"unique"` // 学号
	ClassID    uint      `json:"class_id"`
	Class      Class     `json:"class" gorm:"foreignKey:ClassID"`
	EnrollYear int       `json:"enroll_year"` // 入学年份
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Course 课程信息
type Course struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Name          string     `json:"name"`
	Code          string     `json:"code" gorm:"unique"`
	Credits       float64    `json:"credits"`       // 学分
	Hours         int        `json:"hours"`         // 学时
	Type          string     `json:"type"`          // 必修/选修/实践
	DepartmentID  uint       `json:"department_id"` // 所属院系
	Department    Department `json:"department" gorm:"foreignKey:DepartmentID"`
	Description   string     `json:"description"`
	Prerequisites []Course   `json:"prerequisites" gorm:"many2many:course_prerequisites"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Semester 学期信息
type Semester struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Current   bool      `json:"current"` // 是否当前学期
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CourseOffering 开课信息
type CourseOffering struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CourseID    uint      `json:"course_id"`
	Course      Course    `json:"course" gorm:"foreignKey:CourseID"`
	SemesterID  uint      `json:"semester_id"`
	Semester    Semester  `json:"semester" gorm:"foreignKey:SemesterID"`
	TeacherID   uint      `json:"teacher_id"`
	Teacher     Teacher   `json:"teacher" gorm:"foreignKey:TeacherID"`
	Capacity    int       `json:"capacity"`    // 容量
	Location    string    `json:"location"`    // 教室
	Schedule    string    `json:"schedule"`    // 上课时间
	Status      string    `json:"status"`      // 状态
	Description string    `json:"description"` // 课程描述
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Enrollment 选课记录
type Enrollment struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	StudentID        uint           `json:"student_id"`
	Student          Student        `json:"student" gorm:"foreignKey:StudentID"`
	CourseOfferingID uint           `json:"course_offering_id"`
	CourseOffering   CourseOffering `json:"course_offering" gorm:"foreignKey:CourseOfferingID"`
	Grade            float64        `json:"grade"`  // 成绩
	Status           string         `json:"status"` // 状态：已选/已退选
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// 成绩组成
type GradeComponent struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CourseOfferingID uint           `json:"course_offering_id"`
	CourseOffering   CourseOffering `json:"course_offering" gorm:"foreignKey:CourseOfferingID"`
	Name             string         `json:"name"`   // 例如：期中考试、期末考试、平时成绩
	Weight           float64        `json:"weight"` // 权重，例如0.3表示30%
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// 学生成绩详情
type GradeDetail struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	EnrollmentID     uint           `json:"enrollment_id"`
	Enrollment       Enrollment     `json:"enrollment" gorm:"foreignKey:EnrollmentID"`
	GradeComponentID uint           `json:"grade_component_id"`
	GradeComponent   GradeComponent `json:"grade_component" gorm:"foreignKey:GradeComponentID"`
	Score            float64        `json:"score"` // 分数
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// 教学评估
type Evaluation struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CourseOfferingID uint           `json:"course_offering_id"`
	CourseOffering   CourseOffering `json:"course_offering" gorm:"foreignKey:CourseOfferingID"`
	StudentID        uint           `json:"student_id"`
	Student          Student        `json:"student" gorm:"foreignKey:StudentID"`
	Rating           int            `json:"rating"`  // 评分(1-5)
	Comment          string         `json:"comment"` // 评价内容
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// 教材信息
type Textbook struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	ISBN        string    `json:"isbn" gorm:"unique"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 教材申请
type TextbookApplication struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CourseOfferingID uint           `json:"course_offering_id"`
	CourseOffering   CourseOffering `json:"course_offering" gorm:"foreignKey:CourseOfferingID"`
	TextbookID       uint           `json:"textbook_id"`
	Textbook         Textbook       `json:"textbook" gorm:"foreignKey:TextbookID"`
	Quantity         int            `json:"quantity"` // 申请数量
	Status           string         `json:"status"`   // 状态：申请中/已批准/已拒绝
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// 实习信息
type Internship struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	StudentID  uint      `json:"student_id"`
	Student    Student   `json:"student" gorm:"foreignKey:StudentID"`
	Company    string    `json:"company"`
	Position   string    `json:"position"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Supervisor string    `json:"supervisor"` // 企业导师
	TeacherID  uint      `json:"teacher_id"` // 校内导师
	Teacher    Teacher   `json:"teacher" gorm:"foreignKey:TeacherID"`
	Status     string    `json:"status"` // 状态
	Report     string    `json:"report"` // 实习报告
	Grade      float64   `json:"grade"`  // 实习成绩
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// 毕业论文
type Thesis struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StudentID uint      `json:"student_id"`
	Student   Student   `json:"student" gorm:"foreignKey:StudentID"`
	TeacherID uint      `json:"teacher_id"` // 指导教师
	Teacher   Teacher   `json:"teacher" gorm:"foreignKey:TeacherID"`
	Title     string    `json:"title"`
	Abstract  string    `json:"abstract"`
	Status    string    `json:"status"` // 状态：选题中/进行中/已完成
	Grade     float64   `json:"grade"`  // 论文成绩
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
