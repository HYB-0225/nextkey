<template>
  <div class="action-buttons">
    <!-- 桌面端显示所有按钮 -->
    <div class="desktop-actions">
      <el-button
        v-for="action in actions"
        :key="action.key"
        size="small"
        :type="action.type || 'default'"
        @click="action.handler"
        class="action-btn"
      >
        <el-icon v-if="action.icon"><component :is="action.icon" /></el-icon>
      </el-button>
    </div>
    
    <!-- 移动端显示下拉菜单 -->
    <el-dropdown class="mobile-actions" trigger="click">
      <el-button size="small" type="primary">
        操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item
            v-for="(action, index) in actions"
            :key="action.key"
            :divided="action.divided"
            @click="action.handler"
          >
            <el-icon v-if="action.icon"><component :is="action.icon" /></el-icon>
            {{ action.label }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup>
import { ArrowDown } from '@element-plus/icons-vue'

defineProps({
  actions: {
    type: Array,
    required: true
  }
})
</script>

<style scoped>
.action-buttons {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.desktop-actions {
  display: flex;
  gap: 4px;
}

.mobile-actions {
  display: none;
}

.action-btn :deep(.el-icon) {
  margin: 0;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-dropdown-menu__item .el-icon) {
  font-size: 16px;
}

@media (max-width: 768px) {
  .desktop-actions {
    display: none;
  }
  
  .mobile-actions {
    display: block;
  }
}
</style>

