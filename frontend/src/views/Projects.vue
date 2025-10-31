<template>
  <div class="page-container">
    <ProjectStatsGrid
      :total-count="projects.length"
      :paid-count="paidProjectsCount"
      :free-count="freeProjectsCount"
    />
    
    <el-card class="main-card card-modern">
      <div class="header-actions action-buttons">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建项目
        </el-button>
        <el-button type="success" @click="handleBatchCreate">
          批量创建
        </el-button>
        <el-button type="danger" @click="handleBatchDelete" :disabled="selectedProjects.length === 0">
          批量删除
        </el-button>
      </div>
      
      <ProjectTable
        :projects="projects"
        :loading="loading"
        :page="page"
        :page-size="pageSize"
        :total="total"
        @selection-change="handleSelectionChange"
        @page-change="handlePageChange"
        @edit="handleEdit"
        @view-cards="handleViewCards"
        @view-vars="handleViewVars"
        @delete="handleDelete"
      />
    </el-card>
    
    <ProjectFormDialog
      v-model:visible="dialogVisible"
      :title="dialogTitle"
      :project-data="currentProject"
      @save="handleSave"
    />
    
    <ProjectBatchCreateDialog
      v-model:visible="batchCreateDialogVisible"
      @save="handleBatchCreateSave"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { getProjects, createProject, updateProject, deleteProject, batchCreateProjects, batchDeleteProjects } from '@/api/project'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useTableSelection } from '@/composables/useTableSelection'
import ProjectStatsGrid from '@/components/projects/ProjectStatsGrid.vue'
import ProjectTable from '@/components/projects/ProjectTable.vue'
import ProjectFormDialog from '@/components/projects/ProjectFormDialog.vue'
import ProjectBatchCreateDialog from '@/components/projects/ProjectBatchCreateDialog.vue'

const router = useRouter()
const loading = ref(false)
const projects = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const dialogVisible = ref(false)
const batchCreateDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentId = ref(null)
const currentProject = ref(null)

const { selectedItems: selectedProjects, handleSelectionChange } = useTableSelection()

const paidProjectsCount = computed(() => {
  return projects.value.filter(p => p.mode === 'paid').length
})

const freeProjectsCount = computed(() => {
  return projects.value.filter(p => p.mode === 'free').length
})

const loadProjects = async () => {
  loading.value = true
  try {
    const res = await getProjects({
      page: page.value,
      page_size: pageSize.value
    })
    projects.value = res.list
    total.value = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (newPage) => {
  page.value = newPage
  selectedProjects.value = []
  loadProjects()
}

const handleCreate = () => {
  isEdit.value = false
  dialogTitle.value = '创建项目'
  currentProject.value = null
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑项目'
  currentId.value = row.id
  currentProject.value = { ...row }
  dialogVisible.value = true
}

const handleSave = async (formData) => {
  try {
    if (isEdit.value) {
      await updateProject(currentId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await createProject(formData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadProjects()
  } catch (error) {
    console.error(error)
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该项目吗?', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteProject(row.id)
      ElMessage.success('删除成功')
      loadProjects()
    } catch (error) {
      console.error(error)
    }
  }).catch(() => {})
}

const handleViewCards = (row) => {
  router.push(`/cards/${row.id}`)
}

const handleViewVars = (row) => {
  router.push(`/cloudvars/${row.id}`)
}

const handleBatchCreate = () => {
  batchCreateDialogVisible.value = true
}

const handleBatchCreateSave = async (jsonData) => {
  try {
    const data = JSON.parse(jsonData)
    
    if (!Array.isArray(data)) {
      ElMessage.error('数据格式错误,必须是数组')
      return
    }
    
    const projectsData = data.map(item => ({
      name: item.name,
      mode: item.mode || 'free',
      version: item.version || '1.0.0',
      token_expire: item.token_expire || 3600,
      enable_hwid: item.enable_hwid !== false,
      enable_ip: item.enable_ip !== false,
      description: item.description || ''
    }))
    
    await batchCreateProjects({ data: projectsData })
    ElMessage.success(`成功创建 ${projectsData.length} 个项目`)
    batchCreateDialogVisible.value = false
    loadProjects()
  } catch (error) {
    if (error instanceof SyntaxError) {
      ElMessage.error('JSON格式错误,请检查')
    } else {
      console.error(error)
    }
  }
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确定要删除选中的 ${selectedProjects.value.length} 个项目吗?`, '批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await batchDeleteProjects({
        ids: selectedProjects.value.map(p => p.id)
      })
      ElMessage.success(`成功删除 ${selectedProjects.value.length} 个项目`)
      selectedProjects.value = []
      loadProjects()
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

.main-card {
  animation: slide-in-up var(--duration-normal) var(--ease-out);
}

.header-actions {
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

/* 移动端适配 */
@media (max-width: 768px) {
  .header-actions {
    margin-bottom: 16px;
  }
}
</style>
