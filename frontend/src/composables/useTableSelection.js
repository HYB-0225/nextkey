import { ref } from 'vue'

export function useTableSelection() {
  const selectedItems = ref([])
  
  const handleSelectionChange = (selection) => {
    selectedItems.value = selection
  }
  
  const clearSelection = () => {
    selectedItems.value = []
  }
  
  return {
    selectedItems,
    handleSelectionChange,
    clearSelection
  }
}

