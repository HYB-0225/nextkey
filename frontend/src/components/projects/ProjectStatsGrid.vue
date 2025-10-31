<template>
  <div class="stats-grid">
    <div class="stat-card pixel-card stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #FF8C42 0%, #FF6B35 100%); border: 2px solid #FFD93D; box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.3);">
        <el-icon><Box /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-label">总项目数</div>
        <div class="stat-value">{{ totalCount }}</div>
      </div>
    </div>
    
    <div class="stat-card pixel-card stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #FFD93D 0%, #FFA400 100%); border: 2px solid #FFE873; box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.3);">
        <el-icon><Money /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-label">付费项目</div>
        <div class="stat-value">{{ paidCount }}</div>
      </div>
    </div>
    
    <div class="stat-card pixel-card stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #FFA400 0%, #FF8C42 100%); border: 2px solid #FFD93D; box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.3);">
        <el-icon><Coin /></el-icon>
      </div>
      <div class="stat-content">
        <div class="stat-label">免费项目</div>
        <div class="stat-value">{{ freeCount }}</div>
      </div>
    </div>
    
    <div class="stat-card pixel-card stagger-item">
      <div class="stat-icon" style="background: linear-gradient(135deg, #D2691E 0%, #A0522D 100%); border: 2px solid #E67E30; box-shadow: 2px 2px 0 0 rgba(0, 0, 0, 0.3);">
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
  background: var(--color-bg);
  border-radius: 0;
  transition: all var(--duration-fast) steps(2);
  border: 3px solid var(--color-border);
  box-shadow: 4px 4px 0 0 rgba(0, 0, 0, 0.1);
}

.stat-card:hover {
  box-shadow: 6px 6px 0 0 rgba(0, 0, 0, 0.15);
  transform: translate(-2px, -2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 24px;
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
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
  }
  
  .stat-card {
    padding: 12px;
    gap: 10px;
    min-height: auto;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
    font-size: 18px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .stat-label {
    font-size: 11px;
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

