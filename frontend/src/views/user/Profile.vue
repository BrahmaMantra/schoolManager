<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>个人信息</span>
            <el-button style="float: right; padding: 3px 0" type="text" @click="handleEdit">编辑</el-button>
          </div>
          <div class="text-center">
            <el-avatar :size="100" icon="el-icon-user-solid"></el-avatar>
            <h3>{{ userInfo.name }}</h3>
            <p>{{ getRoleName(userInfo.role) }}</p>
          </div>
          <div class="user-info">
            <p><span class="label">用户名:</span> {{ userInfo.username }}</p>
            <p><span class="label">邮箱:</span> {{ userInfo.email }}</p>
            <p><span class="label">电话:</span> {{ userInfo.phone }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card>
          <div slot="header" class="clearfix">
            <span>修改密码</span>
          </div>
          <el-form
            ref="passwordForm"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
          >
            <el-form-item label="原密码" prop="oldPassword">
              <el-input v-model="passwordForm.oldPassword" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="passwordForm.newPassword" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirmPassword">
              <el-input v-model="passwordForm.confirmPassword" type="password" show-password></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitPasswordForm">修改密码</el-button>
              <el-button @click="resetPasswordForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
        
        <el-card style="margin-top: 20px">
          <div slot="header" class="clearfix">
            <span>最近登录记录</span>
          </div>
          <el-table :data="loginRecords" style="width: 100%">
            <el-table-column prop="date" label="日期" width="180"></el-table-column>
            <el-table-column prop="ip" label="IP地址"></el-table-column>
            <el-table-column prop="browser" label="浏览器"></el-table-column>
            <el-table-column prop="os" label="操作系统"></el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 编辑个人信息对话框 -->
    <el-dialog title="编辑个人信息" :visible.sync="dialogVisible" width="40%">
      <el-form
        ref="userForm"
        :model="userForm"
        :rules="userRules"
        label-width="100px"
      >
        <el-form-item label="姓名" prop="name">
          <el-input v-model="userForm.name"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email"></el-input>
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="userForm.phone"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitUserForm">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Profile',
  data() {
    // 密码确认校验规则
    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== this.passwordForm.newPassword) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    
    return {
      userInfo: {
        id: 1,
        username: 'admin',
        name: '系统管理员',
        role: 'admin',
        email: 'admin@university.edu',
        phone: '13800138000'
      },
      dialogVisible: false,
      userForm: {
        name: '',
        email: '',
        phone: ''
      },
      userRules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
          { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        phone: [
          { required: true, message: '请输入电话号码', trigger: 'blur' },
          { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
        ]
      },
      passwordForm: {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      },
      passwordRules: {
        oldPassword: [
          { required: true, message: '请输入原密码', trigger: 'blur' }
        ],
        newPassword: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 6, message: '密码长度不能小于6个字符', trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, message: '请再次输入新密码', trigger: 'blur' },
          { validator: validateConfirmPassword, trigger: 'blur' }
        ]
      },
      loginRecords: [
        {
          date: '2023-04-20 10:30:24',
          ip: '192.168.1.100',
          browser: 'Chrome 89.0.4389.82',
          os: 'Windows 10'
        },
        {
          date: '2023-04-19 15:45:12',
          ip: '192.168.1.100',
          browser: 'Chrome 89.0.4389.82',
          os: 'Windows 10'
        },
        {
          date: '2023-04-18 09:12:58',
          ip: '192.168.1.100',
          browser: 'Chrome 89.0.4389.82',
          os: 'Windows 10'
        }
      ]
    }
  },
  created() {
    // 此处可以从vuex中获取用户信息
    // this.userInfo = this.$store.state.user
  },
  methods: {
    // 获取角色名称
    getRoleName(role) {
      const roleMap = {
        admin: '系统管理员',
        academic: '教务处人员',
        department: '院系管理员',
        teacher: '教师',
        student: '学生',
        supervisor: '教学督导',
        finance: '财务人员'
      }
      return roleMap[role] || role
    },
    
    // 编辑个人信息
    handleEdit() {
      this.userForm = {
        name: this.userInfo.name,
        email: this.userInfo.email,
        phone: this.userInfo.phone
      }
      this.dialogVisible = true
    },
    
    // 提交个人信息表单
    submitUserForm() {
      this.$refs.userForm.validate(valid => {
        if (valid) {
          // 此处应该调用API更新用户信息
          // 模拟更新成功
          this.userInfo = {
            ...this.userInfo,
            name: this.userForm.name,
            email: this.userForm.email,
            phone: this.userForm.phone
          }
          
          this.$message.success('个人信息更新成功')
          this.dialogVisible = false
        }
      })
    },
    
    // 提交修改密码表单
    submitPasswordForm() {
      this.$refs.passwordForm.validate(valid => {
        if (valid) {
          // 此处应该调用API修改密码
          // 模拟修改成功
          this.$message.success('密码修改成功')
          this.resetPasswordForm()
        }
      })
    },
    
    // 重置密码表单
    resetPasswordForm() {
      this.$refs.passwordForm.resetFields()
    }
  }
}
</script>

<style scoped>
.profile-container {
  padding: 20px;
}

.text-center {
  text-align: center;
  margin-bottom: 20px;
}

.user-info {
  margin-top: 20px;
}

.user-info p {
  margin: 10px 0;
}

.label {
  font-weight: bold;
  margin-right: 10px;
  display: inline-block;
  width: 70px;
}
</style> 