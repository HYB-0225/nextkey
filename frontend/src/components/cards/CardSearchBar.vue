<template>
  <el-card class="search-bar" shadow="never">
    <div class="search-container">
      <!-- 第一行：主要搜索 -->
      <div class="search-row">
        <el-input 
          v-model="searchForm.keyword" 
          placeholder="搜索卡密" 
          clearable
          class="search-input"
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-input 
          v-model="searchForm.card_type" 
          placeholder="类型" 
          clearable
          class="search-input-sm"
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        />
        
        <el-select 
          v-model="searchForm.activated" 
          placeholder="激活状态" 
          clearable
          class="search-select-sm"
          @change="handleSearch"
        >
          <el-option label="已激活" value="true" />
          <el-option label="未激活" value="false" />
        </el-select>
        
        <el-input 
          v-model="searchForm.note" 
          placeholder="备注" 
          clearable
          class="search-input-sm"
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        />
        
        <el-button type="primary" @click="handleSearch" class="search-btn">
          搜索
        </el-button>
        
        <el-button @click="handleReset" class="reset-btn">
          重置
        </el-button>
        
        <el-button 
          text 
          @click="showAdvanced = !showAdvanced" 
          class="toggle-btn"
        >
          {{ showAdvanced ? '收起' : '展开' }}
          <el-icon :class="{ 'rotate': showAdvanced }">
            <ArrowDown />
          </el-icon>
        </el-button>
      </div>
      
      <!-- 第二行：高级搜索（可展开/收起） -->
      <transition name="slide-fade">
        <div v-show="showAdvanced" class="advanced-row">
          <el-input 
            v-model="searchForm.custom_data" 
            placeholder="专属信息" 
            clearable
            class="search-input-sm"
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
          
          <el-input 
            v-model="searchForm.hwid" 
            placeholder="设备码" 
            clearable
            class="search-input-sm"
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
          
          <el-input 
            v-model="searchForm.ip" 
            placeholder="IP地址" 
            clearable
            class="search-input-sm"
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
          
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            clearable
            class="date-picker"
            @change="handleDateChange"
          />
        </div>
      </transition>
    </div>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { Search, ArrowDown } from '@element-plus/icons-vue'

const emit = defineEmits(['search', 'reset'])

const showAdvanced = ref(false)

const searchForm = ref({
  keyword: '',
  card_type: '',
  activated: '',
  note: '',
  custom_data: '',
  hwid: '',
  ip: ''
})

const dateRange = ref(null)

const handleDateChange = () => {
  handleSearch()
}

const handleSearch = () => {
  const params = { ...searchForm.value }
  
  if (dateRange.value && dateRange.value.length === 2) {
    params.start_time = dateRange.value[0].toISOString().split('T')[0] + ' 00:00:00'
    params.end_time = dateRange.value[1].toISOString().split('T')[0] + ' 23:59:59'
  }
  
  emit('search', params)
}

const handleReset = () => {
  searchForm.value = {
    keyword: '',
    card_type: '',
    activated: '',
    note: '',
    custom_data: '',
    hwid: '',
    ip: ''
  }
  dateRange.value = null
  emit('reset')
}
</script>

<style scoped>
.search-bar {
  margin-bottom: 20px;
  border-radius: var(--radius-lg);
}

.search-bar :deep(.el-card__body) {
  padding: 16px;
}

.search-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.search-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.advanced-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  padding-top: 4px;
}

.search-input {
  flex: 1;
  min-width: 200px;
  max-width: 300px;
}

.search-input-sm {
  width: 160px;
}

.search-select-sm {
  width: 140px;
}

.date-picker {
  width: 280px;
}

.search-btn {
  min-width: 80px;
}

.reset-btn {
  min-width: 80px;
}

.toggle-btn {
  margin-left: auto;
  color: var(--color-primary);
}

.toggle-btn :deep(.el-icon) {
  margin-left: 4px;
  transition: transform var(--duration-normal) var(--ease-out);
}

.toggle-btn :deep(.el-icon.rotate) {
  transform: rotate(180deg);
}

/* 展开收起动画 */
.slide-fade-enter-active {
  transition: all var(--duration-normal) var(--ease-out);
}

.slide-fade-leave-active {
  transition: all var(--duration-fast) var(--ease-in);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 平板适配 */
@media (min-width: 769px) and (max-width: 1200px) {
  .search-input {
    min-width: 180px;
    max-width: 250px;
  }
  
  .search-input-sm {
    width: 140px;
  }
  
  .search-select-sm {
    width: 120px;
  }
  
  .date-picker {
    width: 240px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .search-bar :deep(.el-card__body) {
    padding: 12px;
  }
  
  .search-container {
    gap: 10px;
  }
  
  .search-row {
    gap: 8px;
  }
  
  .advanced-row {
    gap: 8px;
  }
  
  .search-input {
    width: 100%;
    min-width: unset;
    max-width: unset;
  }
  
  .search-input-sm {
    flex: 1;
    min-width: 120px;
  }
  
  .search-select-sm {
    flex: 1;
    min-width: 100px;
  }
  
  .date-picker {
    width: 100%;
  }
  
  .search-btn,
  .reset-btn {
    flex: 1;
    min-width: 60px;
  }
  
  .toggle-btn {
    width: 100%;
    margin-left: 0;
  }
}

/* 输入框样式增强 */
.search-bar :deep(.el-input__wrapper) {
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-out);
}

.search-bar :deep(.el-input__wrapper:hover) {
  border-color: var(--color-primary);
}

.search-bar :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--color-primary) inset;
}

.search-bar :deep(.el-button) {
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--duration-fast) var(--ease-out);
}

.search-bar :deep(.el-button:hover) {
  transform: translateY(-1px);
}

.search-bar :deep(.el-button:active) {
  transform: translateY(0);
}
</style>

