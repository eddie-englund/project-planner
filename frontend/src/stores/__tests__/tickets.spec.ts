import { ref } from 'vue'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTicketsStore } from '../tickets'
import { useApi } from '@/composables/useApi'
import type { TicketWithTopic } from '@/types/ticket'

vi.mock('@/composables/useApi', () => ({ useApi: vi.fn() }))

const mockApi = vi.mocked(useApi)

const ticket1: TicketWithTopic = {
  id: 'tk1',
  topicId: 't1',
  statusId: null,
  title: 'Fix bug',
  body: '',
  urls: [],
  createdAt: '2026-01-01T00:00:00Z',
  topicColor: '#10b981',
  topicTitle: 'Backend',
}

const ticket2: TicketWithTopic = {
  id: 'tk2',
  topicId: 't2',
  statusId: 's1',
  title: 'Design UI',
  body: '',
  urls: [],
  createdAt: '2026-01-02T00:00:00Z',
  topicColor: '#3b82f6',
  topicTitle: 'Frontend',
}

function setupGet<T>(data: T, error: Error | null = null) {
  const result = { data: ref(data), error: ref(error) }
  mockApi.mockReturnValue({ json: vi.fn().mockResolvedValue(result) } as any)
  return result
}

describe('useTicketsStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('fetchAllByProject', () => {
    it('populates projectItems on success', async () => {
      setupGet([ticket1, ticket2])
      const store = useTicketsStore()
      await store.fetchAllByProject('p1')
      expect(store.projectItems).toEqual([ticket1, ticket2])
      expect(store.isLoadingProject).toBe(false)
      expect(store.error).toBeNull()
    })

    it('sets error on failure', async () => {
      setupGet(null, new Error('network'))
      const store = useTicketsStore()
      await store.fetchAllByProject('p1')
      expect(store.error).toBe('Failed to load tickets')
      expect(store.projectItems).toHaveLength(0)
    })

    it('calls correct project-scoped URL', async () => {
      setupGet([])
      const store = useTicketsStore()
      await store.fetchAllByProject('proj-99')
      expect(mockApi).toHaveBeenCalledWith('/projects/proj-99/tickets')
    })

    it('sets isLoadingProject to false after fetch', async () => {
      setupGet([ticket1])
      const store = useTicketsStore()
      await store.fetchAllByProject('p1')
      expect(store.isLoadingProject).toBe(false)
    })
  })
})
