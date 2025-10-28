<template>
  <div class="modern-dialog theme-success">
    <el-dialog 
      v-model="dialogVisible" 
      title="卡密生成成功" 
      :width="isMobile ? '95%' : '700px'"
      :fullscreen="isMobile"
      :close-on-click-modal="false"
      @close="handleClose"
    >
      <div class="success-header">
        <el-icon :size="48" color="#67C23A"><SuccessFilled /></el-icon>
        <div class="success-text">
          <h3>成功生成 {{ cards.length }} 个卡密</h3>
          <p>您可以导出或复制这些卡密</p>
        </div>
      </div>
      
      <div class="cards-list">
        <el-table :data="cards" max-height="400" style="width: 100%">
          <el-table-column type="index" label="序号" width="60" />
          <el-table-column prop="card_key" label="卡密" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="card-key-text">{{ row.card_key }}</span>
            </template>
          </el-table-column>
          <el-table-column label="有效时长" width="120">
            <template #default="{ row }">
              {{ formatDuration(row.duration) }}
            </template>
          </el-table-column>
          <el-table-column prop="card_type" label="类型" width="100" />
        </el-table>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-dropdown split-button type="primary" @click="handleExport('txt')" style="margin-right: 10px;">
            导出卡密
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleExport('json')">导出为 JSON</el-dropdown-item>
                <el-dropdown-item @click="handleExport('txt')">导出为 TXT</el-dropdown-item>
                <el-dropdown-item @click="handleExport('csv')">导出为 CSV</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button type="success" @click="handleCopyAll">
            <el-icon><CopyDocument /></el-icon>
            复制全部
          </el-button>
          <el-button @click="handleClose">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { SuccessFilled, CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useResponsive } from '@/composables/useResponsive'
import { formatDuration } from '@/composables/useDuration'
import { exportToJSON, exportToTXT, exportToCSV } from '@/utils/export'

const { isMobile } = useResponsive()

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  cards: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:visible'])

const dialogVisible = ref(false)

watch(() => props.visible, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

const handleClose = () => {
  dialogVisible.value = false
}

const handleExport = (format) => {
  const timestamp = new Date().toISOString().slice(0, 19).replace(/[:-]/g, '').replace('T', '_')
  const filename = `cards_created_${timestamp}`
  
  switch (format) {
    case 'json':
      exportToJSON(props.cards, filename)
      break
    case 'txt':
      exportToTXT(props.cards, filename)
      break
    case 'csv':
      exportToCSV(props.cards, filename)
      break
  }
  
  ElMessage.success(`成功导出 ${props.cards.length} 个卡密`)
}

const handleCopyAll = async () => {
  const cardKeys = props.cards.map(card => card.card_key).join('\n')
  
  try {
    await navigator.clipboard.writeText(cardKeys)
    ElMessage.success('已复制到剪贴板')
  } catch (err) {
    const textarea = document.createElement('textarea')
    textarea.value = cardKeys
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    
    try {
      document.execCommand('copy')
      ElMessage.success('已复制到剪贴板')
    } catch (e) {
      ElMessage.error('复制失败，请手动复制')
    }
    
    document.body.removeChild(textarea)
  }
}
</script>

<style scoped>
.success-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding: 20px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border-radius: var(--radius-lg);
}

.success-text h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: var(--color-text-primary);
}

.success-text p {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.cards-list {
  margin-bottom: 20px;
}

.card-key-text {
  font-family: 'Courier New', monospace;
  font-weight: 500;
  color: var(--color-primary);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .success-header {
    flex-direction: column;
    text-align: center;
    padding: 16px;
  }
  
  .dialog-footer {
    flex-direction: column;
  }
  
  .dialog-footer :deep(.el-button),
  .dialog-footer :deep(.el-dropdown) {
    width: 100%;
  }
}
</style>

