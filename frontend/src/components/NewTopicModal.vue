<script setup lang="ts">
import { ref, watch } from "vue";
import { useTopicsStore } from "@/stores/topics";
import AppButton from "@/components/AppButton.vue";
import type { Topic } from "@/types/topic";

const props = defineProps<{ open: boolean; projectId: string; nextIndex: number; topic?: Topic | null }>();
const emit = defineEmits<{ close: []; created: []; updated: [] }>();

const store = useTopicsStore();

const PALETTE = [
  '#C28B82',
  '#82B0C2',
  '#A38EC8',
  '#82C2A8',
  '#C2A882',
  '#C282A8',
  '#A8C282',
  '#82B8C2',
];

const title = ref("");
const selectedColor = ref<string>(PALETTE[0] ?? "#7C6F8E");
const submitting = ref(false);

watch(() => props.open, (isOpen) => {
  if (isOpen && props.topic) {
    title.value = props.topic.title;
    selectedColor.value = props.topic.color;
  } else if (isOpen) {
    title.value = "";
    selectedColor.value = PALETTE[0] ?? "#7C6F8E";
  }
});

async function submit() {
  if (!title.value.trim()) return;
  submitting.value = true;
  if (props.topic) {
    await store.update(props.projectId, props.topic.id, {
      title: title.value.trim(),
      color: selectedColor.value,
    });
    submitting.value = false;
    emit("updated");
    emit("close");
  } else {
    await store.create(props.projectId, {
      title: title.value.trim(),
      color: selectedColor.value,
      index: props.nextIndex,
    });
    submitting.value = false;
    emit("created");
    emit("close");
  }
}

function close() {
  title.value = "";
  selectedColor.value = PALETTE[0] ?? "#7C6F8E";
  emit("close");
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
        @click.self="close"
      >
        <div class="w-full max-w-sm rounded-2xl bg-zinc-900 p-6 shadow-2xl ring-1 ring-white/10">
          <h2 class="mb-5 text-base font-semibold text-zinc-100">
            {{ topic ? 'Edit topic' : 'New topic' }}
          </h2>

          <div class="mb-4">
            <label class="mb-1.5 block text-xs font-medium text-zinc-400">Title</label>
            <input
              v-model="title"
              type="text"
              placeholder="Topic name"
              class="w-full rounded-lg bg-zinc-800 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-500 outline-none ring-1 ring-zinc-700 transition focus:ring-zinc-500"
              @keydown.enter="submit"
              @keydown.esc="close"
            />
          </div>

          <div class="mb-6">
            <label class="mb-2 block text-xs font-medium text-zinc-400">Color</label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="color in PALETTE"
                :key="color"
                type="button"
                class="h-7 w-7 rounded-full transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40 cursor-pointer"
                :class="
                  selectedColor === color ? 'scale-110 ring-2 ring-white/60' : 'hover:scale-105'
                "
                :style="{ backgroundColor: color }"
                @click="selectedColor = color"
              />
            </div>
          </div>

          <div class="flex justify-end gap-2">
            <AppButton variant="ghost" @click="close">Cancel</AppButton>
            <AppButton
              variant="primary"
              :color="selectedColor"
              :disabled="submitting || !title.trim()"
              @click="submit"
            >
              {{ submitting ? (topic ? 'Saving…' : 'Creating…') : (topic ? 'Save' : 'Create') }}
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
