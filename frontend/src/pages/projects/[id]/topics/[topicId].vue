<script setup lang="ts">
import { ref, computed, onMounted, watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useTitle } from '@vueuse/core'
import { VueDraggable } from 'vue-draggable-plus'
import { useProjectsStore } from '@/stores/projects'
import { useTopicsStore } from '@/stores/topics'
import { useTicketsStore } from '@/stores/tickets'
import { useStatusesStore } from '@/stores/statuses'
import TicketCard from '@/components/TicketCard.vue'
import TicketDrawer from '@/components/TicketDrawer.vue'
import NewTicketModal from '@/components/NewTicketModal.vue'
import AppButton from '@/components/AppButton.vue'
import type { Ticket, Status } from '@/types/ticket'

const route = useRoute('/projects/[id]/topics/[topicId]')
const projectId = route.params.id
const topicId = route.params.topicId

const projectsStore = useProjectsStore()
const topicsStore = useTopicsStore()
const ticketsStore = useTicketsStore()
const statusesStore = useStatusesStore()

const { items: tickets, isLoading, error, search, sort } = storeToRefs(ticketsStore)
const { items: statuses } = storeToRefs(statusesStore)

const project = computed(() => projectsStore.items.find((p) => p.id === projectId))
const topic = computed(() => topicsStore.items.find((t) => t.id === topicId))

const showModal = ref(false)
const drawerTicket = ref<Ticket | null>(null)

useTitle(computed(() => topic.value?.title ?? 'Topic'))

// Kanban columns: each status + uncategorized
const allColumns = computed(() => {
  const cols: Array<{ status: Status | null; id: string | null }> = statuses.value.map((s) => ({
    status: s as Status | null,
    id: s.id as string | null,
  }))
  const hasUncategorized = tickets.value.some((t) => t.statusId === null)
  if (hasUncategorized) cols.push({ status: null, id: null })
  return cols
})

// Per-column ticket arrays for VueDraggable v-model
// These are reactive local copies; we sync from store, and update store on drop
const columnTickets = ref<Record<string, Ticket[]>>({})

watchEffect(() => {
  const next: Record<string, Ticket[]> = {}
  for (const col of allColumns.value) {
    const key = col.id ?? 'null'
    const filtered = ticketsStore.ticketsForStatus(col.id)
    next[key] = [...filtered.value]
  }
  columnTickets.value = next
})

function onDragEnd(event: { item: HTMLElement; to: HTMLElement; from: HTMLElement }) {
  if (event.to === event.from) return
  const ticketId = event.item.dataset['ticketId']
  if (!ticketId) return
  const raw = event.to.dataset['statusId']
  const newStatusId = raw === 'null' ? null : (raw ?? null)
  ticketsStore.updateStatus(projectId, topicId, ticketId, newStatusId)
}

onMounted(async () => {
  search.value = ''
  if (projectsStore.items.length === 0) await projectsStore.fetchAll()
  if (topicsStore.items.length === 0) await topicsStore.fetchAll(projectId)
  await Promise.all([
    statusesStore.fetchAll(projectId),
    ticketsStore.fetchAll(projectId, topicId),
  ])
})
</script>

<template>
  <div class="flex min-h-screen flex-col bg-zinc-950">
    <!-- Topic color stripe -->
    <div class="h-1 w-full shrink-0" :style="{ backgroundColor: topic?.color ?? '#52525b' }" />

    <!-- Sticky header -->
    <div class="sticky top-1 z-10 bg-zinc-950/95 px-6 pb-3 pt-5 backdrop-blur-sm sm:px-8">
      <div class="flex flex-wrap items-center gap-2">
        <!-- Back -->
        <RouterLink
          :to="`/projects/${projectId}`"
          class="flex items-center justify-center rounded-lg p-1.5 text-zinc-500 transition hover:bg-zinc-800 hover:text-zinc-200"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
          </svg>
        </RouterLink>

        <!-- Title -->
        <div class="min-w-0 flex-1">
          <p v-if="project" class="text-xs text-zinc-600">{{ project.title }}</p>
          <h1 class="truncate text-lg font-bold text-zinc-100">{{ topic?.title ?? 'Topic' }}</h1>
        </div>

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
          <button v-if="search" class="shrink-0 text-zinc-500 hover:text-zinc-300 transition" @click="search = ''">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Sort -->
        <select
          v-model="sort"
          class="rounded-xl bg-zinc-900 px-3 py-1.5 text-xs text-zinc-300 outline-none ring-1 ring-zinc-800 transition hover:ring-zinc-600 cursor-pointer"
        >
          <option value="newest">Newest first</option>
          <option value="oldest">Oldest first</option>
          <option value="alpha">A–Z</option>
        </select>

        <!-- New ticket -->
        <AppButton variant="secondary" @click="showModal = true">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          New ticket
        </AppButton>
      </div>
    </div>

    <!-- Loading / Error -->
    <div v-if="isLoading" class="px-8 py-10 text-sm text-zinc-500">Loading…</div>
    <div v-else-if="error" class="px-8 py-10 text-sm font-medium text-red-400">{{ error }}</div>

    <!-- Empty state (no tickets at all) -->
    <template v-else-if="tickets.length === 0">
      <div class="flex flex-1 flex-col items-center justify-center py-24 text-center">
        <div
          class="mb-4 flex h-12 w-12 items-center justify-center rounded-2xl"
          :style="{ backgroundColor: `${topic?.color ?? '#52525b'}22` }"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5" :style="{ color: topic?.color ?? '#71717a' }">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 0 0 2.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 0 0-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 0 0 .75-.75 2.25 2.25 0 0 0-.1-.664m-5.8 0A2.251 2.251 0 0 1 13.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25ZM6.75 12h.008v.008H6.75V12Zm0 3h.008v.008H6.75V15Zm0 3h.008v.008H6.75V18Z" />
          </svg>
        </div>
        <p class="mb-1 text-sm font-medium text-zinc-400">No tickets yet</p>
        <p class="mb-6 text-xs text-zinc-600">Track work, ideas, and links for this topic</p>
        <AppButton variant="outline" @click="showModal = true">Add first ticket</AppButton>
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
            @update:model-value="(v: Ticket[]) => (columnTickets[col.id ?? 'null'] = v)"
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
              :status="statuses.find((s) => s.id === ticket.statusId)"
              :muted="col.status?.isTerminal ?? false"
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
      :topic-id="topicId"
      @close="drawerTicket = null"
    />

    <!-- New ticket modal -->
    <NewTicketModal
      :open="showModal"
      :project-id="projectId"
      :topic-id="topicId"
      :statuses="statuses"
      @close="showModal = false"
      @created="showModal = false"
    />
  </div>
</template>
