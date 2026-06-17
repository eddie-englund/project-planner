<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { marked } from 'marked'
import { useTicketsStore } from '@/stores/tickets'
import { useToast } from '@/composables/useToast'
import AppButton from '@/components/AppButton.vue'
import type { Ticket, Status } from '@/types/ticket'

const props = defineProps<{
  ticket: Ticket | null
  statuses: Status[]
  projectId: string
  topicId: string
}>()

const emit = defineEmits<{ close: [] }>()

const store = useTicketsStore()
const toast = useToast()

// --- Local copy + draft state (must be declared before the immediate watch) ---
const local = ref<Ticket | null>(null)
const editingTitle = ref(false)
const titleDraft = ref('')
const titleInput = ref<HTMLInputElement | null>(null)
const editingBody = ref(false)
const bodyDraft = ref('')
const bodyInput = ref<HTMLTextAreaElement | null>(null)

watch(
  () => props.ticket,
  (t) => {
    local.value = t ? { ...t } : null
    titleDraft.value = t?.title ?? ''
    bodyDraft.value = t?.body ?? ''
    editingTitle.value = false
    editingBody.value = false
  },
  { immediate: true }
)

// --- Status helpers ---
const currentStatus = computed(() => props.statuses.find((s) => s.id === local.value?.statusId))
const isTerminal = computed(() => currentStatus.value?.isTerminal ?? false)
const statusColor = computed(() => currentStatus.value?.color ?? '#52525b')

// --- Confirm dialog ---
const confirm = ref<{ message: string; onConfirm: () => void } | null>(null)

function withConfirm(message: string, action: () => void) {
  confirm.value = { message, onConfirm: action }
}

function acceptConfirm() {
  confirm.value?.onConfirm()
  confirm.value = null
}

function cancelConfirm() {
  confirm.value = null
}

// --- Status change ---
async function changeStatus(statusId: string) {
  if (!local.value) return
  local.value = { ...local.value, statusId }
  await store.update(props.projectId, props.topicId, local.value.id, { statusId })
  toast.show('Status updated')
}

function requestStatusChange(s: Status) {
  if (s.id === local.value?.statusId) return
  withConfirm(
    `Change status to "${s.name}"?`,
    () => changeStatus(s.id)
  )
}

function requestToggleClosed() {
  if (!local.value) return
  if (isTerminal.value) {
    const open = props.statuses.find((s) => !s.isTerminal)
    if (open) withConfirm('Reopen this ticket?', () => changeStatus(open.id))
  } else {
    const closed = props.statuses.find((s) => s.isTerminal)
    if (closed) withConfirm('Close this ticket?', () => changeStatus(closed.id))
  }
}

// --- Inline title editing ---
const titleDirty = computed(() => titleDraft.value !== local.value?.title)

function startEditTitle() {
  titleDraft.value = local.value?.title ?? ''
  editingTitle.value = true
  setTimeout(() => titleInput.value?.focus(), 0)
}

async function saveTitle() {
  if (!local.value || !titleDirty.value) {
    editingTitle.value = false
    return
  }
  local.value = { ...local.value, title: titleDraft.value }
  await store.update(props.projectId, props.topicId, local.value.id, { title: titleDraft.value })
  editingTitle.value = false
  toast.show('Saved')
}

function cancelTitle() {
  titleDraft.value = local.value?.title ?? ''
  editingTitle.value = false
}

// --- Inline body editing ---
const bodyDirty = computed(() => bodyDraft.value !== local.value?.body)

function startEditBody() {
  bodyDraft.value = local.value?.body ?? ''
  editingBody.value = true
  setTimeout(() => bodyInput.value?.focus(), 0)
}

async function saveBody() {
  if (!local.value || !bodyDirty.value) {
    editingBody.value = false
    return
  }
  local.value = { ...local.value, body: bodyDraft.value }
  await store.update(props.projectId, props.topicId, local.value.id, { body: bodyDraft.value })
  editingBody.value = false
  toast.show('Saved')
}

function cancelBody() {
  bodyDraft.value = local.value?.body ?? ''
  editingBody.value = false
}

const parsedBody = computed(() =>
  local.value?.body ? (marked.parse(local.value.body) as string) : ''
)

// --- URL editing ---
const editingUrls = ref(false)
const urlInputs = ref<string[]>([])

function startEditUrls() {
  editingUrls.value = true
  urlInputs.value = [...(local.value?.urls ?? [])]
}

function addUrl() { urlInputs.value.push('') }
function removeUrl(i: number) { urlInputs.value.splice(i, 1) }

