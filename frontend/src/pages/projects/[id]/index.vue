<script setup lang="ts">
import { ref, computed, onMounted, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useTitle } from '@vueuse/core'
import { VueDraggable } from 'vue-draggable-plus'
import { useProjectsStore } from '@/stores/projects'
import { useTopicsStore } from '@/stores/topics'
import { useTicketsStore } from '@/stores/tickets'
import { useStatusesStore } from '@/stores/statuses'
import TicketCard from '@/components/TicketCard.vue'
import TicketDrawer from '@/components/TicketDrawer.vue'
import NewProjectModal from '@/components/NewProjectModal.vue'
import AppButton from '@/components/AppButton.vue'
import type { TicketWithTopic, Status } from '@/types/ticket'
import type { Project } from '@/types/project'

const route = useRoute('/projects/[id]/')
const router = useRouter()
const projectId = route.params.id

const projectsStore = useProjectsStore()
const topicsStore = useTopicsStore()
const ticketsStore = useTicketsStore()
const statusesStore = useStatusesStore()

const { projectItems, isLoadingProject: isLoading, error } = storeToRefs(ticketsStore)
const { items: statuses } = storeToRefs(statusesStore)
const { items: topics } = storeToRefs(topicsStore)

const project = computed(() => projectsStore.items.find((p) => p.id === projectId))

const search = ref('')
const activeTopicIds = ref<string[]>([])
const drawerTicket = ref<TicketWithTopic | null>(null)
const editingProject = ref<Project | null>(null)
const confirmDeleteProject = ref<{ deleting: boolean } | null>(null)

useTitle(computed(() => project.value ? `${project.value.title} — Kanban` : 'Kanban'))

const filteredTickets = computed(() => {
  let list = projectItems.value
  if (activeTopicIds.value.length > 0)
    list = list.filter((t) => activeTopicIds.value.includes(t.topicId))
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(
      (t) => t.title.toLowerCase().includes(q) || t.body.toLowerCase().includes(q)
    )
  }
  return list
})

const allColumns = computed(() => {
  const sorted = [...statuses.value].sort((a, b) => a.position - b.position)
  const cols: Array<{ status: Status | null; id: string | null }> = sorted.map((s) => ({
    status: s as Status | null,
    id: s.id as string | null,
  }))
  const hasUncategorized = filteredTickets.value.some((t) => t.statusId === null)
  if (hasUncategorized) cols.push({ status: null, id: null })
  return cols
})

const columnTickets = ref<Record<string, TicketWithTopic[]>>({})

watchEffect(() => {
  const next: Record<string, TicketWithTopic[]> = {}
  for (const col of allColumns.value) {
    const key = col.id ?? 'null'
    next[key] = filteredTickets.value.filter((t) => t.statusId === (col.id ?? null))
  }
  columnTickets.value = next
})

function onDragEnd(event: { item: HTMLElement; to: HTMLElement; from: HTMLElement }) {
  if (event.to === event.from) return
  const ticketId = event.item.dataset['ticketId']
  if (!ticketId) return
  const ticket = projectItems.value.find((t) => t.id === ticketId)
  if (!ticket) return
  const raw = event.to.dataset['statusId']
  const newStatusId = raw === 'null' ? null : (raw ?? null)
  ticketsStore.updateStatus(projectId, ticket.topicId, ticketId, newStatusId)
}

function toggleTopic(topicId: string) {
  const idx = activeTopicIds.value.indexOf(topicId)
  if (idx === -1) activeTopicIds.value.push(topicId)
  else activeTopicIds.value.splice(idx, 1)
}

function clearTopicFilter() {
  activeTopicIds.value = []
}

function openEditProject() {
  if (project.value) editingProject.value = project.value
}

function closeProjectModal() {
  editingProject.value = null
}

async function doDeleteProject() {
  if (!confirmDeleteProject.value) return
  confirmDeleteProject.value.deleting = true
  await projectsStore.remove(projectId)
  router.push('/')
}

