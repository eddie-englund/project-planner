<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useTopicsStore } from '@/stores/topics'
import { useTitle } from '@vueuse/core'
import { useProjectsStore } from '@/stores/projects'
import TopicCard from '@/components/TopicCard.vue'
import NewProjectCard from '@/components/NewProjectCard.vue'
import NewTopicModal from '@/components/NewTopicModal.vue'
import AppButton from '@/components/AppButton.vue'

const route = useRoute('/projects/[id]/')
const router = useRouter()
const projectId = route.params.id

const topicsStore = useTopicsStore()
const projectsStore = useProjectsStore()
const { filtered, search, items: topics, isLoading, error } = storeToRefs(topicsStore)

const project = computed(() => projectsStore.items.find((p) => p.id === projectId))
const nextIndex = computed(() => topics.value.length)
const showModal = ref(false)

useTitle(computed(() => project.value ? `${project.value.title} — Topics` : 'Project'))

onMounted(async () => {
  topicsStore.search = ''
  if (projectsStore.items.length === 0) await projectsStore.fetchAll()
  await topicsStore.fetchAll(projectId)
})
</script>

<template>
  <div class="min-h-screen bg-zinc-950">
    <!-- Sticky header -->
    <div class="sticky top-0 z-10 bg-zinc-950/95 px-6 pb-4 pt-6 backdrop-blur-sm sm:px-10 sm:pt-10">
      <div class="mb-4 flex items-center gap-3">
        <RouterLink
          to="/"
          class="flex items-center justify-center rounded-lg p-1.5 text-zinc-500 transition hover:bg-zinc-800 hover:text-zinc-200"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
          </svg>
        </RouterLink>
        <h1 class="flex-1 text-xl font-semibold text-zinc-200">
          {{ project?.title ?? 'Project' }}
        </h1>
        <AppButton variant="secondary" @click="showModal = true">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          New topic
        </AppButton>
      </div>

      <!-- Search bar -->
      <div class="flex items-center gap-2 rounded-xl bg-zinc-900 px-3 py-2 ring-1 ring-zinc-800 transition focus-within:ring-zinc-600">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 shrink-0 text-zinc-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 1 1 5 11a6 6 0 0 1 12 0z" />
        </svg>
        <input
          v-model="search"
          type="text"
          placeholder="Filter topics…"
          class="w-full bg-transparent text-sm text-zinc-100 placeholder-zinc-500 outline-none"
        />
        <button
          v-if="search"
          class="shrink-0 text-zinc-500 transition hover:text-zinc-300"
          @click="topicsStore.search = ''"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="px-6 pb-10 sm:px-10">
      <p v-if="isLoading" class="text-sm text-zinc-500">Loading…</p>
      <p v-else-if="error" class="text-sm font-medium text-red-400">{{ error }}</p>

      <template v-else>
        <!-- Empty search result -->
        <p v-if="filtered.length === 0 && search" class="py-16 text-center text-sm text-zinc-500">
          No topics match "{{ search }}"
        </p>

        <!-- No topics at all -->
        <div v-else-if="filtered.length === 0" class="flex flex-col items-center justify-center py-24 text-center">
          <p class="mb-1 text-sm font-medium text-zinc-400">No topics yet</p>
          <p class="mb-6 text-xs text-zinc-600">Break this project into topics to get started</p>
          <AppButton variant="outline" @click="showModal = true">Add first topic</AppButton>
        </div>

        <!-- Grid -->
        <div v-else class="grid grid-cols-3 gap-4 sm:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
          <TopicCard
            v-for="(topic, i) in filtered"
            :key="topic.id"
            :topic="topic"
            class="animate-fade-up"
            :style="{ '--delay': `${Math.min(i * 40, 300)}ms` }"
            @click="router.push(`/projects/${projectId}/topics/${topic.id}`)"
          />
          <NewProjectCard v-if="!search" aspect="landscape" @click="showModal = true" />
        </div>
      </template>
    </div>

    <NewTopicModal
      :open="showModal"
      :project-id="projectId"
      :next-index="nextIndex"
      @close="showModal = false"
      @created="showModal = false"
    />
  </div>
</template>
