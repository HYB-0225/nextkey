<template>
  <div class="page-container">
    <el-card>
      <div class="header-actions">
        <el-select 
          v-model="selectedProjectId" 
          placeholder="选择项目" 
          :style="isMobile ? 'width: 100%; order: -1;' : 'width: 300px; margin-right: 10px;'" 
          @change="loadCloudVars"
        >
          <el-option v-for="project in projects" :key="project.id" :label="project.name" :value="project.id" />
        </el-select>
        <el-button type="primary" @click="handleCreate" :disabled="!selectedProjectId">
          <el-icon><Plus /></el-icon>
          添加变量
        </el-button>
        <el-button type="success" @click="handleBatchImport" :disabled="!selectedProjectId" style="margin-left: 10px;">
          批量导入
        </el-button>
        <el-button type="danger" @click="handleBatchDelete" :disabled="selectedVars.length === 0">
          批量删除
        </el-button>
      </div>
      
      <CloudVarTable
        :cloud-vars="cloudVars"
        :loading="loading"
        @selection-change="handleSelectionChange"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </el-card>
    
    <CloudVarFormDialog
      v-model:visible="dialogVisible"
      :title="dialogTitle"
      :is-edit="isEdit"
      :var-data="currentVar"
      @save="handleSave"
    />
    
    <CloudVarBatchImportDialog
      v-model:visible="batchImportDialogVisible"
      @save="handleBatchImportSave"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { getProjects } from '@/api/project'
import { getCloudVars, setCloudVar, deleteCloudVar, batchSetCloudVars, batchDeleteCloudVars } from '@/api/cloudvar'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useTableSelection } from '@/composables/useTableSelection'
import { useResponsive } from '@/composables/useResponsive'
import CloudVarTable from '@/components/cloudvars/CloudVarTable.vue'
import CloudVarFormDialog from '@/components/cloudvars/CloudVarFormDialog.vue'
import CloudVarBatchImportDialog from '@/components/cloudvars/CloudVarBatchImportDialog.vue'

const route = useRoute()
const { isMobile } = useResponsive()
const loading = ref(false)
const projects = ref([])
const selectedProjectId = ref(null)
const cloudVars = ref([])
const dialogVisible = ref(false)
const batchImportDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentVar = ref(null)

const { selectedItems: selectedVars, handleSelectionChange } = useTableSelection()

const loadProjects = async () => {
  try {
    projects.value = await getProjects()
    if (route.params.projectId) {
      selectedProjectId.value = parseInt(route.params.projectId)
      loadCloudVars()
    }
  } catch (error) {
    console.error(error)
  }
}

const loadCloudVars = async () => {
  if (!selectedProjectId.value) return
  loading.value = true
  try {
    cloudVars.value = await getCloudVars({ project_id: selectedProjectId.value })
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  dialogTitle.value = '添加变量'
  currentVar.value = null
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑变量'
  currentVar.value = { ...row }
  dialogVisible.value = true
}

const handleSave = async (formData) => {
  try {
    await setCloudVar({
      project_id: selectedProjectId.value,
      key: formData.key,
      value: formData.value
    })
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadCloudVars()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该变量吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteCloudVar(row.id)
      ElMessage.success('删除成功')
      loadCloudVars()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleBatchImport = () => {
  batchImportDialogVisible.value = true
}

const handleBatchImportSave = async (jsonData) => {
  try {
    const data = JSON.parse(jsonData)
    
    if (!Array.isArray(data)) {
      ElMessage.error('数据格式错误,必须是数组')
      return
    }
    
    const varsData = data.map(item => ({
      project_id: selectedProjectId.value,
      key: item.key,
      value: item.value
    }))
    
    await batchSetCloudVars({ data: varsData })
    ElMessage.success(`成功导入 ${varsData.length} 个变量`)
    batchImportDialogVisible.value = false
    loadCloudVars()
  } catch (error) {
    if (error instanceof SyntaxError) {
      ElMessage.error('JSON格式错误,请检查')
    } else {
      console.error(error)
    }
  }
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确定要删除选中的 ${selectedVars.value.length} 个变量吗?`, '批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await batchDeleteCloudVars({
        ids: selectedVars.value.map(v => v.id)
      })
      ElMessage.success(`成功删除 ${selectedVars.value.length} 个变量`)
      selectedVars.value = []
      loadCloudVars()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
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
  animation: slide-in-up var(--duration-normal) var(--ease-out);
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
