<template>
  <div class="page-container">
    <el-card>
      <div class="header-actions">
        <el-select v-model="selectedProjectId" placeholder="选择项目" style="width: 300px; margin-right: 10px;" @change="loadCards">
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
        <el-button type="danger" @click="handleBatchDelete" :disabled="selectedCards.length === 0">
          批量删除
        </el-button>
      </div>
      
      <el-table :data="cards" style="width: 100%; margin-top: 20px;" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="card_key" label="卡密" width="200" />
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
        <el-table-column prop="note" label="备注" />
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
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
    </el-card>
    
    <el-dialog v-model="createDialogVisible" title="生成卡密" width="600px">
      <el-form :model="createForm" label-width="120px">
        <el-form-item label="生成方式">
          <el-radio-group v-model="createForm.mode">
            <el-radio label="batch">批量生成</el-radio>
            <el-radio label="custom">自定义卡密</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="自定义卡密" v-if="createForm.mode === 'custom'">
          <el-input v-model="createForm.card_key" placeholder="输入完整卡密" />
        </el-form-item>
        
        <template v-if="createForm.mode === 'batch'">
          <el-form-item label="前缀">
            <el-input v-model="createForm.prefix" />
          </el-form-item>
          <el-form-item label="后缀">
            <el-input v-model="createForm.suffix" />
          </el-form-item>
          <el-form-item label="生成数量">
            <el-input-number v-model="createForm.count" :min="1" :max="1000" />
          </el-form-item>
        </template>
        
        <el-form-item label="有效时长">
          <div style="display: flex; gap: 10px; align-items: center;">
            <el-input-number v-model="createForm.duration_value" :min="0" style="width: 150px;" />
            <el-select v-model="createForm.duration_unit" style="width: 100px;">
              <el-option label="秒" value="second" />
              <el-option label="天" value="day" />
              <el-option label="周" value="week" />
              <el-option label="月" value="month" />
              <el-option label="季" value="quarter" />
              <el-option label="年" value="year" />
            </el-select>
          </div>
          <div style="color: #999; font-size: 12px; margin-top: 5px;">0表示永久</div>
        </el-form-item>
        
        <el-form-item label="卡密类型">
          <el-input v-model="createForm.card_type" placeholder="normal" />
        </el-form-item>
        
        <el-form-item label="设备码上限">
          <el-input-number v-model="createForm.max_hwid" :min="-1" />
          <div style="color: #999; font-size: 12px; margin-top: 5px;">-1表示无限制</div>
        </el-form-item>
        
        <el-form-item label="IP上限">
          <el-input-number v-model="createForm.max_ip" :min="-1" />
          <div style="color: #999; font-size: 12px; margin-top: 5px;">-1表示无限制</div>
        </el-form-item>
        
        <el-form-item label="备注">
          <el-input v-model="createForm.note" type="textarea" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateSave">确定</el-button>
      </template>
    </el-dialog>
    
    <el-dialog v-model="editDialogVisible" title="编辑卡密" width="600px">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="有效时长">
          <div style="display: flex; gap: 10px; align-items: center;">
            <el-input-number v-model="editForm.duration_value" :min="0" style="width: 150px;" />
            <el-select v-model="editForm.duration_unit" style="width: 100px;">
              <el-option label="秒" value="second" />
              <el-option label="天" value="day" />
              <el-option label="周" value="week" />
              <el-option label="月" value="month" />
              <el-option label="季" value="quarter" />
              <el-option label="年" value="year" />
            </el-select>
          </div>
          <div style="color: #999; font-size: 12px; margin-top: 5px;">0表示永久</div>
        </el-form-item>
        
        <el-form-item label="卡密类型">
          <el-input v-model="editForm.card_type" placeholder="normal" />
        </el-form-item>
        
        <el-form-item label="设备码上限">
          <el-input-number v-model="editForm.max_hwid" :min="-1" />
          <div style="color: #999; font-size: 12px; margin-top: 5px;">-1表示无限制</div>
        </el-form-item>
        
        <el-form-item label="IP上限">
          <el-input-number v-model="editForm.max_ip" :min="-1" />
          <div style="color: #999; font-size: 12px; margin-top: 5px;">-1表示无限制</div>
        </el-form-item>
        
        <el-form-item label="备注">
          <el-input v-model="editForm.note" type="textarea" />
        </el-form-item>
        
        <el-form-item label="专属信息">
          <el-input v-model="editForm.custom_data" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleEditSave">保存</el-button>
      </template>
    </el-dialog>
    
    <el-dialog v-model="detailDialogVisible" title="卡密详情" width="700px">
      <el-descriptions :column="2" border v-if="currentCard">
        <el-descriptions-item label="卡密">{{ currentCard.card_key }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentCard.activated ? 'success' : 'info'">
            {{ currentCard.activated ? '已激活' : '未激活' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="时长">{{ formatDuration(currentCard.duration) }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ currentCard.card_type }}</el-descriptions-item>
        <el-descriptions-item label="设备码列表" :span="2">
          <el-tag v-for="hwid in currentCard.hwid_list" :key="hwid" style="margin-right: 5px;">
            {{ hwid }}
          </el-tag>
          <span v-if="!currentCard.hwid_list || currentCard.hwid_list.length === 0">无</span>
        </el-descriptions-item>
        <el-descriptions-item label="IP列表" :span="2">
          <el-tag v-for="ip in currentCard.ip_list" :key="ip" style="margin-right: 5px;">
            {{ ip }}
          </el-tag>
          <span v-if="!currentCard.ip_list || currentCard.ip_list.length === 0">无</span>
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ currentCard.note || '无' }}</el-descriptions-item>
        <el-descriptions-item label="专属信息" :span="2">{{ currentCard.custom_data || '无' }}</el-descriptions-item>
      </el-descriptions>
      
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
    
    <el-dialog v-model="batchUpdateDialogVisible" title="批量修改卡密" width="600px">
      <el-alert type="info" :closable="false" style="margin-bottom: 20px;">
        已选择 {{ selectedCards.length }} 个卡密,只修改填写的字段
      </el-alert>
      <el-form :model="batchUpdateForm" label-width="120px">
        <el-form-item label="有效时长">
          <el-checkbox v-model="batchUpdateForm.update_duration" style="margin-bottom: 10px;">
            修改时长
          </el-checkbox>
          <div v-if="batchUpdateForm.update_duration" style="display: flex; gap: 10px; align-items: center;">
            <el-input-number v-model="batchUpdateForm.duration_value" :min="0" style="width: 150px;" />
            <el-select v-model="batchUpdateForm.duration_unit" style="width: 100px;">
              <el-option label="秒" value="second" />
              <el-option label="天" value="day" />
              <el-option label="周" value="week" />
              <el-option label="月" value="month" />
              <el-option label="季" value="quarter" />
              <el-option label="年" value="year" />
            </el-select>
          </div>
        </el-form-item>
        
        <el-form-item label="卡密类型">
          <el-checkbox v-model="batchUpdateForm.update_card_type" style="margin-bottom: 10px;">
            修改类型
          </el-checkbox>
          <el-input v-if="batchUpdateForm.update_card_type" v-model="batchUpdateForm.card_type" placeholder="normal" />
        </el-form-item>
        
        <el-form-item label="设备码上限">
          <el-checkbox v-model="batchUpdateForm.update_max_hwid" style="margin-bottom: 10px;">
            修改设备码上限
          </el-checkbox>
          <el-input-number v-if="batchUpdateForm.update_max_hwid" v-model="batchUpdateForm.max_hwid" :min="-1" />
        </el-form-item>
        
        <el-form-item label="IP上限">
          <el-checkbox v-model="batchUpdateForm.update_max_ip" style="margin-bottom: 10px;">
            修改IP上限
          </el-checkbox>
          <el-input-number v-if="batchUpdateForm.update_max_ip" v-model="batchUpdateForm.max_ip" :min="-1" />
        </el-form-item>
        
        <el-form-item label="备注">
          <el-checkbox v-model="batchUpdateForm.update_note" style="margin-bottom: 10px;">
            修改备注
          </el-checkbox>
          <el-input v-if="batchUpdateForm.update_note" v-model="batchUpdateForm.note" type="textarea" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="batchUpdateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBatchUpdateSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getProjects } from '@/api/project'
