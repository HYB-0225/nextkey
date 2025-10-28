<template>
  <div class="page-container">
    <el-card>
      <div class="header-actions">
        <el-select 
          v-model="selectedProjectId" 
          placeholder="选择项目" 
          :style="isMobile ? 'width: 100%; order: -1;' : 'width: 300px; margin-right: 10px;'" 
          @change="loadCards"
        >
          <el-option
            v-for="project in projects"
            :key="project.id"
            :label="project.name"
            :value="project.id"
          />
        </el-select>
        <el-button type="primary" @click="handleCreate" :disabled="!selectedProjectId">
          <el-icon><Plus /></el-icon>
          生成卡密
        </el-button>
        <el-button type="warning" @click="handleBatchUpdate" :disabled="selectedCards.length === 0" style="margin-left: 10px;">
          批量修改
        </el-button>
        <el-button type="warning" @click="handleBatchFreeze" :disabled="selectedCards.length === 0">
          批量冻结
        </el-button>
        <el-button type="success" @click="handleBatchUnfreeze" :disabled="selectedCards.length === 0">
          批量恢复
        </el-button>
        <el-button type="success" @click="handleBatchExport" :disabled="selectedCards.length === 0">
          批量导出
        </el-button>
        <el-button type="danger" @click="handleBatchDelete" :disabled="selectedCards.length === 0">
          批量删除
        </el-button>
      </div>
      
      <CardSearchBar
        @search="handleSearch"
        @reset="handleSearchReset"
      />
      
      <CardTable
        :cards="cards"
        :loading="loading"
        :page="page"
        :page-size="pageSize"
        :total="total"
        @selection-change="handleSelectionChange"
        @page-change="handlePageChange"
        @edit="handleEdit"
        @view="handleView"
        @delete="handleDelete"
        @freeze="handleFreeze"
        @unfreeze="handleUnfreeze"
      />
    </el-card>
    
    <CardCreateDialog
      v-model:visible="createDialogVisible"
      @save="handleCreateSave"
    />
    
    <CardEditDialog
      v-model:visible="editDialogVisible"
      :card-data="editCardData"
      @save="handleEditSave"
    />
    
    <CardDetailDialog
      v-model:visible="detailDialogVisible"
      :card="currentCard"
    />
    
    <CardBatchUpdateDialog
      v-model:visible="batchUpdateDialogVisible"
      :selected-count="selectedCards.length"
      @save="handleBatchUpdateSave"
    />
    
    <CardCreatedDialog
      v-model:visible="createdDialogVisible"
      :cards="createdCards"
    />
    
    <el-dialog
      v-model="exportDialogVisible"
      title="选择导出格式"
      :width="isMobile ? '95%' : '400px'"
    >
      <div style="display: flex; flex-direction: column; gap: 10px;">
        <el-button type="primary" @click="handleExportFormat('json')">
          导出为 JSON
        </el-button>
        <el-button type="primary" @click="handleExportFormat('txt')">
          导出为 TXT (仅卡密)
        </el-button>
        <el-button type="primary" @click="handleExportFormat('csv')">
          导出为 CSV
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { getProjects } from '@/api/project'
import { getCards, createCards, deleteCard, updateCard, batchUpdateCards, batchDeleteCards, freezeCard, unfreezeCard, batchFreezeCards, batchUnfreezeCards } from '@/api/card'
import { ElMessage, ElMessageBox } from 'element-plus'
import { usePagination } from '@/composables/usePagination'
import { useTableSelection } from '@/composables/useTableSelection'
import { useResponsive } from '@/composables/useResponsive'
import { secondsToUnitValue, unitValueToSeconds } from '@/composables/useDuration'
import { exportToJSON, exportToTXT, exportToCSV } from '@/utils/export'
import CardTable from '@/components/cards/CardTable.vue'
import CardSearchBar from '@/components/cards/CardSearchBar.vue'
import CardCreateDialog from '@/components/cards/CardCreateDialog.vue'
import CardEditDialog from '@/components/cards/CardEditDialog.vue'
import CardDetailDialog from '@/components/cards/CardDetailDialog.vue'
import CardBatchUpdateDialog from '@/components/cards/CardBatchUpdateDialog.vue'
import CardCreatedDialog from '@/components/cards/CardCreatedDialog.vue'

