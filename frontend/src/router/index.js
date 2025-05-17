import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'

Vue.use(VueRouter)

// Layout components
import Layout from '@/views/layout/Layout.vue'
import Login from '@/views/Login.vue'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '首页', icon: 'el-icon-s-home' }
      },
      {
        path: 'department',
        name: 'Department',
        component: () => import('@/views/department/DepartmentList.vue'),
        meta: { title: '院系管理', icon: 'el-icon-office-building', roles: ['admin', 'academic', 'department'] }
      },
      {
        path: 'major',
        name: 'Major',
        component: () => import('@/views/major/MajorList.vue'),
        meta: { title: '专业管理', icon: 'el-icon-collection', roles: ['admin', 'academic', 'department'] }
      },
      {
        path: 'course',
        name: 'Course',
        component: () => import('@/views/course/CourseList.vue'),
        meta: { title: '课程管理', icon: 'el-icon-reading', roles: ['admin', 'academic', 'department', 'teacher'] }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人中心', icon: 'el-icon-user' }
      }
    ]
  },
  {
    path: '*',
    redirect: '/404'
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // This route requires auth, check if logged in
    if (!token) {
      // Not logged in, redirect to login page
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      // Check if route requires specific roles
      const userRole = store.getters.userRole
      console.log('当前用户角色:', userRole)
      console.log('当前路由:', to.path)
      
      // 如果store中没有用户信息，先加载用户信息
      if (!userRole) {
        store.dispatch('checkAuth').then(() => {
          // 重新应用路由逻辑
          router.push(to.path)
        }).catch(() => {
          next({ path: '/login' })
        })
        return
      }
      
      if (to.meta && to.meta.roles) {
        console.log('路由所需角色:', to.meta.roles)
        // 判断用户是否有权限访问该路由
        if (!to.meta.roles.includes(userRole)) {
          console.log('用户没有权限，重定向到首页')
          next({ path: '/dashboard' })
        } else {
          console.log('用户有权限，允许访问')
          next()
        }
      } else {
        // 无需特定角色，允许访问
        next()
      }
    }
  } else {
    // Route does not require auth
    if (token && to.path === '/login') {
      // Already logged in, redirect to dashboard
      next({ path: '/dashboard' })
    } else {
      next()
    }
  }
})

export default router 