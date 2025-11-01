<template>
  <div>
    <!-- 桌面端表格视图 -->
    <template v-if="!isMobile">
      <el-table :data="cards" style="width: 100%;" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="卡密" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <CopyableText :text="row.card_key" success-message="卡密已复制">
              {{ row.card_key }}
            </CopyableText>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'frozen'" type="danger" size="small">
              已冻结
            </el-tag>
            <el-tag v-else-if="row.is_online" type="success" size="small" effect="dark">
              在线
            </el-tag>
            <el-tag v-else-if="row.is_expired && row.status === 'activated'" type="warning" size="small">
              已过期
            </el-tag>
            <el-tag v-else-if="row.status === 'activated'" type="success" size="small">
              已激活
            </el-tag>
            <el-tag v-else type="info" size="small">
              未激活
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="到期时间" width="220">
          <template #default="{ row }">
            {{ formatExpireTime(row) }}
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
        <el-table-column label="操作" width="300" :fixed="isDesktop ? 'right' : false" class-name="action-column">
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
    </template>
    
    <!-- 移动端卡片视图 -->
    <CardListMobile
      v-else
      :cards="cards"
      :loading="loading"
      :page="page"
      :page-size="pageSize"
      :total="total"
      :selected-cards="selectedCards"
      @selection-change="handleSelectionChange"
      @page-change="handlePageChange"
      @edit="(row) => $emit('edit', row)"
      @view="(row) => $emit('view', row)"
      @delete="(row) => $emit('delete', row)"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Edit, View, Delete, Lock, Unlock } from '@element-plus/icons-vue'
import { formatDuration, formatExpireTime } from '@/composables/useDuration'
import ActionButtons from '@/components/common/ActionButtons.vue'
import CopyableText from '@/components/common/CopyableText.vue'
import CardListMobile from './CardListMobile.vue'
import { useResponsive } from '@/composables/useResponsive'

const { isMobile, isDesktop } = useResponsive()

const selectedCards = ref([])

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

const emit = defineEmits(['selection-change', 'page-change', 'edit', 'view', 'delete', 'freeze', 'unfreeze'])

const handleSelectionChange = (selection) => {
  selectedCards.value = selection
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
    key: 'freeze',
    icon: row.frozen ? Unlock : Lock,
    label: row.frozen ? '恢复' : '冻结',
    type: row.frozen ? 'success' : 'warning',
    divided: true,
    handler: () => row.frozen ? emit('unfreeze', row) : emit('freeze', row)
  },
  {
    key: 'delete',
    icon: Delete,
    label: '删除',
    type: 'danger',
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
  background: linear-gradient(135deg, #FFD93D 0%, #FFA400 100%);
  border: none;
  color: #fff;
}

:deep(.el-tag.el-tag--info) {
  background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%);
  border: none;
  color: #fff;
}

:deep(.el-tag.el-tag--danger) {
  background: linear-gradient(135deg, #D2691E 0%, #A0522D 100%);
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

