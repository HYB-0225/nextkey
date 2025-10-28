<template>
  <div class="page-container">
    <el-card>
      <div class="header-actions">
        <el-select v-model="selectedProjectId" placeholder="选择项目" style="width: 300px; margin-right: 10px;" @change="loadCloudVars">
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
      
      <el-table :data="cloudVars" style="width: 100%; margin-top: 20px;" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="key" label="变量名" width="200" />
        <el-table-column prop="value" label="值" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="变量名">
          <el-input v-model="form.key" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="值">
          <el-input v-model="form.value" type="textarea" :rows="5" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
    
    <el-dialog v-model="batchImportDialogVisible" title="批量导入变量" width="700px">
      <el-alert type="info" :closable="false" style="margin-bottom: 15px;">
        请输入JSON格式的变量数据,格式: [{"key": "变量名", "value": "值"}, ...]
      </el-alert>
      <el-input
        v-model="batchImportData"
        type="textarea"
        :rows="12"
        placeholder='[
  {"key": "api_url", "value": "https://api.example.com"},
  {"key": "app_name", "value": "MyApp"},
  {"key": "version", "value": "1.0.0"}
]'
      />
      
      <template #footer>
        <el-button @click="batchImportDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBatchImportSave">确定导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getProjects } from '@/api/project'
import { getCloudVars, setCloudVar, deleteCloudVar, batchSetCloudVars, batchDeleteCloudVars } from '@/api/cloudvar'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const loading = ref(false)
const projects = ref([])
const selectedProjectId = ref(null)
const cloudVars = ref([])
const dialogVisible = ref(false)
const batchImportDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const selectedVars = ref([])
const batchImportData = ref('')

const form = ref({ key: '', value: '' })

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
  form.value = { key: '', value: '' }
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑变量'
  form.value = { ...row }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    await setCloudVar({
      project_id: selectedProjectId.value,
      key: form.value.key,
      value: form.value.value
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

const handleSelectionChange = (selection) => {
  selectedVars.value = selection
}

const handleBatchImport = () => {
  batchImportData.value = ''
  batchImportDialogVisible.value = true
}

const handleBatchImportSave = async () => {
  try {
    const data = JSON.parse(batchImportData.value)
    
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
  max-width: 1400px;
}

.header-actions {
  display: flex;
  justify-content: flex-start;
}
</style>

