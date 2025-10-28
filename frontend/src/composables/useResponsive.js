import { ref, onMounted, onUnmounted } from 'vue'

export function useResponsive() {
  const isMobile = ref(false)
  const isTablet = ref(false)
  const isDesktop = ref(false)
  const windowWidth = ref(0)
  
  const mobileQuery = window.matchMedia('(max-width: 767px)')
  const tabletQuery = window.matchMedia('(min-width: 768px) and (max-width: 1023px)')
  const desktopQuery = window.matchMedia('(min-width: 1024px)')
  
  const updateResponsiveState = () => {
    windowWidth.value = window.innerWidth
    isMobile.value = mobileQuery.matches
    isTablet.value = tabletQuery.matches
    isDesktop.value = desktopQuery.matches
  }
  
  onMounted(() => {
    updateResponsiveState()
    
    mobileQuery.addEventListener('change', updateResponsiveState)
    tabletQuery.addEventListener('change', updateResponsiveState)
    desktopQuery.addEventListener('change', updateResponsiveState)
    window.addEventListener('resize', updateResponsiveState)
  })
  
  onUnmounted(() => {
    mobileQuery.removeEventListener('change', updateResponsiveState)
    tabletQuery.removeEventListener('change', updateResponsiveState)
    desktopQuery.removeEventListener('change', updateResponsiveState)
    window.removeEventListener('resize', updateResponsiveState)
  })
  
  return {
    isMobile,
    isTablet,
    isDesktop,
    windowWidth
  }
}

