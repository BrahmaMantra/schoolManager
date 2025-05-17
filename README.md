# 高校教学管理系统 (University Teaching Management System)

基于Vue.js + Node.js (前端) 和 Golang + SQLite (后端) 的高校教学管理系统。

## 依赖
- npm
- sqlite
- golang
## 项目结构

```
.
├── frontend/         # 前端Vue.js项目
│   ├── public/       # 静态资源
│   └── src/          # 源代码
│       ├── components/  # 组件
│       ├── views/       # 页面
│       ├── router/      # 路由
│       ├── store/       # 状态管理
│       ├── api/         # API请求
│       ├── utils/       # 工具函数
│       └── assets/      # 图片等资源
│
├── backend/          # 后端Golang项目
│   ├── api/          # API路由定义
│   ├── controllers/  # 控制器
│   ├── models/       # 数据模型
│   ├── db/           # 数据库操作
│   ├── middleware/   # 中间件
│   ├── utils/        # 工具函数
│   └── config/       # 配置文件
│
└── design.md         # 需求文档
```

## 功能模块

1. 基础信息管理 - 院系专业管理、人员信息管理
2. 教学计划管理 - 教学计划制定、学期开课计划
3. 排课管理 - 智能排课、调停课管理
4. 选课管理 - 选课系统、选课结果处理
5. 成绩管理 - 成绩录入、成绩分析与统计
6. 教学质量评估 - 评教系统、教学评估分析
7. 实践教学管理 - 实习管理、毕业论文管理
8. 教材管理
9. 系统管理 - 用户权限管理、系统配置

## 技术栈

- 前端: Vue.js, Element UI, ECharts
- 后端: Golang
- 数据库: SQLite

## 部署与运行

### 前端

```bash
cd frontend
npm install
npm run serve
```

### 后端

```bash
cd backend
go mod tidy
go run main.go
```

## 用户角色

- 系统管理员
- 教务处人员
- 院系管理员
- 教师
- 学生
- 教学督导
- 财务人员 