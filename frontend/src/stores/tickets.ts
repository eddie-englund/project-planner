import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useApi } from '@/composables/useApi'
import type { Ticket, CreateTicketPayload } from '@/types/ticket'

export const useTicketsStore = defineStore('tickets', () => {
  const items = ref<Ticket[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const search = ref('')
  const sort = ref<'newest' | 'oldest' | 'alpha'>('newest')

  function ticketsForStatus(statusId: string | null) {
    return computed(() => {
      let list = items.value.filter((t) => t.statusId === statusId)
      if (search.value) {
        const q = search.value.toLowerCase()
        list = list.filter(
          (t) => t.title.toLowerCase().includes(q) || t.body.toLowerCase().includes(q)
        )
      }
      if (sort.value === 'oldest') return [...list].sort((a, b) => a.createdAt.localeCompare(b.createdAt))
      if (sort.value === 'alpha') return [...list].sort((a, b) => a.title.localeCompare(b.title))
      return [...list].sort((a, b) => b.createdAt.localeCompare(a.createdAt))
    })
  }

  async function fetchAll(projectId: string, topicId: string) {
    isLoading.value = true
    error.value = null
    const { data, error: err } = await useApi(
      `/projects/${projectId}/topics/${topicId}/tickets`
    ).json<Ticket[]>()
    isLoading.value = false
    if (err.value) {
      error.value = 'Failed to load tickets'
      return
    }
    items.value = data.value ?? []
  }

  async function create(projectId: string, topicId: string, payload: CreateTicketPayload) {
    const { data, error: err } = await useApi(
      `/projects/${projectId}/topics/${topicId}/tickets`
    ).post(payload).json<Ticket>()
    if (err.value || !data.value) throw new Error('Failed to create ticket')
    items.value.push(data.value)
    return data.value
  }

  async function update(
    projectId: string,
    topicId: string,
    ticketId: string,
    payload: { title?: string; body?: string; urls?: string[]; statusId?: string | null }
  ) {
    const ticket = items.value.find((t) => t.id === ticketId)
    if (!ticket) return
    const merged = {
      title: ticket.title,
      body: ticket.body,
      urls: ticket.urls,
      statusId: ticket.statusId,
      ...payload,
    }
    const { data, error: err } = await useApi(
      `/projects/${projectId}/topics/${topicId}/tickets/${ticketId}`
    ).put(merged).json<Ticket>()
    if (err.value || !data.value) throw new Error('Failed to update ticket')
    const idx = items.value.findIndex((t) => t.id === ticketId)
    if (idx !== -1) items.value[idx] = data.value
    return data.value
  }

  async function updateStatus(
    projectId: string,
    topicId: string,
    ticketId: string,
    statusId: string | null
  ) {
    return update(projectId, topicId, ticketId, { statusId })
  }

  return { items, isLoading, error, search, sort, ticketsForStatus, fetchAll, create, update, updateStatus }
})
