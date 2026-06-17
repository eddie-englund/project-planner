<script setup lang="ts">
withDefaults(defineProps<{ aspect?: 'square' | 'landscape' }>(), { aspect: 'square' })
defineEmits<{ click: [] }>()
</script>

<template>
  <button
    :class="[
      'flex w-full cursor-pointer items-center justify-center border-2 border-dashed border-zinc-600 bg-zinc-800/50 transition-all duration-200 hover:border-zinc-400 hover:bg-zinc-700/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40',
      aspect === 'landscape' ? 'aspect-[4/3] rounded-xl' : 'aspect-square rounded-2xl',
    ]"
    @click="$emit('click')"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-8 w-8 text-zinc-400"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      stroke-width="1.5"
    >
      <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
    </svg>
  </button>
</template>
