<template>
  <el-container class="app-container">
    <el-aside width="220px" class="sidebar-container">
      <div class="logo">
        <span>高校教学管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        unique-opened
        :collapse="isCollapse"
        router
      >
        <sidebar-item
          v-for="route in routes"
          :key="route.path"
          :item="route"
          :base-path="route.path"
        />
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header class="app-header">
        <div class="header-left">
          <i 
            :class="isCollapse ? 'el-icon-s-unfold' : 'el-icon-s-fold'" 
            @click="toggleSideBar"
            class="toggle-button"
          ></i>
          <breadcrumb />
        </div>
        <div class="header-right">
          <span style="margin-right: 15px;">角色: {{ userRole }}</span>
          <el-dropdown trigger="click" @command="handleCommand">
            <span class="el-dropdown-link">
              {{ username }}<i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="profile">个人中心</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="app-main">
        <transition name="fade-transform" mode="out-in">
          <router-view />
        </transition>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import SidebarItem from './SidebarItem'
import Breadcrumb from './Breadcrumb'

export default {
  name: 'Layout',
  components: {
    SidebarItem,
    Breadcrumb
  },
  data() {
    return {
      isCollapse: false
    }
  },
  computed: {
    username() {
      return this.$store.state.user ? this.$store.state.user.name : 'User'
    },
    userRole() {
      return this.$store.getters.userRole || '未知'
    },
    activeMenu() {
      const route = this.$route
      return route.path
    },
    routes() {
      // Filter routes based on user role
      const userRole = this.$store.getters.userRole
      
      // Return routes from the router that have children
      const mainRoute = this.$router.options.routes.find(route => route.path === '/')
      
      if (mainRoute && mainRoute.children) {
        return mainRoute.children.filter(route => {
          // Show route if no roles specified or if user has required role
          if (!route.meta || !route.meta.roles) {
            return true
          }
          return route.meta.roles.includes(userRole)
        })
      }
      
      return []
    }
  },
  methods: {
    toggleSideBar() {
      this.isCollapse = !this.isCollapse
    },
    handleCommand(command) {
      if (command === 'logout') {
        this.$confirm('确认退出登录?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.$store.dispatch('logout')
          this.$router.push('/login')
          this.$message({
            type: 'success',
            message: '已退出登录'
          })
        }).catch(() => {})
      } else if (command === 'profile') {
        this.$router.push('/profile')
      }
    }
  },
  created() {
    // Check authentication status when component is created
    if (this.$store.getters.isAuthenticated && !this.$store.state.user) {
      this.$store.dispatch('checkAuth').catch(() => {
        this.$router.push('/login')
      })
    }
  }
}
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.sidebar-container {
  background-color: #304156;
  transition: width 0.28s;
  height: 100%;
  overflow: hidden;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  overflow: hidden;
}

.sidebar-menu {
  height: calc(100% - 60px);
  border-right: none;
}

.app-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.toggle-button {
  font-size: 20px;
  cursor: pointer;
  margin-right: 15px;
}

.header-right {
  color: #606266;
}

.el-dropdown-link {
  cursor: pointer;
  font-size: 14px;
}

.app-main {
  height: calc(100vh - 60px);
  overflow: auto;
  padding: 20px;
  background-color: #f0f2f5;
}

.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter,
.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style> 