<script setup lang="ts">
import type { Ticket, Status } from '@/types/ticket'

const props = defineProps<{
  ticket: Ticket
  status?: Status
  muted?: boolean
  topicColor?: string
  topicTitle?: string
}>()

defineEmits<{ click: [] }>()
</script>

<template>
  <div
    :data-ticket-id="ticket.id"
    :class="[
      'rounded-xl bg-zinc-900 p-4 ring-1 ring-zinc-800 transition cursor-grab active:cursor-grabbing select-none',
      'hover:ring-zinc-700',
      muted ? 'opacity-50' : '',
    ]"
    @click="$emit('click')"
  >
    <div class="mb-1.5 flex items-start justify-between gap-2">
      <p
        :class="[
          'text-sm font-semibold leading-snug',
          muted ? 'text-zinc-500 line-through' : 'text-zinc-100',
        ]"
      >
        {{ ticket.title }}
      </p>
      <span
        v-if="status"
        class="shrink-0 rounded-full px-2 py-0.5 text-xs font-medium text-white/90"
        :style="{ backgroundColor: status.color }"
      >
        {{ status.name }}
      </span>
    </div>

    <p v-if="ticket.body" class="mb-3 line-clamp-2 text-xs leading-relaxed text-zinc-500">
      {{ ticket.body }}
    </p>

    <div v-if="ticket.urls.length > 0" class="flex flex-col gap-1">
      <a
        v-for="url in ticket.urls"
        :key="url"
        :href="url"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-1.5 text-xs text-zinc-600 transition hover:text-zinc-400"
        @click.stop
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
        </svg>
        <span class="truncate">{{ url }}</span>
      </a>
    </div>

    <div v-if="topicColor && topicTitle" class="mt-2 flex items-center gap-1.5">
      <span class="h-2 w-2 shrink-0 rounded-full" :style="{ backgroundColor: topicColor }" />
      <span class="truncate text-xs text-zinc-500">{{ topicTitle }}</span>
    </div>
  </div>
</template>
