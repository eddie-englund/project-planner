<script setup lang="ts">
import type { Topic } from '@/types/topic'
import AppButton from '@/components/AppButton.vue'

defineProps<{ topic: Topic }>()
defineEmits<{ click: []; edit: []; delete: [] }>()
</script>

<template>
  <div
    class="group relative flex aspect-[4/3] w-full cursor-pointer flex-col items-center justify-center overflow-hidden rounded-xl p-4 shadow-lg transition-all duration-200 hover:scale-[1.03] hover:shadow-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40"
    :style="{ backgroundColor: topic.color }"
    role="button"
    tabindex="0"
    @click="$emit('click')"
    @keydown.enter.space.prevent="$emit('click')"
  >
    <div class="pointer-events-none absolute inset-x-0 top-0 h-px bg-white/20" />
    <span class="relative text-center text-sm font-medium leading-snug text-white/90 drop-shadow-sm">
      {{ topic.title }}
    </span>
    <div class="absolute right-2 top-2 flex gap-1 opacity-0 transition-opacity duration-150 group-hover:opacity-100">
      <AppButton variant="icon" @click.stop="$emit('edit')">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125" />
        </svg>
      </AppButton>
      <AppButton variant="icon" @click.stop="$emit('delete')">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
        </svg>
      </AppButton>
    </div>
  </div>
</template>
