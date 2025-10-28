<template>
  <div class="page-container">
    <el-card>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建项目
        </el-button>
        <el-button type="success" @click="handleBatchCreate" style="margin-left: 10px;">
          批量创建
        </el-button>
        <el-button type="danger" @click="handleBatchDelete" :disabled="selectedProjects.length === 0">
          批量删除
        </el-button>
      </div>
      
      <el-table :data="projects" style="width: 100%; margin-top: 20px;" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="项目名称" />
        <el-table-column prop="uuid" label="项目UUID" width="280" />
        <el-table-column prop="mode" label="模式" width="100">
          <template #default="{ row }">
            <el-tag :type="row.mode === 'paid' ? 'success' : 'info'">
              {{ row.mode === 'paid' ? '付费' : '免费' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column label="配置" width="180">
          <template #default="{ row }">
            <el-tag v-if="row.enable_hwid" size="small" style="margin-right: 5px">机器码</el-tag>
            <el-tag v-if="row.enable_ip" size="small">IP验证</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleViewCards(row)">卡密</el-button>
            <el-button size="small" @click="handleViewVars(row)">变量</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form :model="form" label-width="120px">
        <el-form-item label="项目名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="模式">
          <el-radio-group v-model="form.mode">
            <el-radio label="free">免费</el-radio>
            <el-radio label="paid">付费</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="版本号">
          <el-input v-model="form.version" />
        </el-form-item>
        <el-form-item label="Token有效期">
          <el-input v-model.number="form.token_expire" type="number">
            <template #append>秒</template>
          </el-input>
        </el-form-item>
        <el-form-item label="启用机器码">
          <el-switch v-model="form.enable_hwid" />
        </el-form-item>
        <el-form-item label="启用IP验证">
          <el-switch v-model="form.enable_ip" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
    
    <el-dialog v-model="batchCreateDialogVisible" title="批量创建项目" width="700px">
      <el-alert type="info" :closable="false" style="margin-bottom: 15px;">
        请输入JSON格式的项目数据,格式: [{"name": "项目名", "mode": "free", ...}, ...]
      </el-alert>
      <el-input
        v-model="batchCreateData"
        type="textarea"
        :rows="12"
        placeholder='[
  {
    "name": "项目1",
    "mode": "free",
    "version": "1.0.0",
    "token_expire": 3600,
    "enable_hwid": true,
    "enable_ip": true,
    "description": "描述"
  }
]'
      />
      
      <template #footer>
        <el-button @click="batchCreateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBatchCreateSave">确定创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getProjects, createProject, updateProject, deleteProject, batchCreateProjects, batchDeleteProjects } from '@/api/project'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const projects = ref([])
const dialogVisible = ref(false)
const batchCreateDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentId = ref(null)
const selectedProjects = ref([])
const batchCreateData = ref('')

const form = ref({
  name: '',
  mode: 'free',
  version: '1.0.0',
  token_expire: 3600,
  enable_hwid: true,
  enable_ip: true,
  description: ''
})

const loadProjects = async () => {
  loading.value = true
  try {
    projects.value = await getProjects()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  dialogTitle.value = '创建项目'
  form.value = {
    name: '',
    mode: 'free',
    version: '1.0.0',
    token_expire: 3600,
    enable_hwid: true,
    enable_ip: true,
    description: ''
  }
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑项目'
  currentId.value = row.id
  form.value = { ...row }
  dialogVisible.value = true
}

const handleSave = async () => {
  try {
    if (isEdit.value) {
      await updateProject(currentId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await createProject(form.value)
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

const handleSelectionChange = (selection) => {
  selectedProjects.value = selection
}

const handleBatchCreate = () => {
  batchCreateData.value = ''
  batchCreateDialogVisible.value = true
}

const handleBatchCreateSave = async () => {
  try {
    const data = JSON.parse(batchCreateData.value)
    
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
  max-width: 1400px;
}

.header-actions {
  display: flex;
  justify-content: flex-end;
}
</style>