onMounted(async () => {
  if (projectsStore.items.length === 0) await projectsStore.fetchAll()
  await Promise.all([
    topicsStore.fetchAll(projectId),
    statusesStore.fetchAll(projectId),
    ticketsStore.fetchAllByProject(projectId),
  ])
})
</script>

<template>
  <div class="flex min-h-screen flex-col bg-zinc-950">
    <!-- Sticky header -->
    <div class="sticky top-0 z-10 bg-zinc-950/95 px-6 pb-3 pt-5 backdrop-blur-sm sm:px-8">
      <div class="flex flex-wrap items-center gap-2">
        <!-- Back -->
        <RouterLink
          to="/"
          class="flex cursor-pointer items-center justify-center rounded-lg p-1.5 text-zinc-500 transition hover:bg-zinc-800 hover:text-zinc-200"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
          </svg>
        </RouterLink>

        <!-- Title -->
        <div class="min-w-0 flex-1">
          <h1 class="truncate text-lg font-bold text-zinc-100">{{ project?.title ?? 'Project' }}</h1>
        </div>

        <!-- Edit / Delete project -->
        <AppButton variant="icon-ghost" title="Edit project" @click="openEditProject">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125" />
          </svg>
        </AppButton>
        <AppButton variant="icon-danger" title="Delete project" @click="confirmDeleteProject = { deleting: false }">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
          </svg>
        </AppButton>

        <!-- Search -->
        <div class="flex items-center gap-2 rounded-xl bg-zinc-900 px-3 py-1.5 ring-1 ring-zinc-800 transition focus-within:ring-zinc-600 w-52">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 shrink-0 text-zinc-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 1 1 5 11a6 6 0 0 1 12 0z" />
          </svg>
          <input
            v-model="search"
            type="text"
            placeholder="Search tickets…"
            class="w-full bg-transparent text-xs text-zinc-100 placeholder-zinc-500 outline-none"
          />
          <button v-if="search" class="shrink-0 cursor-pointer text-zinc-500 hover:text-zinc-300 transition" @click="search = ''">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- View toggle -->
        <div class="flex rounded-xl bg-zinc-900 p-0.5 ring-1 ring-zinc-800 text-xs font-medium">
          <span class="rounded-lg px-3 py-1.5 bg-zinc-800 text-zinc-200">Kanban</span>
          <RouterLink
            :to="`/projects/${projectId}/topics`"
            class="cursor-pointer rounded-lg px-3 py-1.5 text-zinc-500 hover:text-zinc-200 transition"
          >
            Topics
          </RouterLink>
        </div>
      </div>

      <!-- Topic filter chips -->
      <div v-if="topics.length > 0" class="mt-3 flex items-center gap-2 overflow-x-auto pb-1">
        <button
          class="shrink-0 cursor-pointer rounded-lg px-3 py-1 text-xs font-medium transition"
          :class="activeTopicIds.length === 0 ? 'bg-zinc-700 text-zinc-100' : 'bg-zinc-900 text-zinc-400 hover:text-zinc-200 ring-1 ring-zinc-800'"
          @click="clearTopicFilter"
        >
          All
        </button>
        <button
          v-for="topic in topics"
          :key="topic.id"
          class="flex shrink-0 cursor-pointer items-center gap-1.5 rounded-lg px-3 py-1 text-xs font-medium transition ring-1"
          :class="activeTopicIds.includes(topic.id)
            ? 'text-zinc-100 ring-transparent'
            : 'bg-zinc-900 text-zinc-400 hover:text-zinc-200 ring-zinc-800'"
          :style="activeTopicIds.includes(topic.id) ? { backgroundColor: topic.color + '33', outlineColor: topic.color } : {}"
          @click="toggleTopic(topic.id)"
        >
          <span class="h-2 w-2 shrink-0 rounded-full" :style="{ backgroundColor: topic.color }" />
          {{ topic.title }}
        </button>
      </div>
    </div>

    <!-- Loading / Error -->
    <div v-if="isLoading" class="px-8 py-10 text-sm text-zinc-500">Loading…</div>
    <div v-else-if="error" class="px-8 py-10 text-sm font-medium text-red-400">{{ error }}</div>

    <!-- Empty state -->
    <template v-else-if="projectItems.length === 0">
      <div class="flex flex-1 flex-col items-center justify-center py-24 text-center">
        <p class="mb-1 text-sm font-medium text-zinc-400">No tickets yet</p>
        <p class="mb-6 text-xs text-zinc-600">
          Add tickets inside a topic to see them here
        </p>
        <RouterLink :to="`/projects/${projectId}/topics`">
          <AppButton variant="outline">Go to Topics</AppButton>
        </RouterLink>
      </div>
    </template>

    <!-- Kanban board -->
    <template v-else>
      <div class="flex flex-1 gap-5 overflow-x-auto px-6 pb-10 pt-4 sm:px-8">
        <div
          v-for="col in allColumns"
          :key="col.id ?? 'null'"
          class="flex w-[300px] shrink-0 flex-col"
        >
          <!-- Column header -->
          <div
            class="mb-3 flex items-center gap-2 border-l-[3px] pl-3"
            :style="col.status ? { borderColor: col.status.color } : { borderColor: '#3f3f46' }"
          >
            <span
              class="h-2 w-2 rounded-full shrink-0"
              :style="{ backgroundColor: col.status?.color ?? '#3f3f46' }"
            />
            <span class="text-xs font-semibold uppercase tracking-widest text-zinc-400">
              {{ col.status?.name ?? 'Uncategorized' }}
            </span>
            <span class="text-xs text-zinc-700">
              {{ columnTickets[col.id ?? 'null']?.length ?? 0 }}
            </span>
          </div>

          <!-- Draggable ticket list -->
          <VueDraggable
            :model-value="columnTickets[col.id ?? 'null'] ?? []"
            @update:model-value="(v: TicketWithTopic[]) => (columnTickets[col.id ?? 'null'] = v)"
            group="tickets"
            :animation="150"
            ghost-class="opacity-20"
            class="flex flex-col gap-2 flex-1 min-h-[60px]"
            :data-status-id="col.id ?? 'null'"
            @end="(e: any) => onDragEnd(e)"
          >
            <TicketCard
              v-for="ticket in columnTickets[col.id ?? 'null']"
              :key="ticket.id"
              :ticket="ticket"
              :muted="col.status?.isTerminal ?? false"
              :topic-color="ticket.topicColor"
              :topic-title="ticket.topicTitle"
              :data-ticket-id="ticket.id"
              @click="drawerTicket = ticket"
            />
          </VueDraggable>

          <!-- Empty column placeholder -->
          <div
            v-if="!columnTickets[col.id ?? 'null']?.length"
            class="rounded-lg border border-dashed border-zinc-800 flex items-center justify-center h-16 text-xs text-zinc-700"
          >
            Drop here
          </div>
        </div>
      </div>
    </template>

    <!-- Ticket drawer -->
    <TicketDrawer
      :ticket="drawerTicket"
      :statuses="statuses"
      :project-id="projectId"
      :topic-id="drawerTicket?.topicId ?? ''"
      @close="drawerTicket = null"
    />

    <NewProjectModal
      :open="!!editingProject"
      :project="editingProject"
      @close="closeProjectModal"
      @updated="closeProjectModal"
    />

    <!-- Delete project confirmation -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="confirmDeleteProject"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
          @click.self="confirmDeleteProject = null"
        >
          <div class="w-full max-w-xs rounded-2xl bg-zinc-900 p-6 shadow-2xl ring-1 ring-white/10">
            <h2 class="mb-2 text-base font-semibold text-zinc-100">Delete project?</h2>
            <p class="mb-6 text-sm text-zinc-400">
              "{{ project?.title }}" and all its topics will be permanently deleted.
            </p>
            <div class="flex justify-end gap-2">
              <AppButton variant="ghost" :disabled="confirmDeleteProject.deleting" @click="confirmDeleteProject = null">
                Cancel
              </AppButton>
              <AppButton
                variant="primary"
                color="#dc2626"
                :disabled="confirmDeleteProject.deleting"
                @click="doDeleteProject"
              >
                {{ confirmDeleteProject.deleting ? 'Deleting…' : 'Delete' }}
              </AppButton>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
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
