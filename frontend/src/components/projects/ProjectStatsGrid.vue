<template>
  <div class="stats-grid">
    <div class="stat-card card-modern hover-lift stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
        <el-icon><Box /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-label">总项目数</div>
        <div class="stat-value">{{ totalCount }}</div>
      </div>
    </div>
    
    <div class="stat-card card-modern hover-lift stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
        <el-icon><Money /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-label">付费项目</div>
        <div class="stat-value">{{ paidCount }}</div>
      </div>
    </div>
    
    <div class="stat-card card-modern hover-lift stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
        <el-icon><Coin /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-label">免费项目</div>
        <div class="stat-value">{{ freeCount }}</div>
      </div>
    </div>
    
    <div class="stat-card card-modern hover-lift stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
        <el-icon><TrendCharts /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-label">活跃项目</div>
        <div class="stat-value">{{ totalCount }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Box, Money, Coin, TrendCharts } from '@element-plus/icons-vue'

defineProps({
  totalCount: {
    type: Number,
    default: 0
  },
  paidCount: {
    type: Number,
    default: 0
  },
  freeCount: {
    type: Number,
    default: 0
  }
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  background: #fff;
  border-radius: var(--radius-lg);
  transition: all var(--duration-normal) var(--ease-out);
  border: 1px solid var(--color-border-light);
  box-shadow: var(--shadow-sm);
}

.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
  font-weight: 500;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1;
}

/* 平板适配 */
@media (min-width: 769px) and (max-width: 1023px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 14px;
  }
  
  .stat-card {
    padding: 20px;
  }
  
  .stat-value {
    font-size: 28px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    margin-bottom: 16px;
  }
  
  .stat-card {
    padding: 16px;
    gap: 12px;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 12px;
  }
}

/* 小屏手机适配 */
@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

/* 大屏优化 - 最多显示4列,超过1400px时保持合理间距 */
@media (min-width: 1400px) {
  .stats-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
  }
}
</style>

