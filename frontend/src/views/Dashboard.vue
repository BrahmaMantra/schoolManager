<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <div slot="header">
            <span>系统概览</span>
          </div>
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="6">
              <div class="stat-card bg-primary">
                <div class="stat-icon">
                  <i class="el-icon-user"></i>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ stats.studentCount }}</div>
                  <div class="stat-title">学生数量</div>
                </div>
              </div>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <div class="stat-card bg-success">
                <div class="stat-icon">
                  <i class="el-icon-reading"></i>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ stats.courseCount }}</div>
                  <div class="stat-title">开设课程</div>
                </div>
              </div>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <div class="stat-card bg-warning">
                <div class="stat-icon">
                  <i class="el-icon-office-building"></i>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ stats.departmentCount }}</div>
                  <div class="stat-title">院系数量</div>
                </div>
              </div>
            </el-col>
            <el-col :xs="24" :sm="12" :md="6">
              <div class="stat-card bg-danger">
                <div class="stat-icon">
                  <i class="el-icon-s-custom"></i>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ stats.teacherCount }}</div>
                  <div class="stat-title">教师数量</div>
                </div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <div slot="header">
            <span>学生分布</span>
          </div>
          <div class="chart-container">
            <div ref="studentChart" style="width: 100%; height: 300px"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <div slot="header">
            <span>课程分类统计</span>
          </div>
          <div class="chart-container">
            <div ref="courseChart" style="width: 100%; height: 300px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <div slot="header">
            <span>系统公告</span>
          </div>
          <el-table :data="announcements" style="width: 100%">
            <el-table-column prop="title" label="标题" width="180"></el-table-column>
            <el-table-column prop="content" label="内容"></el-table-column>
            <el-table-column prop="date" label="发布时间" width="180"></el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import * as echarts from 'echarts'

export default {
  name: 'Dashboard',
  data() {
    return {
      stats: {
        studentCount: 1200,
        courseCount: 150,
        departmentCount: 8,
        teacherCount: 80
      },
      announcements: [
        {
          title: '系统上线通知',
          content: '高校教学管理系统正式上线，请各位用户及时使用。',
          date: '2023-04-01'
        },
        {
          title: '选课通知',
          content: '2023学年第二学期选课将于下周开始，请学生提前做好准备。',
          date: '2023-04-02'
        },
        {
          title: '系统维护通知',
          content: '系统将于本周六晚上10点至次日凌晨2点进行例行维护，请各位安排好使用时间。',
          date: '2023-04-03'
        }
      ],
      studentChart: null,
      courseChart: null
    }
  },
  mounted() {
    this.initCharts()
  },
  methods: {
    initCharts() {
      // 初始化学生分布图表
      this.studentChart = echarts.init(this.$refs.studentChart)
      const studentOption = {
        title: {
          text: '学生专业分布'
        },
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
          orient: 'vertical',
          right: 10,
          data: ['计算机科学', '电子工程', '数学', '物理', '化学', '其他']
        },
        series: [
          {
            name: '学生人数',
            type: 'pie',
            radius: ['50%', '70%'],
            avoidLabelOverlap: false,
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: '30',
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: false
            },
            data: [
              { value: 300, name: '计算机科学' },
              { value: 250, name: '电子工程' },
              { value: 180, name: '数学' },
              { value: 150, name: '物理' },
              { value: 120, name: '化学' },
              { value: 200, name: '其他' }
            ]
          }
        ]
      }
      this.studentChart.setOption(studentOption)

      // 初始化课程分类图表
      this.courseChart = echarts.init(this.$refs.courseChart)
      const courseOption = {
        title: {
          text: '课程类型分布'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        legend: {
          data: ['必修课', '选修课', '实践课']
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: ['计算机学院', '物理学院', '数学学院', '化学学院', '外语学院']
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '必修课',
            type: 'bar',
            data: [20, 15, 18, 12, 10]
          },
          {
            name: '选修课',
            type: 'bar',
            data: [15, 12, 10, 8, 18]
          },
          {
            name: '实践课',
            type: 'bar',
            data: [8, 6, 4, 10, 5]
          }
        ]
      }
      this.courseChart.setOption(courseOption)

      // 窗口大小改变时重绘图表
      window.addEventListener('resize', () => {
        this.studentChart.resize()
        this.courseChart.resize()
      })
    }
  },
  beforeDestroy() {
    window.removeEventListener('resize', () => {
      this.studentChart.resize()
      this.courseChart.resize()
    })
  }
}
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.stat-card {
  display: flex;
  padding: 20px;
  border-radius: 4px;
  margin-bottom: 20px;
  color: white;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 48px;
  margin-right: 20px;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 5px;
}

.stat-title {
  font-size: 16px;
}

.bg-primary {
  background-color: #409EFF;
}

.bg-success {
  background-color: #67C23A;
}

.bg-warning {
  background-color: #E6A23C;
}

.bg-danger {
  background-color: #F56C6C;
}

.chart-container {
  position: relative;
  width: 100%;
  height: 300px;
}
</style> 