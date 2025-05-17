import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

// Set axios baseURL
axios.defaults.baseURL = 'http://localhost:8080/api'

// Add token to each request if available
axios.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

export default new Vuex.Store({
  state: {
    token: localStorage.getItem('token') || '',
    user: null,
    departments: [],
    courses: []
  },
  getters: {
    isAuthenticated: state => !!state.token,
    userRole: state => {
      // 尝试从state.user中获取角色
      if (state.user && state.user.role) {
        return state.user.role;
      }
      
      // 如果没有，尝试从localStorage获取
      const savedRole = localStorage.getItem('userRole');
      if (savedRole) {
        return savedRole;
      }
      
      // 默认返回普通用户角色
      return null;
    },
    departments: state => state.departments,
    courses: state => state.courses
  },
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
    },
    SET_USER(state, user) {
      state.user = user
    },
    CLEAR_AUTH(state) {
      state.token = ''
      state.user = null
    },
    SET_DEPARTMENTS(state, departments) {
      state.departments = departments
    },
    SET_COURSES(state, courses) {
      state.courses = courses
    }
  },
  actions: {
    // Authentication actions
    async login({ commit }, credentials) {
      try {
        const response = await axios.post('/login', credentials)
        const { token, user } = response.data
        
        console.log('登录成功，用户信息:', user);
        console.log('用户角色:', user.role);
        
        // Save token to localStorage
        localStorage.setItem('token', token)
        localStorage.setItem('userRole', user.role) // 额外保存角色信息
        
        // Set token in axios headers
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        
        commit('SET_TOKEN', token)
        commit('SET_USER', user)
        
        return Promise.resolve(user)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    logout({ commit }) {
      // Remove token from localStorage
      localStorage.removeItem('token')
      
      // Remove token from axios headers
      delete axios.defaults.headers.common['Authorization']
      
      commit('CLEAR_AUTH')
    },
    
    // Check user authentication status
    async checkAuth({ commit, state }) {
      if (!state.token) {
        return Promise.reject(new Error('No token found'))
      }
      
      try {
        const response = await axios.get('/user')
        console.log('获取用户信息成功:', response.data);
        
        // 如果后端没有返回角色信息，尝试从localStorage获取
        if (response.data && !response.data.role) {
          response.data.role = localStorage.getItem('userRole') || 'unknown';
        }
        
        commit('SET_USER', response.data)
        return Promise.resolve(response.data)
      } catch (error) {
        commit('CLEAR_AUTH')
        localStorage.removeItem('token')
        localStorage.removeItem('userRole')
        return Promise.reject(error)
      }
    },
    
    // Department actions
    async fetchDepartments({ commit }) {
      try {
        const response = await axios.get('/departments')
        commit('SET_DEPARTMENTS', response.data)
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async createDepartment({ dispatch }, department) {
      try {
        const response = await axios.post('/departments', department)
        // Refresh department list
        dispatch('fetchDepartments')
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async updateDepartment({ dispatch }, { id, department }) {
      try {
        const response = await axios.put(`/departments/${id}`, department)
        // Refresh department list
        dispatch('fetchDepartments')
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async deleteDepartment({ dispatch }, id) {
      try {
        const response = await axios.delete(`/departments/${id}`)
        // Refresh department list
        dispatch('fetchDepartments')
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },

    // Course actions
    async fetchCourses({ commit }) {
      try {
        const response = await axios.get('/courses')
        commit('SET_COURSES', response.data)
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async getCourse({ commit }, id) {
      try {
        const response = await axios.get(`/courses/${id}`)
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async createCourse({ dispatch }, course) {
      try {
        const response = await axios.post('/courses', course)
        // Refresh course list
        dispatch('fetchCourses')
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async updateCourse({ dispatch }, { id, course }) {
      try {
        const response = await axios.put(`/courses/${id}`, course)
        // Refresh course list
        dispatch('fetchCourses')
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    
    async deleteCourse({ dispatch }, id) {
      try {
        const response = await axios.delete(`/courses/${id}`)
        // Refresh course list
        dispatch('fetchCourses')
        return Promise.resolve(response.data)
      } catch (error) {
        return Promise.reject(error)
      }
    }
  }
}) 