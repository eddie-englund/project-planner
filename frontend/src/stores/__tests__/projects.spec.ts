import { ref } from 'vue'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useProjectsStore } from '../projects'
import { useApi } from '@/composables/useApi'
import type { Project } from '@/types/project'

vi.mock('@/composables/useApi', () => ({ useApi: vi.fn() }))

const mockApi = vi.mocked(useApi)

const project1: Project = { id: 'p1', title: 'Alpha', color: '#aaa', imageUrl: null, createdBy: 'u1', createdAt: '' }
const project2: Project = { id: 'p2', title: 'Beta', color: '#bbb', imageUrl: null, createdBy: 'u1', createdAt: '' }

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

describe('useProjectsStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('fetchAll', () => {
    it('populates items on success', async () => {
      setupGet([project1, project2])
      const store = useProjectsStore()
      await store.fetchAll()
      expect(store.items).toEqual([project1, project2])
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('sets error on failure', async () => {
      setupGet(null, new Error('network'))
      const store = useProjectsStore()
      await store.fetchAll()
      expect(store.error).toBe('network')
      expect(store.items).toHaveLength(0)
    })
  })

  describe('create', () => {
    it('appends new project to items', async () => {
      setupMutation(project1)
      const store = useProjectsStore()
      await store.create({ title: 'Alpha', color: '#aaa' })
      expect(store.items).toContainEqual(project1)
    })

    it('sets error on failure', async () => {
      setupMutation(null, new Error('fail'))
      const store = useProjectsStore()
      await store.create({ title: 'X', color: '#000' })
      expect(store.error).toBe('fail')
    })
  })

  describe('update', () => {
    it('replaces matching item', async () => {
      const updated = { ...project1, title: 'Alpha Updated' }
      setupMutation(updated)
      const store = useProjectsStore()
      store.items = [project1, project2]
      await store.update('p1', { title: 'Alpha Updated', color: '#aaa' })
      expect(store.items[0].title).toBe('Alpha Updated')
      expect(store.items[1]).toEqual(project2)
    })
  })

  describe('remove', () => {
    it('removes item by id', async () => {
      setupMutation(null)
      const store = useProjectsStore()
      store.items = [project1, project2]
      await store.remove('p1')
      expect(store.items).toEqual([project2])
    })
  })

  describe('filtered', () => {
    it('returns all when search is empty', async () => {
      setupGet([project1, project2])
      const store = useProjectsStore()
      await store.fetchAll()
      expect(store.filtered).toHaveLength(2)
    })

    it('filters by title case-insensitively', async () => {
      setupGet([project1, project2])
      const store = useProjectsStore()
      await store.fetchAll()
      store.search = 'alph'
      expect(store.filtered).toEqual([project1])
    })
  })
})