const route = useRoute()
const { isMobile } = useResponsive()
const loading = ref(false)
const projects = ref([])
const selectedProjectId = ref(null)
const cards = ref([])
const searchParams = ref({})

const { page, pageSize, total, handlePageChange: onPageChange } = usePagination()
const { selectedItems: selectedCards, handleSelectionChange } = useTableSelection()

const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const batchUpdateDialogVisible = ref(false)
const exportDialogVisible = ref(false)
const createdDialogVisible = ref(false)
const currentCard = ref(null)
const editCardData = ref(null)
const createdCards = ref([])

const loadProjects = async () => {
  try {
    projects.value = await getProjects()
    if (route.params.projectId) {
      selectedProjectId.value = parseInt(route.params.projectId)
      loadCards()
    }
  } catch (error) {
    console.error(error)
  }
}

const loadCards = async () => {
  if (!selectedProjectId.value) return
  
  loading.value = true
  try {
    const data = await getCards({
      project_id: selectedProjectId.value,
      page: page.value,
      page_size: pageSize.value,
      ...searchParams.value
    })
    cards.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = (params) => {
  searchParams.value = params
  page.value = 1
  loadCards()
}

const handleSearchReset = () => {
  searchParams.value = {}
  page.value = 1
  loadCards()
}

const handlePageChange = (newPage) => {
  onPageChange(newPage)
  loadCards()
}

const handleCreate = () => {
  createDialogVisible.value = true
}

const handleCreateSave = async (formData) => {
  try {
    const duration = unitValueToSeconds(formData.duration_value, formData.duration_unit)
    
    const params = {
      project_id: selectedProjectId.value,
      duration: duration,
      card_type: formData.card_type,
      max_hwid: formData.max_hwid,
      max_ip: formData.max_ip,
      note: formData.note
    }
    
    if (formData.mode === 'custom') {
      params.card_key = formData.card_key
      params.count = 1
    } else {
      params.prefix = formData.prefix
      params.suffix = formData.suffix
      params.count = formData.count
    }
    
    const result = await createCards(params)
    createDialogVisible.value = false
    
    createdCards.value = result
    createdDialogVisible.value = true
    
    loadCards()
  } catch (error) {
    console.error(error)
  }
}

const handleEdit = (row) => {
  const { value, unit } = secondsToUnitValue(row.duration)
  editCardData.value = {
    id: row.id,
    duration_value: value,
    duration_unit: unit,
    card_type: row.card_type,
    max_hwid: row.max_hwid,
    max_ip: row.max_ip,
    note: row.note,
    custom_data: row.custom_data || '',
    hwid_list: row.hwid_list || [],
    ip_list: row.ip_list || []
  }
  editDialogVisible.value = true
}

const handleEditSave = async (formData) => {
  try {
    const duration = unitValueToSeconds(formData.duration_value, formData.duration_unit)
    
    await updateCard(formData.id, {
      duration: duration,
      card_type: formData.card_type,
      max_hwid: formData.max_hwid,
      max_ip: formData.max_ip,
      note: formData.note,
      custom_data: formData.custom_data,
      hwid_list: formData.hwid_list,
      ip_list: formData.ip_list
    })
    ElMessage.success('保存成功')
    editDialogVisible.value = false
    loadCards()
  } catch (error) {
    console.error(error)
  }
}

const handleView = (row) => {
  currentCard.value = { ...row }
  detailDialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该卡密吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteCard(row.id)
      ElMessage.success('删除成功')
      loadCards()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleBatchUpdate = () => {
  batchUpdateDialogVisible.value = true
}

const handleBatchUpdateSave = async (formData) => {
  try {
    const updateData = {}
    
    if (formData.update_duration) {
      updateData.duration = unitValueToSeconds(formData.duration_value, formData.duration_unit)
    }
    if (formData.update_card_type) {
      updateData.card_type = formData.card_type
    }
    if (formData.update_max_hwid) {
      updateData.max_hwid = formData.max_hwid
    }
    if (formData.update_max_ip) {
      updateData.max_ip = formData.max_ip
    }
    if (formData.update_note) {
      updateData.note = formData.note
    }
    
    if (Object.keys(updateData).length === 0) {
      ElMessage.warning('请至少选择一个要修改的字段')
      return
    }
    
    await batchUpdateCards({
      ids: selectedCards.value.map(c => c.id),
      data: updateData
    })
    
    ElMessage.success(`成功更新 ${selectedCards.value.length} 个卡密`)
    batchUpdateDialogVisible.value = false
    selectedCards.value = []
    loadCards()
  } catch (error) {
    console.error(error)
  }
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确定要删除选中的 ${selectedCards.value.length} 个卡密吗?`, '批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await batchDeleteCards({
        ids: selectedCards.value.map(c => c.id)
      })
      ElMessage.success(`成功删除 ${selectedCards.value.length} 个卡密`)
      selectedCards.value = []
      loadCards()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleFreeze = (row) => {
  ElMessageBox.confirm('确定要冻结该卡密吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await freezeCard(row.id)
      ElMessage.success('冻结成功')
      loadCards()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleUnfreeze = (row) => {
  ElMessageBox.confirm('确定要恢复该卡密吗?', '确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
      await unfreezeCard(row.id)
      ElMessage.success('恢复成功')
      loadCards()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleBatchFreeze = () => {
  ElMessageBox.confirm(`确定要冻结选中的 ${selectedCards.value.length} 个卡密吗?`, '批量冻结', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await batchFreezeCards({
        ids: selectedCards.value.map(c => c.id)
      })
      ElMessage.success(`成功冻结 ${selectedCards.value.length} 个卡密`)
      selectedCards.value = []
      loadCards()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleBatchUnfreeze = () => {
  ElMessageBox.confirm(`确定要恢复选中的 ${selectedCards.value.length} 个卡密吗?`, '批量恢复', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    try {
      await batchUnfreezeCards({
        ids: selectedCards.value.map(c => c.id)
      })
      ElMessage.success(`成功恢复 ${selectedCards.value.length} 个卡密`)
      selectedCards.value = []
      loadCards()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleBatchExport = () => {
  exportDialogVisible.value = true
}

const handleExportFormat = (format) => {
  const timestamp = new Date().toISOString().slice(0, 19).replace(/[:-]/g, '').replace('T', '_')
  const filename = `cards_export_${timestamp}`
  
  switch (format) {
    case 'json':
      exportToJSON(selectedCards.value, filename)
      break
    case 'txt':
      exportToTXT(selectedCards.value, filename)
      break
    case 'csv':
      exportToCSV(selectedCards.value, filename)
      break
  }
  
  exportDialogVisible.value = false
  ElMessage.success(`成功导出 ${selectedCards.value.length} 个卡密`)
}

onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.page-container {
  width: 100%;
}

.header-actions {
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 20px;
}

:deep(.el-card) {
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-sm);
  transition: all var(--duration-normal) var(--ease-out);
}

:deep(.el-card:hover) {
  box-shadow: var(--shadow-md);
}

:deep(.el-select .el-input__wrapper) {
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

:deep(.el-select .el-input__wrapper:hover) {
  border-color: var(--color-primary);
}

:deep(.el-button) {
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--duration-fast) var(--ease-out);
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
}

:deep(.el-button:active) {
  transform: translateY(0);
}

/* 平板适配 */
@media (min-width: 769px) and (max-width: 1023px) {
  .header-actions {
    gap: 10px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .header-actions {
    margin-bottom: 16px;
    gap: 8px;
  }
  
  .header-actions :deep(.el-select) {
    width: 100%;
    order: -1;
  }
  
  .header-actions :deep(.el-button) {
    flex: 1;
  }
}
</style>
