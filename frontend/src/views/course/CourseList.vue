<template>
  <div class="course-container">
    <el-card>
      <div slot="header" class="clearfix">
        <span>课程管理</span>
        <el-button 
          style="float: right; padding: 3px 0" 
          type="text" 
          icon="el-icon-plus"
          @click="handleAdd"
        >
          新增课程
        </el-button>
      </div>
      
      <el-table
        v-loading="loading"
        :data="courses"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="课程名称"></el-table-column>
        <el-table-column prop="code" label="课程代码" width="120"></el-table-column>
        <el-table-column prop="credits" label="学分" width="80"></el-table-column>
        <el-table-column prop="hours" label="学时" width="80"></el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template slot-scope="scope">
            <el-tag :type="getTypeTag(scope.row.type)">{{ scope.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="department.name" label="所属院系"></el-table-column>
        <el-table-column label="操作" width="200" align="center">
          <template slot-scope="scope">
            <el-button
              size="mini"
              type="primary"
              icon="el-icon-edit"
              @click="handleEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              size="mini"
              type="danger"
              icon="el-icon-delete"
              @click="handleDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑课程对话框 -->
    <el-dialog
      :title="dialogStatus === 'add' ? '新增课程' : '编辑课程'"
      :visible.sync="dialogVisible"
      width="40%"
    >
      <el-form
        ref="courseForm"
        :model="courseForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="courseForm.name"></el-input>
        </el-form-item>
        <el-form-item label="课程代码" prop="code">
          <el-input v-model="courseForm.code"></el-input>
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="学分" prop="credits">
              <el-input-number v-model="courseForm.credits" :min="0" :max="10" :step="0.5"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="学时" prop="hours">
              <el-input-number v-model="courseForm.hours" :min="0" :max="100" :step="2"></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="课程类型" prop="type">
          <el-select v-model="courseForm.type" placeholder="请选择课程类型" style="width: 100%">
            <el-option label="必修课" value="必修课"></el-option>
            <el-option label="选修课" value="选修课"></el-option>
            <el-option label="实践课" value="实践课"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="所属院系" prop="departmentId">
          <el-select v-model="courseForm.departmentId" placeholder="请选择院系" style="width: 100%">
            <el-option
              v-for="item in departments"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input type="textarea" v-model="courseForm.description" rows="3"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'CourseList',
  data() {
    return {
      loading: false,
      courses: [
        {
          id: 1,
          name: '计算机导论',
          code: 'CS101',
          credits: 3,
          hours: 48,
          type: '必修课',
          department: { id: 1, name: '计算机学院' },
          description: '计算机基础入门课程'
        },
        {
          id: 2,
          name: '数据结构',
          code: 'CS201',
          credits: 4,
          hours: 64,
          type: '必修课',
          department: { id: 1, name: '计算机学院' },
          description: '数据结构基础课程'
        },
        {
          id: 3,
          name: '高等数学',
          code: 'MA101',
          credits: 5,
          hours: 80,
          type: '必修课',
          department: { id: 2, name: '数学学院' },
          description: '微积分、线性代数等基础数学知识'
        },
        {
          id: 4,
          name: 'Web前端开发',
          code: 'CS301',
          credits: 2.5,
          hours: 40,
          type: '选修课',
          department: { id: 1, name: '计算机学院' },
          description: 'HTML, CSS, JavaScript基础教学'
        }
      ],
      departments: [],
      dialogVisible: false,
      dialogStatus: 'add', // 'add' or 'edit'
      courseForm: {
        id: undefined,
        name: '',
        code: '',
        credits: 3,
        hours: 48,
        type: '必修课',
        departmentId: undefined,
        description: ''
      },
      rules: {
        name: [
          { required: true, message: '请输入课程名称', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入课程代码', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        credits: [
          { required: true, message: '请输入学分', trigger: 'blur' }
        ],
        hours: [
          { required: true, message: '请输入学时', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择课程类型', trigger: 'change' }
        ],
        departmentId: [
          { required: true, message: '请选择所属院系', trigger: 'change' }
        ]
      }
    }
  },
  created() {
    this.fetchDepartments()
  },
  methods: {
    // 获取所有院系
    async fetchDepartments() {
      try {
        await this.$store.dispatch('fetchDepartments')
        this.departments = this.$store.getters.departments
      } catch (error) {
        this.$message.error('获取院系数据失败')
        console.error(error)
      }
    },

    // 获取课程类型对应的标签类型
    getTypeTag(type) {
      switch (type) {
        case '必修课':
          return 'primary'
        case '选修课':
          return 'success'
        case '实践课':
          return 'warning'
        default:
          return 'info'
      }
    },

    // 新增课程
    handleAdd() {
      this.dialogStatus = 'add'
      this.courseForm = {
        id: undefined,
        name: '',
        code: '',
        credits: 3,
        hours: 48,
        type: '必修课',
        departmentId: undefined,
        description: ''
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs['courseForm'].clearValidate()
      })
    },

    // 编辑课程
    handleEdit(row) {
      this.dialogStatus = 'edit'
      this.courseForm = {
        id: row.id,
        name: row.name,
        code: row.code,
        credits: row.credits,
        hours: row.hours,
        type: row.type,
        departmentId: row.department.id,
        description: row.description
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs['courseForm'].clearValidate()
      })
    },

    // 删除课程
    handleDelete(row) {
      this.$confirm('确认删除该课程吗？删除后不可恢复', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        // 模拟删除成功
        this.courses = this.courses.filter(item => item.id !== row.id)
        this.$message.success('删除成功')
      }).catch(() => {})
    },

    // 提交表单
    submitForm() {
      this.$refs['courseForm'].validate(valid => {
        if (valid) {
          // 模拟提交成功
          if (this.dialogStatus === 'add') {
            // 生成假ID
            const newId = Math.max(...this.courses.map(c => c.id)) + 1
            const department = this.departments.find(d => d.id === this.courseForm.departmentId)
            
            this.courses.push({
              id: newId,
              name: this.courseForm.name,
              code: this.courseForm.code,
              credits: this.courseForm.credits,
              hours: this.courseForm.hours,
              type: this.courseForm.type,
              department: {
                id: department.id,
                name: department.name
              },
              description: this.courseForm.description
            })
            
            this.$message.success('新增课程成功')
          } else {
            // 查找并更新
            const index = this.courses.findIndex(c => c.id === this.courseForm.id)
            const department = this.departments.find(d => d.id === this.courseForm.departmentId)
            
            if (index !== -1) {
              this.courses[index] = {
                ...this.courses[index],
                name: this.courseForm.name,
                code: this.courseForm.code,
                credits: this.courseForm.credits,
                hours: this.courseForm.hours,
                type: this.courseForm.type,
                department: {
                  id: department.id,
                  name: department.name
                },
                description: this.courseForm.description
              }
            }
            
            this.$message.success('更新课程成功')
          }
          
          this.dialogVisible = false
        }
      })
    }
  }
}
</script>

<style scoped>
.course-container {
  padding: 20px;
}
</style> 