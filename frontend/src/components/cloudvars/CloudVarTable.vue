<template>
  <div>
    <!-- 桌面端表格视图 -->
    <el-table 
      v-if="!isMobile"
      :data="cloudVars" 
      style="width: 100%;" 
      v-loading="loading" 
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="key" label="变量名" min-width="200" show-overflow-tooltip />
      <el-table-column prop="value" label="值" min-width="250" show-overflow-tooltip />
      <el-table-column label="操作" width="150" :fixed="isDesktop ? 'right' : false" class-name="action-column">
        <template #default="{ row }">
          <ActionButtons :actions="getRowActions(row)" />
        </template>
      </el-table-column>
    </el-table>
    
    <!-- 移动端卡片视图 -->
    <CloudVarCardList
      v-else
      :cloud-vars="cloudVars"
      :loading="loading"
      :selected-vars="selectedVars"
      :page="page"
      :page-size="pageSize"
      :total="total"
      @selection-change="handleSelectionChange"
      @page-change="(newPage) => $emit('page-change', newPage)"
      @edit="(row) => $emit('edit', row)"
      @delete="(row) => $emit('delete', row)"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Edit, Delete } from '@element-plus/icons-vue'
import ActionButtons from '@/components/common/ActionButtons.vue'
import CloudVarCardList from './CloudVarCardList.vue'
import { useResponsive } from '@/composables/useResponsive'

const { isMobile, isDesktop } = useResponsive()

const selectedVars = ref([])

defineProps({
  cloudVars: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 20
  },
  total: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['selection-change', 'page-change', 'edit', 'delete'])

const handleSelectionChange = (selection) => {
  selectedVars.value = selection
  emit('selection-change', selection)
}

const getRowActions = (row) => [
  {
    key: 'edit',
    icon: Edit,
    label: '编辑',
    handler: () => emit('edit', row)
  },
  {
    key: 'delete',
    icon: Delete,
    label: '删除',
    type: 'danger',
    divided: true,
    handler: () => emit('delete', row)
  }
]
</script>

<style scoped>
:deep(.el-table) {
  border-radius: var(--radius-md);
  overflow: hidden;
}

:deep(.el-table th) {
  background: var(--color-bg-tertiary);
  color: var(--color-text-primary);
  font-weight: 600;
}

:deep(.el-table tr) {
  transition: all var(--duration-fast) var(--ease-out);
}

:deep(.el-table tr:hover) {
  background: rgba(102, 126, 234, 0.04);
  transform: scale(1.001);
}

/* 平板适配 */
@media (min-width: 769px) and (max-width: 1023px) {
  :deep(.el-table__cell) {
    padding: 10px 8px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  :deep(.el-table__cell) {
    padding: 8px 4px;
  }
  
  :deep(.action-column) {
    width: 100px !important;
  }
}
</style>

