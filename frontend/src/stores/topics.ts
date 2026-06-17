import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useApi } from '@/composables/useApi'
import type { Topic, CreateTopicPayload } from '@/types/topic'

export const useTopicsStore = defineStore('topics', () => {
  const items = ref<Topic[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const search = ref('')
  const filtered = computed(() =>
    items.value.filter(t => t.title.toLowerCase().includes(search.value.toLowerCase()))
  )

  async function fetchAll(projectId: string) {
    items.value = []
    isLoading.value = true
    error.value = null
    const { data, error: fetchError } = await useApi(`/projects/${projectId}/topics`).json<Topic[]>()
    isLoading.value = false
    if (fetchError.value) {
      error.value = fetchError.value.message
      return
    }
    if (data.value) {
      items.value = data.value
    }
  }

  async function create(projectId: string, payload: CreateTopicPayload) {
    const { data, error: fetchError } = await useApi(`/projects/${projectId}/topics`)
      .post(payload)
      .json<Topic>()
    if (fetchError.value) {
      error.value = fetchError.value.message
      return
    }
    if (data.value) {
      items.value.push(data.value)
    }
  }

  return { items, filtered, search, isLoading, error, fetchAll, create }
})
