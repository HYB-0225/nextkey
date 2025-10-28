import { ref } from 'vue'

export function usePagination(initialPageSize = 20) {
  const page = ref(1)
  const pageSize = ref(initialPageSize)
  const total = ref(0)
  
  const handlePageChange = (newPage) => {
    page.value = newPage
  }
  
  const resetPage = () => {
    page.value = 1
  }
  
  return {
    page,
    pageSize,
    total,
    handlePageChange,
    resetPage
  }
}

