<template>
  <div class="major-container">
    <el-card>
      <div slot="header" class="clearfix">
        <span>专业管理</span>
        <el-button 
          style="float: right; padding: 3px 0" 
          type="text" 
          icon="el-icon-plus"
          @click="handleAdd"
        >
          新增专业
        </el-button>
      </div>
      
      <el-table
        v-loading="loading"
        :data="majors"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="专业名称"></el-table-column>
        <el-table-column prop="code" label="专业代码" width="120"></el-table-column>
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

    <!-- 新增/编辑专业对话框 -->
    <el-dialog
      :title="dialogStatus === 'add' ? '新增专业' : '编辑专业'"
      :visible.sync="dialogVisible"
      width="30%"
    >
      <el-form
        ref="majorForm"
        :model="majorForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="专业名称" prop="name">
          <el-input v-model="majorForm.name"></el-input>
        </el-form-item>
        <el-form-item label="专业代码" prop="code">
          <el-input v-model="majorForm.code"></el-input>
        </el-form-item>
        <el-form-item label="所属院系" prop="departmentId">
          <el-select v-model="majorForm.departmentId" placeholder="请选择院系" style="width: 100%">
            <el-option
              v-for="item in departments"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            ></el-option>
          </el-select>
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
  name: 'MajorList',
  data() {
    return {
      loading: false,
      majors: [
        {
          id: 1,
          name: '计算机科学与技术',
          code: 'CS001',
          department: { id: 1, name: '计算机学院' }
        },
        {
          id: 2,
          name: '软件工程',
          code: 'SE001',
          department: { id: 1, name: '计算机学院' }
        },
        {
          id: 3,
          name: '数学',
          code: 'MA001',
          department: { id: 2, name: '数学学院' }
        }
      ],
      departments: [],
      dialogVisible: false,
      dialogStatus: 'add', // 'add' or 'edit'
      majorForm: {
        id: undefined,
        name: '',
        code: '',
        departmentId: undefined
      },
      rules: {
        name: [
          { required: true, message: '请输入专业名称', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入专业代码', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
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

    // 新增专业
    handleAdd() {
      this.dialogStatus = 'add'
      this.majorForm = {
        id: undefined,
        name: '',
        code: '',
        departmentId: undefined
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs['majorForm'].clearValidate()
      })
    },

    // 编辑专业
    handleEdit(row) {
      this.dialogStatus = 'edit'
      this.majorForm = {
        id: row.id,
        name: row.name,
        code: row.code,
        departmentId: row.department.id
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs['majorForm'].clearValidate()
      })
    },

    // 删除专业
    handleDelete(row) {
      this.$confirm('确认删除该专业吗？删除后不可恢复', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        // 模拟删除成功
        this.majors = this.majors.filter(item => item.id !== row.id)
        this.$message.success('删除成功')
      }).catch(() => {})
    },

    // 提交表单
    submitForm() {
      this.$refs['majorForm'].validate(valid => {
        if (valid) {
          // 模拟提交成功
          if (this.dialogStatus === 'add') {
            // 生成假ID
            const newId = Math.max(...this.majors.map(m => m.id)) + 1
            const department = this.departments.find(d => d.id === this.majorForm.departmentId)
            
            this.majors.push({
              id: newId,
              name: this.majorForm.name,
              code: this.majorForm.code,
              department: {
                id: department.id,
                name: department.name
              }
            })
            
            this.$message.success('新增专业成功')
          } else {
            // 查找并更新
            const index = this.majors.findIndex(m => m.id === this.majorForm.id)
            const department = this.departments.find(d => d.id === this.majorForm.departmentId)
            
            if (index !== -1) {
              this.majors[index] = {
                ...this.majors[index],
                name: this.majorForm.name,
                code: this.majorForm.code,
                department: {
                  id: department.id,
                  name: department.name
                }
              }
            }
            
            this.$message.success('更新专业成功')
          }
          
          this.dialogVisible = false
        }
      })
    }
  }
}
</script>

<style scoped>
.major-container {
  padding: 20px;
}
</style> 