import { getCards, createCards, deleteCard, updateCard, batchUpdateCards, batchDeleteCards } from '@/api/card'
import { ElMessage, ElMessageBox } from 'element-plus'

// 时间单位定义
const TIME_UNITS = {
  second: { label: '秒', value: 1 },
  day: { label: '天', value: 86400 },
  week: { label: '周', value: 604800 },
  month: { label: '月', value: 2592000 },
  quarter: { label: '季', value: 7776000 },
  year: { label: '年', value: 31536000 }
}

// 秒转换为最合适的单位
const formatDuration = (seconds) => {
  if (seconds === 0) return '永久'
  
  const units = [
    { name: '年', value: 31536000 },
    { name: '季', value: 7776000 },
    { name: '月', value: 2592000 },
    { name: '周', value: 604800 },
    { name: '天', value: 86400 }
  ]
  
  for (const unit of units) {
    if (seconds % unit.value === 0) {
      return `${seconds / unit.value}${unit.name}`
    }
  }
  
  return `${seconds}秒`
}

// 秒转换为数值和单位
const secondsToUnitValue = (seconds) => {
  if (seconds === 0) return { value: 0, unit: 'day' }
  
  const units = [
    { key: 'year', value: 31536000 },
    { key: 'quarter', value: 7776000 },
    { key: 'month', value: 2592000 },
    { key: 'week', value: 604800 },
    { key: 'day', value: 86400 }
  ]
  
  for (const unit of units) {
    if (seconds % unit.value === 0) {
      return { value: seconds / unit.value, unit: unit.key }
    }
  }
  
  return { value: seconds, unit: 'second' }
}

