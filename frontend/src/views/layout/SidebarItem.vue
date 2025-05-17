<template>
  <div>
    <el-menu-item v-if="!item.children" :index="resolvePath(item.path)">
      <i v-if="item.meta && item.meta.icon" :class="item.meta.icon"></i>
      <span slot="title">{{ item.meta && item.meta.title ? item.meta.title : item.name }}</span>
    </el-menu-item>
    
    <el-submenu v-else ref="submenu" :index="resolvePath(item.path)">
      <template slot="title">
        <i v-if="item.meta && item.meta.icon" :class="item.meta.icon"></i>
        <span slot="title">{{ item.meta && item.meta.title ? item.meta.title : item.name }}</span>
      </template>
      
      <sidebar-item
        v-for="child in item.children"
        :key="child.path"
        :item="child"
        :base-path="resolvePath(child.path)"
      />
    </el-submenu>
  </div>
</template>

<script>
export default {
  name: 'SidebarItem',
  props: {
    item: {
      type: Object,
      required: true
    },
    basePath: {
      type: String,
      default: ''
    }
  },
  methods: {
    resolvePath(routePath) {
      if (this.isExternalLink(routePath)) {
        return routePath
      }
      // 简单的路径拼接逻辑，替代Node.js的path.resolve
      if (routePath.charAt(0) === '/') {
        return routePath
      }
      return this.basePath.replace(/\/+$/, '') + '/' + routePath.replace(/^\/+/, '')
    },
    isExternalLink(path) {
      return /^(https?:|mailto:|tel:)/.test(path)
    }
  }
}
</script> 