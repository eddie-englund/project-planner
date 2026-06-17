import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useApi } from '@/composables/useApi'
import type { Status } from '@/types/ticket'

export const useStatusesStore = defineStore('statuses', () => {
  const items = ref<Status[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function fetchAll(projectId: string) {
    isLoading.value = true
    error.value = null
    const { data, error: err } = await useApi(`/projects/${projectId}/statuses`).json<Status[]>()
    isLoading.value = false
    if (err.value) {
      error.value = 'Failed to load statuses'
      return
    }
    items.value = data.value ?? []
  }

  return { items, isLoading, error, fetchAll }
})