// 数值和单位转换为秒
const unitValueToSeconds = (value, unit) => {
  return value * (TIME_UNITS[unit]?.value || 1)
}

const route = useRoute()
const loading = ref(false)
const projects = ref([])
const selectedProjectId = ref(null)
const cards = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const batchUpdateDialogVisible = ref(false)
const currentCard = ref(null)
const selectedCards = ref([])

const createForm = ref({
  mode: 'batch',
  card_key: '',
  prefix: '',
  suffix: '',
  count: 10,
  duration_value: 30,
  duration_unit: 'day',
  card_type: 'normal',
  max_hwid: -1,
  max_ip: -1,
  note: ''
})

const editForm = ref({
  duration_value: 30,
  duration_unit: 'day',
  card_type: 'normal',
  max_hwid: -1,
  max_ip: -1,
  note: '',
  custom_data: ''
})

const batchUpdateForm = ref({
  update_duration: false,
  duration_value: 30,
  duration_unit: 'day',
  update_card_type: false,
  card_type: 'normal',
  update_max_hwid: false,
  max_hwid: -1,
  update_max_ip: false,
  max_ip: -1,
  update_note: false,
  note: ''
})

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
      page_size: pageSize.value
    })
    cards.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  createForm.value = {
    mode: 'batch',
    card_key: '',
    prefix: '',
    suffix: '',
    count: 10,
    duration_value: 30,
    duration_unit: 'day',
    card_type: 'normal',
    max_hwid: -1,
    max_ip: -1,
    note: ''
  }
  createDialogVisible.value = true
}

const handleCreateSave = async () => {
  try {
    const duration = unitValueToSeconds(createForm.value.duration_value, createForm.value.duration_unit)
    
    const params = {
      project_id: selectedProjectId.value,
      duration: duration,
      card_type: createForm.value.card_type,
      max_hwid: createForm.value.max_hwid,
      max_ip: createForm.value.max_ip,
      note: createForm.value.note
    }
    
    if (createForm.value.mode === 'custom') {
      params.card_key = createForm.value.card_key
      params.count = 1
    } else {
      params.prefix = createForm.value.prefix
      params.suffix = createForm.value.suffix
      params.count = createForm.value.count
    }
    
    await createCards(params)
    ElMessage.success('生成成功')
    createDialogVisible.value = false
    loadCards()
  } catch (error) {
    console.error(error)
  }
}

const handleEdit = (row) => {
  const { value, unit } = secondsToUnitValue(row.duration)
  editForm.value = {
    id: row.id,
    duration_value: value,
    duration_unit: unit,
    card_type: row.card_type,
    max_hwid: row.max_hwid,
    max_ip: row.max_ip,
    note: row.note,
    custom_data: row.custom_data || ''
  }
  editDialogVisible.value = true
}

const handleEditSave = async () => {
  try {
    const duration = unitValueToSeconds(editForm.value.duration_value, editForm.value.duration_unit)
    
    await updateCard(editForm.value.id, {
      duration: duration,
      card_type: editForm.value.card_type,
      max_hwid: editForm.value.max_hwid,
      max_ip: editForm.value.max_ip,
      note: editForm.value.note,
      custom_data: editForm.value.custom_data
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

const handlePageChange = (newPage) => {
  page.value = newPage
  loadCards()
}

const handleSelectionChange = (selection) => {
  selectedCards.value = selection
}

const handleBatchUpdate = () => {
  batchUpdateForm.value = {
    update_duration: false,
    duration_value: 30,
    duration_unit: 'day',
    update_card_type: false,
    card_type: 'normal',
    update_max_hwid: false,
    max_hwid: -1,
    update_max_ip: false,
    max_ip: -1,
    update_note: false,
    note: ''
  }
  batchUpdateDialogVisible.value = true
}

const handleBatchUpdateSave = async () => {
  try {
    const updateData = {}
    
    if (batchUpdateForm.value.update_duration) {
      updateData.duration = unitValueToSeconds(batchUpdateForm.value.duration_value, batchUpdateForm.value.duration_unit)
    }
    if (batchUpdateForm.value.update_card_type) {
      updateData.card_type = batchUpdateForm.value.card_type
    }
    if (batchUpdateForm.value.update_max_hwid) {
      updateData.max_hwid = batchUpdateForm.value.max_hwid
    }
    if (batchUpdateForm.value.update_max_ip) {
      updateData.max_ip = batchUpdateForm.value.max_ip
    }
    if (batchUpdateForm.value.update_note) {
      updateData.note = batchUpdateForm.value.note
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

