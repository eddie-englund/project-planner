import { ref } from 'vue'

const message = ref<string | null>(null)
let timer: ReturnType<typeof setTimeout> | null = null

export function useToast() {
  function show(msg: string, duration = 2000) {
    message.value = msg
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => { message.value = null }, duration)
  }

  function hide() {
    message.value = null
    if (timer) clearTimeout(timer)
  }

  return { message, show, hide }
}
