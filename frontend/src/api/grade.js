import axios from 'axios'

const apiBase = '/api'

export default {
  /**
   * 获取学生成绩
   * @param {Object} params - 请求参数，可选
   * @returns {Promise} - 包含成绩数据的Promise
   */
  getStudentGrades(params = {}) {
    return axios.get(`${apiBase}/grades/student`, { params })
  },

  /**
   * 获取课程的所有学生成绩（教师用）
   * @param {number} courseOfferingId - 课程开设ID
   * @returns {Promise} - 包含课程成绩数据的Promise
   */
  getCourseGrades(courseOfferingId) {
    return axios.get(`${apiBase}/grades/course`, { 
      params: { course_offering_id: courseOfferingId } 
    })
  },

  /**
   * 更新学生成绩（教师用）
   * @param {number} id - 成绩记录ID
   * @param {number} grade - 新的成绩值
   * @returns {Promise} - 更新结果的Promise
   */
  updateGrade(id, grade) {
    return axios.put(`${apiBase}/grades/${id}`, { grade })
  },

  /**
   * 创建学生成绩记录（教师用）
   * @param {Object} gradeData - 成绩数据
   * @returns {Promise} - 创建结果的Promise
   */
  createGrade(gradeData) {
    return axios.post(`${apiBase}/grades`, gradeData)
  }
} 