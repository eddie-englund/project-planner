import { createFetch } from '@vueuse/core'

export const useApi = createFetch({
  baseUrl: import.meta.env.VITE_API_BASE_URL,
  fetchOptions: {
    mode: 'cors',
  },
})
