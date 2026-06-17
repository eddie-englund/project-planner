<script setup lang="ts">
import { ref } from 'vue'
import { useTicketsStore } from '@/stores/tickets'
import AppButton from '@/components/AppButton.vue'
import type { Status } from '@/types/ticket'

const props = defineProps<{
  open: boolean
  projectId: string
  topicId: string
  statuses: Status[]
}>()
const emit = defineEmits<{ close: []; created: [] }>()

const store = useTicketsStore()

const title = ref('')
const body = ref('')
const urls = ref<string[]>([])
const selectedStatusId = ref<string>('')
const submitting = ref(false)

function addUrl() {
  urls.value.push('')
}

function removeUrl(i: number) {
  urls.value.splice(i, 1)
}

async function submit() {
  if (!title.value.trim()) return
  submitting.value = true
  await store.create(props.projectId, props.topicId, {
    title: title.value.trim(),
    body: body.value.trim(),
    urls: urls.value.filter((u) => u.trim()),
    statusId: selectedStatusId.value || undefined,
  })
  submitting.value = false
  reset()
  emit('created')
  emit('close')
}

function reset() {
  title.value = ''
  body.value = ''
  urls.value = []
  selectedStatusId.value = props.statuses[0]?.id ?? ''
}

function close() {
  reset()
  emit('close')
}

// set default status when statuses load
import { watch } from 'vue'
watch(
  () => props.statuses,
  (s) => { if (s.length > 0 && !selectedStatusId.value) selectedStatusId.value = s[0]?.id ?? '' },
  { immediate: true }
)
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
        @click.self="close"
      >
        <div class="w-full max-w-2xl rounded-2xl bg-zinc-900 p-6 shadow-2xl ring-1 ring-white/10">
          <h2 class="mb-5 text-base font-semibold text-zinc-100">New ticket</h2>

          <!-- Title -->
          <div class="mb-4">
            <label class="mb-1.5 block text-xs font-medium text-zinc-400">Title</label>
            <input
              v-model="title"
              type="text"
              placeholder="Ticket title"
              class="w-full rounded-lg bg-zinc-800 px-3 py-2.5 text-base text-zinc-100 placeholder-zinc-500 outline-none ring-1 ring-zinc-700 transition focus:ring-zinc-500"
              @keydown.esc="close"
            />
          </div>

          <!-- Body -->
          <div class="mb-4">
            <label class="mb-1.5 block text-xs font-medium text-zinc-400">Description <span class="text-zinc-600">(optional)</span></label>
            <textarea
              v-model="body"
              rows="12"
              placeholder="Write in markdown…"
              class="w-full resize-none rounded-lg bg-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-500 outline-none ring-1 ring-zinc-700 transition focus:ring-zinc-500"
              @keydown.esc="close"
            />
            <p class="mt-1.5 text-xs text-zinc-600">Markdown supported</p>
          </div>

          <!-- Status -->
          <div v-if="statuses.length > 0" class="mb-4">
            <label class="mb-2 block text-xs font-medium text-zinc-400">Status</label>
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="s in statuses"
                :key="s.id"
                type="button"
                class="cursor-pointer rounded-full px-3 py-1 text-xs font-medium transition"
                :style="selectedStatusId === s.id ? { backgroundColor: s.color, color: '#fff' } : {}"
                :class="selectedStatusId === s.id ? '' : 'bg-zinc-800 text-zinc-400 hover:text-zinc-200'"
                @click="selectedStatusId = s.id"
              >{{ s.name }}</button>
            </div>
          </div>

          <!-- URLs -->
          <div class="mb-6">
            <div class="mb-2 flex items-center justify-between">
              <label class="text-xs font-medium text-zinc-400">Links <span class="text-zinc-600">(optional)</span></label>
              <button
                type="button"
                class="text-xs text-zinc-500 transition hover:text-zinc-300"
                @click="addUrl"
              >
                + Add link
              </button>
            </div>
            <div class="flex flex-col gap-2">
              <div v-for="(_, i) in urls" :key="i" class="flex items-center gap-2">
                <input
                  v-model="urls[i]"
                  type="url"
                  placeholder="https://…"
                  class="min-w-0 flex-1 rounded-lg bg-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-500 outline-none ring-1 ring-zinc-700 transition focus:ring-zinc-500"
                />
                <button
                  type="button"
                  class="shrink-0 text-zinc-600 transition hover:text-zinc-400"
                  @click="removeUrl(i)"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2">
            <AppButton variant="ghost" @click="close">Cancel</AppButton>
            <AppButton
              variant="secondary"
              :disabled="submitting || !title.trim()"
              @click="submit"
            >
              {{ submitting ? 'Creating…' : 'Create ticket' }}
            </AppButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.15s ease;
}
.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.15s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.97);
}
</style>
