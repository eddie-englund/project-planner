import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useApi } from '@/composables/useApi'
import type { Topic, CreateTopicPayload, UpdateTopicPayload } from '@/types/topic'

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

  async function update(projectId: string, id: string, payload: UpdateTopicPayload) {
    const { data, error: fetchError } = await useApi(`/projects/${projectId}/topics/${id}`)
      .put(payload)
      .json<Topic>()
    if (fetchError.value) {
      error.value = fetchError.value.message
      return
    }
    if (data.value) {
      const idx = items.value.findIndex((t) => t.id === id)
      if (idx !== -1) items.value[idx] = data.value
    }
  }

  async function remove(projectId: string, id: string) {
    const { error: fetchError } = await useApi(`/projects/${projectId}/topics/${id}`).delete()
    if (fetchError.value) {
      error.value = fetchError.value.message
      return
    }
    items.value = items.value.filter((t) => t.id !== id)
  }

  return { items, filtered, search, isLoading, error, fetchAll, create, update, remove }
})
