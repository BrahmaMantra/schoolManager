<template>
  <div class="department-container">
    <el-card>
      <div slot="header" class="clearfix">
        <span>院系管理</span>
        <el-button 
          style="float: right; padding: 3px 0" 
          type="text" 
          icon="el-icon-plus"
          @click="handleAdd"
        >
          新增院系
        </el-button>
      </div>
      
      <el-table
        v-loading="loading"
        :data="departments"
        border
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="院系名称"></el-table-column>
        <el-table-column prop="code" label="院系代码" width="120"></el-table-column>
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

    <!-- 新增/编辑院系对话框 -->
    <el-dialog
      :title="dialogStatus === 'add' ? '新增院系' : '编辑院系'"
      :visible.sync="dialogVisible"
      width="30%"
    >
      <el-form
        ref="departmentForm"
        :model="departmentForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="院系名称" prop="name">
          <el-input v-model="departmentForm.name"></el-input>
        </el-form-item>
        <el-form-item label="院系代码" prop="code">
          <el-input v-model="departmentForm.code"></el-input>
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
  name: 'DepartmentList',
  data() {
    return {
      loading: false,
      departments: [],
      dialogVisible: false,
      dialogStatus: 'add', // 'add' or 'edit'
      departmentForm: {
        id: undefined,
        name: '',
        code: ''
      },
      rules: {
        name: [
          { required: true, message: '请输入院系名称', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入院系代码', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
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
      this.loading = true
      try {
        await this.$store.dispatch('fetchDepartments')
        this.departments = this.$store.getters.departments
      } catch (error) {
        this.$message.error('获取院系数据失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },

    // 新增院系
    handleAdd() {
      this.dialogStatus = 'add'
      this.departmentForm = {
        id: undefined,
        name: '',
        code: ''
      }
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs['departmentForm'].clearValidate()
      })
    },

    // 编辑院系
    handleEdit(row) {
      this.dialogStatus = 'edit'
      this.departmentForm = Object.assign({}, row)
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs['departmentForm'].clearValidate()
      })
    },

    // 删除院系
    handleDelete(row) {
      this.$confirm('确认删除该院系吗？删除后不可恢复', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async () => {
          try {
            await this.$store.dispatch('deleteDepartment', row.id)
            this.$message.success('删除成功')
          } catch (error) {
            this.$message.error(error.response?.data?.error || '删除失败')
            console.error(error)
          }
        })
        .catch(() => {})
    },

    // 提交表单
    submitForm() {
      this.$refs['departmentForm'].validate(async valid => {
        if (valid) {
          try {
            if (this.dialogStatus === 'add') {
              await this.$store.dispatch('createDepartment', this.departmentForm)
              this.$message.success('新增院系成功')
            } else {
              await this.$store.dispatch('updateDepartment', {
                id: this.departmentForm.id,
                department: {
                  name: this.departmentForm.name,
                  code: this.departmentForm.code
                }
              })
              this.$message.success('更新院系成功')
            }
            this.dialogVisible = false
          } catch (error) {
            this.$message.error(error.response?.data?.error || '操作失败')
            console.error(error)
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.department-container {
  padding: 20px;
}
</style> 