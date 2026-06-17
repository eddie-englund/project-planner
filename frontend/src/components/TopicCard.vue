<script setup lang="ts">
import type { Topic } from '@/types/topic'

defineProps<{ topic: Topic }>()
defineEmits<{ click: [] }>()
</script>

<template>
  <button
    class="group relative flex aspect-[4/3] w-full cursor-pointer flex-col items-center justify-center overflow-hidden rounded-xl p-4 shadow-lg transition-all duration-200 hover:scale-[1.03] hover:shadow-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40"
    :style="{ backgroundColor: topic.color }"
    @click="$emit('click')"
  >
    <!-- subtle top-edge highlight for depth -->
    <div class="pointer-events-none absolute inset-x-0 top-0 h-px bg-white/20" />
    <span class="relative text-center text-sm font-medium leading-snug text-white/90 drop-shadow-sm">
      {{ topic.title }}
    </span>
  </button>
</template>
