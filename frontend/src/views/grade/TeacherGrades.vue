<template>
  <div class="teacher-grades-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>成绩管理</span>
      </div>
      
      <div class="filter-container">
        <el-form :inline="true" class="demo-form-inline">
          <el-form-item label="选择课程">
            <el-select v-model="selectedCourse" placeholder="请选择课程" @change="fetchCourseGrades">
              <el-option
                v-for="course in teacherCourses"
                :key="course.id"
                :label="`${course.course.name} (${course.semester.name})`"
                :value="course.id">
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      
      <el-table
        v-loading="loading"
        :data="courseGrades"
        style="width: 100%"
        border>
        <el-table-column
          prop="student.studentId"
          label="学号"
          width="120">
        </el-table-column>
        <el-table-column
          prop="student.user.name"
          label="学生姓名"
          width="120">
        </el-table-column>
        <el-table-column
          prop="student.class.name"
          label="班级"
          width="120">
        </el-table-column>
        <el-table-column
          prop="grade"
          label="成绩"
          width="120">
          <template slot-scope="scope">
            <el-input-number 
              v-if="editingGrade === scope.row.id"
              v-model="editedGrade"
              :min="0"
              :max="100"
              :precision="1"
              size="mini"
              style="width: 100px;">
            </el-input-number>
            <span v-else :class="{ 'grade-fail': scope.row.grade < 60 }">
              {{ scope.row.grade || '未录入' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          label="状态"
          width="100">
          <template slot-scope="scope">
            <el-tag :type="statusType(scope.row.status)">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          label="操作">
          <template slot-scope="scope">
            <el-button
              v-if="editingGrade !== scope.row.id"
              size="mini"
              type="primary"
              @click="startEditGrade(scope.row)">
              {{ scope.row.grade ? '修改成绩' : '录入成绩' }}
            </el-button>
            <el-button
              v-if="editingGrade === scope.row.id"
              size="mini"
              type="success"
              @click="saveGrade(scope.row)">
              保存
            </el-button>
            <el-button
              v-if="editingGrade === scope.row.id"
              size="mini"
              @click="cancelEdit">
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="summary-info" v-if="courseGrades.length > 0">
        <p>班级平均分: <span class="grade-highlight">{{ averageGrade }}</span></p>
        <p>最高分: <span class="grade-highlight">{{ highestGrade }}</span></p>
        <p>最低分: <span class="grade-highlight">{{ lowestGrade }}</span></p>
        <p>及格率: <span class="grade-highlight">{{ passRate }}%</span></p>
      </div>
      
      <div class="empty-data" v-if="courseGrades.length === 0 && !loading && selectedCourse">
        <el-empty description="暂无学生选课数据"></el-empty>
      </div>
      
      <div class="course-select-tip" v-if="!selectedCourse && !loading">
        <el-alert
          title="请先选择一门课程"
          type="info"
          center
          show-icon>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script>
import gradeApi from '@/api/grade'

export default {
  name: 'TeacherGrades',
  data() {
    return {
      teacherCourses: [
        { id: 1, course: { name: '数据结构' }, semester: { name: '2023-2024学年第一学期' } },
        { id: 2, course: { name: '计算机网络' }, semester: { name: '2023-2024学年第一学期' } },
        { id: 3, course: { name: '高等数学' }, semester: { name: '2023-2024学年第一学期' } },
        { id: 4, course: { name: '操作系统' }, semester: { name: '2023-2024学年第二学期' } },
        { id: 5, course: { name: '软件工程' }, semester: { name: '2023-2024学年第二学期' } }
      ],
      selectedCourse: null,
      courseGrades: [],
      loading: false,
      editingGrade: null,
      editedGrade: 0
    }
  },
  computed: {
    averageGrade() {
      if (this.courseGrades.length === 0) return 0
      const validGrades = this.courseGrades.filter(g => g.grade)
      if (validGrades.length === 0) return 0
      const sum = validGrades.reduce((acc, curr) => acc + curr.grade, 0)
      return (sum / validGrades.length).toFixed(2)
    },
    highestGrade() {
      if (this.courseGrades.length === 0) return 0
      const validGrades = this.courseGrades.filter(g => g.grade)
      if (validGrades.length === 0) return 0
      return Math.max(...validGrades.map(g => g.grade))
    },
    lowestGrade() {
      if (this.courseGrades.length === 0) return 0
      const validGrades = this.courseGrades.filter(g => g.grade)
      if (validGrades.length === 0) return 0
      return Math.min(...validGrades.map(g => g.grade))
    },
    passRate() {
      if (this.courseGrades.length === 0) return 0
      const validGrades = this.courseGrades.filter(g => g.grade)
      if (validGrades.length === 0) return 0
      const passCount = validGrades.filter(g => g.grade >= 60).length
      return Math.round((passCount / validGrades.length) * 100)
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
    async fetchCourseGrades() {
      if (!this.selectedCourse) return
      
      this.loading = true
      try {
        const response = await gradeApi.getCourseGrades(this.selectedCourse)
        this.courseGrades = response.data.data
      } catch (error) {
        this.$message.error('获取课程成绩信息失败')
        console.error('获取课程成绩失败:', error)
        
        // 模拟数据供演示使用
        if (this.selectedCourse === 1) {
          this.courseGrades = [
            { id: 1, student: { studentId: '2023001', user: { name: '张三' }, class: { name: '计科1班' } }, grade: 85.5, status: '已完成' },
            { id: 3, student: { studentId: '2023002', user: { name: '李四' }, class: { name: '计科1班' } }, grade: 78.5, status: '已完成' },
            { id: 4, student: { studentId: '2023003', user: { name: '王五' }, class: { name: '计科2班' } }, grade: 92.0, status: '已完成' },
            { id: 5, student: { studentId: '2023004', user: { name: '赵六' }, class: { name: '计科2班' } }, grade: 56.5, status: '未通过' }
          ]
        } else if (this.selectedCourse === 2) {
          this.courseGrades = [
            { id: 2, student: { studentId: '2023001', user: { name: '张三' }, class: { name: '计科1班' } }, grade: 92.0, status: '已完成' },
            { id: 6, student: { studentId: '2023003', user: { name: '王五' }, class: { name: '计科2班' } }, grade: 84.5, status: '已完成' }
          ]
        } else if (this.selectedCourse === 3) {
          this.courseGrades = [
            { id: 7, student: { studentId: '2023001', user: { name: '张三' }, class: { name: '计科1班' } }, grade: 95.5, status: '已完成' },
            { id: 8, student: { studentId: '2023002', user: { name: '李四' }, class: { name: '计科1班' } }, grade: 88.0, status: '已完成' },
            { id: 9, student: { studentId: '2023005', user: { name: '孙七' }, class: { name: '数学1班' } }, grade: 91.5, status: '已完成' },
            { id: 10, student: { studentId: '2023006', user: { name: '周八' }, class: { name: '数学1班' } }, grade: 76.0, status: '已完成' }
          ]
        } else if (this.selectedCourse === 4) {
          this.courseGrades = [
            { id: 11, student: { studentId: '2023003', user: { name: '王五' }, class: { name: '计科2班' } }, grade: 81.0, status: '已完成' },
            { id: 12, student: { studentId: '2023004', user: { name: '赵六' }, class: { name: '计科2班' } }, grade: 70.5, status: '已完成' },
            { id: 13, student: { studentId: '2023007', user: { name: '吴九' }, class: { name: '软工1班' } }, grade: 88.5, status: '已完成' },
            { id: 14, student: { studentId: '2023008', user: { name: '郑十' }, class: { name: '软工1班' } }, grade: 59.0, status: '未通过' }
          ]
        } else if (this.selectedCourse === 5) {
          this.courseGrades = [
            { id: 15, student: { studentId: '2023007', user: { name: '吴九' }, class: { name: '软工1班' } }, grade: 93.0, status: '已完成' },
            { id: 16, student: { studentId: '2023008', user: { name: '郑十' }, class: { name: '软工1班' } }, grade: 89.5, status: '已完成' },
            { id: 17, student: { studentId: '2023009', user: { name: '冯十一' }, class: { name: '软工2班' } }, grade: 77.0, status: '已完成' },
            { id: 18, student: { studentId: '2023010', user: { name: '陈十二' }, class: { name: '软工2班' } }, grade: 65.5, status: '已完成' }
          ]
        }
      } finally {
        this.loading = false
      }
    },
    startEditGrade(grade) {
      this.editingGrade = grade.id
      this.editedGrade = grade.grade || 0
    },
    cancelEdit() {
      this.editingGrade = null
      this.editedGrade = 0
    },
    async saveGrade(grade) {
      try {
        await gradeApi.updateGrade(grade.id, this.editedGrade)
        this.$message.success('成绩保存成功')
        
        // Update the local data
        const index = this.courseGrades.findIndex(g => g.id === grade.id)
        if (index !== -1) {
          this.courseGrades[index].grade = this.editedGrade
          this.courseGrades[index].status = this.editedGrade >= 60 ? '已完成' : '未通过'
        }
        
        this.editingGrade = null
      } catch (error) {
        this.$message.error('保存成绩失败')
        console.error('保存成绩失败:', error)
      }
    }
  },
  created() {
    // In a real app, this would fetch the teacher's courses
    // this.fetchTeacherCourses()
  }
}
</script>

<style scoped>
.teacher-grades-container {
  padding: 20px;
}
.filter-container {
  margin-bottom: 20px;
}
.summary-info {
  margin-top: 20px;
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  display: flex;
  justify-content: space-around;
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
.empty-data, .course-select-tip {
  margin: 30px 0;
}
</style> 