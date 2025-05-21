<template>
  <div class="student-grades-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>我的成绩</span>
      </div>
      
      <el-table
        v-loading="loading"
        :data="grades"
        style="width: 100%"
        border>
        <el-table-column
          prop="courseOffering.course.code"
          label="课程代码"
          width="120">
        </el-table-column>
        <el-table-column
          prop="courseOffering.course.name"
          label="课程名称"
          width="180">
        </el-table-column>
        <el-table-column
          prop="courseOffering.teacher.user.name"
          label="任课教师"
          width="120">
        </el-table-column>
        <el-table-column
          prop="courseOffering.semester.name"
          label="学期"
          width="180">
        </el-table-column>
        <el-table-column
          prop="courseOffering.course.credits"
          label="学分"
          width="80">
        </el-table-column>
        <el-table-column
          prop="grade"
          label="成绩"
          width="80">
          <template slot-scope="scope">
            <span :class="{ 'grade-fail': scope.row.grade < 60 }">
              {{ scope.row.grade }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          label="成绩状态"
          width="100">
          <template slot-scope="scope">
            <el-tag :type="statusType(scope.row.status)">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="summary-info" v-if="grades.length > 0">
        <p>总平均分: <span class="grade-highlight">{{ averageGrade }}</span></p>
        <p>已获得学分: <span class="grade-highlight">{{ totalCredits }}</span></p>
      </div>
      
      <div class="empty-data" v-if="grades.length === 0 && !loading">
        <el-empty description="暂无成绩数据"></el-empty>
      </div>
    </el-card>
  </div>
</template>

<script>
import gradeApi from '@/api/grade'

export default {
  name: 'StudentGrades',
  data() {
    return {
      grades: [],
      loading: false
    }
  },
  computed: {
    averageGrade() {
      if (this.grades.length === 0) return 0
      const sum = this.grades.reduce((acc, curr) => acc + curr.grade, 0)
      return (sum / this.grades.length).toFixed(2)
    },
    totalCredits() {
      if (this.grades.length === 0) return 0
      return this.grades.reduce((acc, curr) => {
        // 只计算及格课程的学分
        if (curr.grade >= 60) {
          return acc + (curr.courseOffering.course.credits || 0)
        }
        return acc
      }, 0)
    }
  },
  methods: {
    statusType(status) {
      switch (status) {
        case '已完成': return 'success'
        case '未完成': return 'info'
        case '未通过': return 'danger'
        default: return 'info'
      }
    },
    async fetchGrades() {
      this.loading = true
      try {
        const response = await gradeApi.getStudentGrades()
        this.grades = response.data.data
      } catch (error) {
        this.$message.error('获取成绩信息失败')
        console.error('获取成绩失败:', error)
        
        // 如果是管理员或其他非学生角色，提供默认数据用于展示
        const currentRole = this.$store.getters.userRole
        if (['admin', 'academic', 'department', 'teacher'].includes(currentRole)) {
          this.grades = [
            {
              id: 1,
              courseOffering: {
                id: 1,
                course: {
                  id: 1,
                  name: '高等数学',
                  code: 'MATH101',
                  credits: 4.0
                },
                teacher: {
                  id: 1,
                  user: {
                    id: 3,
                    name: '王同学'
                  }
                },
                semester: {
                  id: 1,
                  name: '2023-2024学年第一学期'
                }
              },
              grade: 95.5,
              status: '已完成'
            },
            {
              id: 2,
              courseOffering: {
                id: 2,
                course: {
                  id: 2,
                  name: '线性代数',
                  code: 'MATH201',
                  credits: 3.0
                },
                teacher: {
                  id: 2,
                  user: {
                    id: 4,
                    name: '李同学'
                  }
                },
                semester: {
                  id: 1,
                  name: '2023-2024学年第一学期'
                }
              },
              grade: 88.0,
              status: '已完成'
            },
            {
              id: 3,
              courseOffering: {
                id: 3,
                course: {
                  id: 3,
                  name: '数据结构',
                  code: 'CS102',
                  credits: 4.0
                },
                teacher: {
                  id: 3,
                  user: {
                    id: 5,
                    name: '张同学'
                  }
                },
                semester: {
                  id: 1,
                  name: '2023-2024学年第一学期'
                }
              },
              grade: 92.5,
              status: '已完成'
            },
            {
              id: 4,
              courseOffering: {
                id: 4,
                course: {
                  id: 4,
                  name: '计算机网络',
                  code: 'CS301',
                  credits: 3.0
                },
                teacher: {
                  id: 4,
                  user: {
                    id: 6,
                    name: '刘同学'
                  }
                },
                semester: {
                  id: 2,
                  name: '2023-2024学年第二学期'
                }
              },
              grade: 78.5,
              status: '已完成'
            },
            {
              id: 5,
              courseOffering: {
                id: 5,
                course: {
                  id: 5,
                  name: '操作系统',
                  code: 'CS401',
                  credits: 4.0
                },
                teacher: {
                  id: 5,
                  user: {
                    id: 7,
                    name: '赵同学'
                  }
                },
                semester: {
                  id: 2,
                  name: '2023-2024学年第二学期'
                }
              },
              grade: 56.0,
              status: '未通过'
            }
          ]
          this.$message({
            type: 'info',
            message: '显示管理员演示数据',
            duration: 3000
          })
        }
      } finally {
        this.loading = false
      }
    }
  },
  created() {
    this.fetchGrades()
  }
}
</script>

<style scoped>
.student-grades-container {
  padding: 20px;
}
.summary-info {
  margin-top: 20px;
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
}
.grade-highlight {
  font-weight: bold;
  font-size: 18px;
  color: #409eff;
}
.grade-fail {
  color: #f56c6c;
  font-weight: bold;
}
.empty-data {
  margin: 30px 0;
}
</style> 