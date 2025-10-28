import { ref } from 'vue'

export function useDialog() {
  const visible = ref(false)
  
  const open = () => {
    visible.value = true
  }
  
  const close = () => {
    visible.value = false
  }
  
  return {
    visible,
    open,
    close
  }
}

