<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'ghost' | 'secondary' | 'outline' | 'icon' | 'icon-ghost' | 'icon-danger'
  color?: string
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'secondary',
  type: 'button',
})

defineEmits<{ click: [MouseEvent] }>()

const base = 'cursor-pointer transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 disabled:cursor-not-allowed'

const variantClasses: Record<NonNullable<Props['variant']>, string> = {
  primary: 'rounded-lg px-4 py-2 text-sm font-medium text-white disabled:opacity-50',
  ghost: 'rounded-lg px-4 py-2 text-sm text-zinc-400 hover:text-zinc-200',
  secondary: 'flex items-center gap-1.5 rounded-xl bg-zinc-800 px-3 py-2 text-sm font-medium text-zinc-200 hover:bg-zinc-700',
  outline: 'rounded-xl border border-zinc-700 bg-zinc-800/60 px-5 py-2.5 text-sm font-medium text-zinc-300 hover:border-zinc-500 hover:bg-zinc-700/60',
  icon: 'rounded-md p-1 text-white/80 hover:text-white hover:bg-black/20 disabled:opacity-40',
  'icon-ghost': 'flex items-center justify-center rounded-lg p-1.5 text-zinc-500 hover:bg-zinc-800 hover:text-zinc-200 disabled:opacity-40',
  'icon-danger': 'flex items-center justify-center rounded-lg p-1.5 text-zinc-500 hover:bg-red-950 hover:text-red-400 disabled:opacity-40',
}
</script>

<template>
  <button
    :type="type"
    :disabled="disabled"
    :class="[base, variantClasses[variant ?? 'secondary']]"
    :style="variant === 'primary' && color ? { backgroundColor: color } : undefined"
    @click="$emit('click', $event)"
  >
    <slot />
  </button>
</template>
