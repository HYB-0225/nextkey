<template>
  <div>
    <el-table :data="cards" style="width: 100%;" v-loading="loading" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="card_key" label="卡密" min-width="200" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.activated ? 'success' : 'info'">
            {{ row.activated ? '已激活' : '未激活' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="有效时长" width="120">
        <template #default="{ row }">
          {{ formatDuration(row.duration) }}
        </template>
      </el-table-column>
      <el-table-column prop="card_type" label="类型" width="100" />
      <el-table-column prop="note" label="备注" min-width="120" show-overflow-tooltip />
      <el-table-column label="设备限制" width="120">
        <template #default="{ row }">
          {{ row.hwid_list?.length || 0 }} / {{ row.max_hwid === -1 ? '∞' : row.max_hwid }}
        </template>
      </el-table-column>
      <el-table-column label="IP限制" width="120">
        <template #default="{ row }">
          {{ row.ip_list?.length || 0 }} / {{ row.max_ip === -1 ? '∞' : row.max_ip }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250" :fixed="isDesktop ? 'right' : false" class-name="action-column">
        <template #default="{ row }">
          <ActionButtons :actions="getRowActions(row)" />
        </template>
      </el-table-column>
    </el-table>
    
    <el-pagination
      v-if="total > 0"
      style="margin-top: 20px; justify-content: flex-end;"
      :current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="handlePageChange"
    />
  </div>
</template>

<script setup>
import { Edit, View, Delete } from '@element-plus/icons-vue'
import { formatDuration } from '@/composables/useDuration'
import ActionButtons from '@/components/common/ActionButtons.vue'
import { useResponsive } from '@/composables/useResponsive'

const { isDesktop } = useResponsive()

const props = defineProps({
  cards: {
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

const emit = defineEmits(['selection-change', 'page-change', 'edit', 'view', 'delete'])

const handleSelectionChange = (selection) => {
  emit('selection-change', selection)
}

const handlePageChange = (newPage) => {
  emit('page-change', newPage)
}

const getRowActions = (row) => [
  {
    key: 'edit',
    icon: Edit,
    label: '编辑',
    handler: () => emit('edit', row)
  },
  {
    key: 'view',
    icon: View,
    label: '详情',
    handler: () => emit('view', row)
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

:deep(.el-table tr:hover) {
  background: rgba(102, 126, 234, 0.04);
}

:deep(.el-tag) {
  border-radius: var(--radius-sm);
  font-weight: 500;
  padding: 4px 12px;
}

:deep(.el-tag.el-tag--success) {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  border: none;
  color: #fff;
}

:deep(.el-tag.el-tag--info) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: #fff;
}

:deep(.el-pagination) {
  display: flex;
  justify-content: center;
}

:deep(.el-pagination button),
:deep(.el-pager li) {
  border-radius: var(--radius-sm);
  transition: all var(--duration-fast) var(--ease-out);
}

:deep(.el-pager li.is-active) {
  background: var(--color-primary);
}

/* 平板适配 */
@media (min-width: 769px) and (max-width: 1023px) {
  :deep(.el-table__cell) {
    padding: 10px 8px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  :deep(.el-table) {
    font-size: 14px;
  }
  
  :deep(.el-table__cell) {
    padding: 8px 4px;
  }
  
  :deep(.action-column) {
    width: 100px !important;
  }
}

/* 小屏手机适配 */
@media (max-width: 480px) {
  :deep(.el-table__header-wrapper),
  :deep(.el-table__body-wrapper) {
    overflow-x: auto;
  }
  
  :deep(.el-table__header th:nth-child(5)),
  :deep(.el-table__body td:nth-child(5)),
  :deep(.el-table__header th:nth-child(6)),
  :deep(.el-table__body td:nth-child(6)),
  :deep(.el-table__header th:nth-child(7)),
  :deep(.el-table__body td:nth-child(7)) {
    display: none;
  }
}
</style>

