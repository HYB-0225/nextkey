<template>
  <div>
    <!-- 桌面端表格视图 -->
    <el-table 
      v-if="!isMobile" 
      :data="projects" 
      style="width: 100%;" 
      v-loading="loading" 
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="项目名称" min-width="120" />
      <el-table-column label="项目UUID" min-width="280" show-overflow-tooltip>
        <template #default="{ row }">
          <CopyableText :text="row.uuid" success-message="UUID已复制">
            {{ row.uuid }}
          </CopyableText>
        </template>
      </el-table-column>
      <el-table-column prop="mode" label="模式" width="100">
        <template #default="{ row }">
          <el-tag :type="row.mode === 'paid' ? 'success' : 'info'">
            {{ row.mode === 'paid' ? '付费' : '免费' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column label="在线人数" width="120">
        <template #default="{ row }">
          <el-tag type="success" effect="plain">
            {{ row.online_count || 0 }} 人在线
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="加密密钥" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <CopyableText :text="row.encryption_key" success-message="密钥已复制" :masked="true">
            {{ row.encryption_key }}
          </CopyableText>
        </template>
      </el-table-column>
      <el-table-column label="配置" width="180">
        <template #default="{ row }">
          <el-tag v-if="row.enable_hwid" size="small" style="margin-right: 5px">机器码</el-tag>
          <el-tag v-if="row.enable_ip" size="small">IP验证</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250" :fixed="isDesktop ? 'right' : false" class-name="action-column">
        <template #default="{ row }">
          <ActionButtons :actions="getRowActions(row)" />
        </template>
      </el-table-column>
    </el-table>
    
    <!-- 移动端卡片视图 -->
    <ProjectCardList
      v-else
      :projects="projects"
      :loading="loading"
      :selected-projects="selectedProjects"
      :page="page"
      :page-size="pageSize"
      :total="total"
      @selection-change="handleSelectionChange"
      @page-change="(newPage) => $emit('page-change', newPage)"
      @edit="(row) => $emit('edit', row)"
      @view-cards="(row) => $emit('view-cards', row)"
      @view-vars="(row) => $emit('view-vars', row)"
      @unbind-link="(row) => $emit('unbind-link', row)"
      @delete="(row) => $emit('delete', row)"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Edit, Ticket, Cloudy, Link, Delete } from '@element-plus/icons-vue'
import ActionButtons from '@/components/common/ActionButtons.vue'
import CopyableText from '@/components/common/CopyableText.vue'
import ProjectCardList from './ProjectCardList.vue'
import { useResponsive } from '@/composables/useResponsive'

const { isMobile, isDesktop } = useResponsive()

const selectedProjects = ref([])

defineProps({
  projects: {
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

const emit = defineEmits(['selection-change', 'page-change', 'edit', 'view-cards', 'view-vars', 'unbind-link', 'delete'])

const handleSelectionChange = (selection) => {
  selectedProjects.value = selection
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
    key: 'cards',
    icon: Ticket,
    label: '卡密',
    handler: () => emit('view-cards', row)
  },
  {
    key: 'vars',
    icon: Cloudy,
    label: '变量',
    handler: () => emit('view-vars', row)
  },
  {
    key: 'unbind-link',
    icon: Link,
    label: '解绑链接',
    handler: () => emit('unbind-link', row)
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
  
  :deep(.el-table__header th:nth-child(3)),
  :deep(.el-table__body td:nth-child(3)),
  :deep(.el-table__header th:nth-child(5)),
  :deep(.el-table__body td:nth-child(5)),
  :deep(.el-table__header th:nth-child(6)),
  :deep(.el-table__body td:nth-child(6)) {
    display: none;
  }
}
</style>