async function saveUrls() {
  editingUrls.value = false
  if (!local.value) return
  const urls = urlInputs.value.filter((u) => u.trim())
  local.value = { ...local.value, urls }
  await store.update(props.projectId, props.topicId, local.value.id, { urls })
  toast.show('Saved')
}

function cancelUrls() { editingUrls.value = false }

// --- Keyboard ---
function handlePanelKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (confirm.value) { cancelConfirm(); return }
    if (editingTitle.value) { cancelTitle(); return }
    if (editingBody.value) { cancelBody(); return }
    emit('close')
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div
        v-if="ticket && local"
        class="fixed inset-0 z-40 flex justify-end"
        tabindex="-1"
        @keydown="handlePanelKeydown"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="emit('close')" />

        <!-- Panel -->
        <div
          class="relative z-50 flex h-full flex-col bg-zinc-900 shadow-2xl ring-1 ring-white/5"
          style="width: clamp(500px, 40vw, 100vw)"
          @click.stop
        >
          <!-- Status color stripe -->
          <div
            class="h-1 w-full shrink-0 transition-colors duration-200"
            :style="{ backgroundColor: statusColor }"
          />

          <!-- Header: status pills + close -->
          <div class="flex items-center gap-2 px-5 py-3 shrink-0 border-b border-zinc-800">
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="s in statuses"
                :key="s.id"
                class="cursor-pointer rounded-full px-3 py-0.5 text-xs font-medium transition"
                :style="local.statusId === s.id ? { backgroundColor: s.color, color: '#fff' } : {}"
                :class="local.statusId === s.id ? '' : 'text-zinc-500 hover:text-zinc-300'"
                @click="requestStatusChange(s)"
              >
                {{ s.name }}
              </button>
            </div>
            <div class="flex-1" />
            <button
              class="cursor-pointer flex h-7 w-7 items-center justify-center rounded-md text-zinc-500 transition hover:bg-zinc-800 hover:text-zinc-200"
              @click="emit('close')"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Scrollable content -->
          <div class="flex-1 overflow-y-auto px-6 pb-6 pt-5">

            <!-- Title -->
            <div class="mb-5">
              <template v-if="editingTitle">
                <input
                  ref="titleInput"
                  v-model="titleDraft"
                  class="w-full rounded-lg bg-zinc-800 px-3 py-2 text-xl font-bold text-zinc-100 outline-none ring-1 ring-zinc-600 transition focus:ring-zinc-400 mb-2"
                  @keydown.enter="saveTitle"
                  @keydown.esc.stop="cancelTitle"
                />
                <div class="flex gap-2">
                  <AppButton variant="secondary" :disabled="!titleDirty" @click="saveTitle">Save</AppButton>
                  <AppButton variant="ghost" @click="cancelTitle">Cancel</AppButton>
                </div>
              </template>
              <h2
                v-else
                class="cursor-text text-xl font-bold leading-snug text-zinc-100 hover:text-white group flex items-center gap-2"
                @click="startEditTitle"
              >
                {{ local.title }}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-zinc-600 opacity-0 group-hover:opacity-100 transition" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Z" />
                </svg>
              </h2>
            </div>

            <!-- Body -->
            <div class="mb-5">
              <label class="mb-1.5 block text-xs font-medium text-zinc-600">Description</label>
              <template v-if="editingBody">
                <textarea
                  ref="bodyInput"
                  v-model="bodyDraft"
                  rows="10"
                  class="w-full resize-none rounded-lg bg-zinc-800 px-3 py-2 text-sm text-zinc-100 outline-none ring-1 ring-zinc-600 transition focus:ring-zinc-400 mb-2"
                  placeholder="Write markdown here…"
                  @keydown.esc.stop="cancelBody"
                />
                <div class="flex gap-2">
                  <AppButton variant="secondary" :disabled="!bodyDirty" @click="saveBody">Save</AppButton>
                  <AppButton variant="ghost" @click="cancelBody">Cancel</AppButton>
                </div>
              </template>
              <div
                v-else-if="parsedBody"
                class="group relative cursor-text"
                @click="startEditBody"
              >
                <div
                  class="rounded-lg px-1 py-1 text-sm text-zinc-300 leading-relaxed
                    [&_h1]:text-lg [&_h1]:font-bold [&_h1]:text-zinc-100 [&_h1]:mb-2
                    [&_h2]:text-base [&_h2]:font-semibold [&_h2]:text-zinc-100 [&_h2]:mb-1.5
                    [&_strong]:font-semibold [&_strong]:text-zinc-100
                    [&_code]:bg-zinc-800 [&_code]:rounded [&_code]:px-1 [&_code]:py-0.5 [&_code]:text-xs [&_code]:font-mono
                    [&_pre]:bg-zinc-800 [&_pre]:rounded-lg [&_pre]:p-3 [&_pre]:overflow-x-auto [&_pre]:mb-2
                    [&_a]:text-blue-400 [&_a]:underline
                    [&_ul]:list-disc [&_ul]:pl-5 [&_ul]:my-1
                    [&_ol]:list-decimal [&_ol]:pl-5 [&_ol]:my-1
                    [&_li]:my-0.5
                    [&_p]:mb-2 [&_p:last-child]:mb-0
                    [&_blockquote]:border-l-2 [&_blockquote]:border-zinc-600 [&_blockquote]:pl-3 [&_blockquote]:text-zinc-400"
                  v-html="parsedBody"
                />
                <svg xmlns="http://www.w3.org/2000/svg" class="absolute top-1 right-1 h-4 w-4 text-zinc-600 opacity-0 group-hover:opacity-100 transition" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Z" />
                </svg>
              </div>
              <button
                v-else
                class="cursor-pointer w-full rounded-lg px-3 py-2 text-left text-sm text-zinc-600 hover:bg-zinc-800 hover:text-zinc-400 transition"
                @click="startEditBody"
              >
                Add description…
              </button>
            </div>

            <!-- URLs -->
            <div class="mb-6">
              <div class="mb-2 flex items-center justify-between">
                <label class="text-xs font-medium text-zinc-600">Links</label>
                <button
                  v-if="!editingUrls"
                  class="cursor-pointer text-xs text-zinc-600 transition hover:text-zinc-400"
                  @click="startEditUrls"
                >
                  {{ local.urls.length > 0 ? 'Edit' : '+ Add link' }}
                </button>
              </div>

              <template v-if="editingUrls">
                <div class="flex flex-col gap-2 mb-2">
                  <div v-for="(_, i) in urlInputs" :key="i" class="flex items-center gap-2">
                    <input
                      v-model="urlInputs[i]"
                      type="url"
                      placeholder="https://…"
                      class="min-w-0 flex-1 rounded-lg bg-zinc-800 px-3 py-1.5 text-sm text-zinc-100 placeholder-zinc-600 outline-none ring-1 ring-zinc-700 transition focus:ring-zinc-500"
                    />
                    <button class="cursor-pointer shrink-0 text-zinc-600 hover:text-zinc-400 transition" @click="removeUrl(i)">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>
                </div>
                <div class="flex gap-2">
                  <button class="cursor-pointer text-xs text-zinc-500 hover:text-zinc-300 transition" @click="addUrl">+ Add link</button>
                  <AppButton variant="secondary" @click="saveUrls">Save</AppButton>
                  <AppButton variant="ghost" @click="cancelUrls">Cancel</AppButton>
                </div>
              </template>

              <div v-else-if="local.urls.length > 0" class="flex flex-col gap-1.5">
                <a
                  v-for="url in local.urls"
                  :key="url"
                  :href="url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex items-center gap-2 text-xs text-zinc-500 transition hover:text-zinc-300"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                  </svg>
                  <span class="truncate">{{ url }}</span>
                </a>
              </div>
            </div>

            <!-- Close / Reopen -->
            <div class="border-t border-zinc-800 pt-4">
              <AppButton variant="outline" class="w-full justify-center" @click="requestToggleClosed">
                {{ isTerminal ? 'Reopen ticket' : 'Close ticket' }}
              </AppButton>
            </div>
          </div>
        </div>

        <!-- Confirm dialog (inside Teleport scope, above drawer) -->
        <Transition name="confirm">
          <div
            v-if="confirm"
            class="absolute inset-0 z-60 flex items-center justify-center p-4"
            @click.self="cancelConfirm"
          >
            <div class="w-full max-w-xs rounded-2xl bg-zinc-900 p-5 shadow-2xl ring-1 ring-white/10">
              <p class="mb-4 text-sm font-medium text-zinc-100">{{ confirm.message }}</p>
              <div class="flex justify-end gap-2">
                <AppButton variant="ghost" @click="cancelConfirm">Cancel</AppButton>
                <AppButton variant="secondary" @click="acceptConfirm">Confirm</AppButton>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.15s ease;
}
.drawer-enter-active .relative,
.drawer-leave-active .relative {
  transition: transform 0.2s ease;
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}
.drawer-enter-from .relative,
.drawer-leave-to .relative {
  transform: translateX(100%);
}

.confirm-enter-active,
.confirm-leave-active {
  transition: opacity 0.1s ease, transform 0.1s ease;
}
.confirm-enter-from,
.confirm-leave-to {
  opacity: 0;
  transform: scale(0.96);
}
</style>
