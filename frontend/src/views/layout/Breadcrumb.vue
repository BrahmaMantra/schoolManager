<template>
  <el-breadcrumb class="app-breadcrumb" separator="/">
    <el-breadcrumb-item v-for="(item, index) in levelList" :key="item.path">
      <span 
        v-if="item.redirect === 'noRedirect' || index === levelList.length - 1" 
        class="no-redirect"
      >
        {{ item.meta && item.meta.title ? item.meta.title : item.name }}
      </span>
      <a v-else @click.prevent="handleLink(item)">
        {{ item.meta && item.meta.title ? item.meta.title : item.name }}
      </a>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script>
export default {
  name: 'Breadcrumb',
  data() {
    return {
      levelList: []
    }
  },
  watch: {
    $route: {
      handler(route) {
        this.getBreadcrumb()
      },
      immediate: true
    }
  },
  methods: {
    getBreadcrumb() {
      // Filter routes with meta.title
      let matched = this.$route.matched.filter(item => item.meta && item.meta.title)
      
      // Add dashboard as first item
      const first = matched[0]
      if (first && first.path !== 'dashboard') {
        matched = [{ path: '/dashboard', meta: { title: '首页' } }].concat(matched)
      }
      
      this.levelList = matched
    },
    handleLink(item) {
      const { path } = item
      this.$router.push(path)
    }
  }
}
</script>

<style scoped>
.app-breadcrumb {
  display: inline-block;
  line-height: 50px;
}

.app-breadcrumb .no-redirect {
  color: #97a8be;
  cursor: text;
}
</style> 