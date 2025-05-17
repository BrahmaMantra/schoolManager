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
    departments: []
  },
  getters: {
    isAuthenticated: state => !!state.token,
    userRole: state => state.user ? state.user.role : null,
    departments: state => state.departments
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
    }
  },
  actions: {
    // Authentication actions
    async login({ commit }, credentials) {
      try {
        const response = await axios.post('/login', credentials)
        const { token, user } = response.data
        
        // Save token to localStorage
        localStorage.setItem('token', token)
        
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
        commit('SET_USER', response.data)
        return Promise.resolve(response.data)
      } catch (error) {
        commit('CLEAR_AUTH')
        localStorage.removeItem('token')
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
    }
  }
}) 