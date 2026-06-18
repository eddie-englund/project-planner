import { ref } from 'vue'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTopicsStore } from '../topics'
import { useApi } from '@/composables/useApi'
import type { Topic } from '@/types/topic'

vi.mock('@/composables/useApi', () => ({ useApi: vi.fn() }))

const mockApi = vi.mocked(useApi)

const topic1: Topic = { id: 't1', projectId: 'p1', index: 0, title: 'Frontend', color: '#abc', imageUrl: null, createdAt: '' }
const topic2: Topic = { id: 't2', projectId: 'p1', index: 1, title: 'Backend', color: '#cba', imageUrl: null, createdAt: '' }

function setupGet<T>(data: T, error: Error | null = null) {
  const result = { data: ref(data), error: ref(error) }
  mockApi.mockReturnValue({ json: vi.fn().mockResolvedValue(result) } as any)
  return result
}

function setupMutation<T>(data: T, error: Error | null = null) {
  const result = { data: ref(data), error: ref(error) }
  const jsonFn = vi.fn().mockResolvedValue(result)
  mockApi.mockReturnValue({
    json: jsonFn,
    post: vi.fn().mockReturnValue({ json: jsonFn }),
    put: vi.fn().mockReturnValue({ json: jsonFn }),
    delete: vi.fn().mockResolvedValue({ error: ref(error) }),
  } as any)
  return result
}

describe('useTopicsStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('fetchAll', () => {
    it('clears then populates items', async () => {
      setupGet([topic1, topic2])
      const store = useTopicsStore()
      store.items = [topic1]
      await store.fetchAll('p1')
      expect(store.items).toEqual([topic1, topic2])
      expect(store.isLoading).toBe(false)
    })

    it('sets error on failure', async () => {
      setupGet(null, new Error('oops'))
      const store = useTopicsStore()
      await store.fetchAll('p1')
      expect(store.error).toBe('oops')
    })

    it('passes project-scoped URL', async () => {
      setupGet([])
      const store = useTopicsStore()
      await store.fetchAll('proj-99')
      expect(mockApi).toHaveBeenCalledWith('/projects/proj-99/topics')
    })
  })

  describe('create', () => {
    it('appends new topic', async () => {
      setupMutation(topic1)
      const store = useTopicsStore()
      await store.create('p1', { title: 'Frontend', color: '#abc', index: 0 })
      expect(store.items).toContainEqual(topic1)
    })
  })

  describe('update', () => {
    it('replaces matching topic', async () => {
      const updated = { ...topic1, title: 'Frontend v2' }
      setupMutation(updated)
      const store = useTopicsStore()
      store.items = [topic1, topic2]
      await store.update('p1', 't1', { title: 'Frontend v2', color: '#abc' })
      expect(store.items[0].title).toBe('Frontend v2')
      expect(store.items[1]).toEqual(topic2)
    })
  })

  describe('remove', () => {
    it('removes topic by id', async () => {
      setupMutation(null)
      const store = useTopicsStore()
      store.items = [topic1, topic2]
      await store.remove('p1', 't1')
      expect(store.items).toEqual([topic2])
    })
  })

  describe('filtered', () => {
    it('filters by title', async () => {
      setupGet([topic1, topic2])
      const store = useTopicsStore()
      await store.fetchAll('p1')
      store.search = 'back'
      expect(store.filtered).toEqual([topic2])
    })
  })
})
