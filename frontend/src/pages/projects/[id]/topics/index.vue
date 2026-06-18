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
import NewProjectModal from '@/components/NewProjectModal.vue'
import AppButton from '@/components/AppButton.vue'
import type { Topic } from '@/types/topic'
import type { Project } from '@/types/project'

const route = useRoute('/projects/[id]/topics/')
const router = useRouter()
const projectId = route.params.id

const topicsStore = useTopicsStore()
const projectsStore = useProjectsStore()
const { filtered, search, items: topics, isLoading, error } = storeToRefs(topicsStore)

const project = computed(() => projectsStore.items.find((p) => p.id === projectId))
const nextIndex = computed(() => topics.value.length)
const showModal = ref(false)
const editingTopic = ref<Topic | null>(null)
const confirmDelete = ref<{ topic: Topic; deleting: boolean } | null>(null)
const editingProject = ref<Project | null>(null)
const confirmDeleteProject = ref<{ deleting: boolean } | null>(null)

useTitle(computed(() => project.value ? `${project.value.title} — Topics` : 'Project'))

onMounted(async () => {
  topicsStore.search = ''
  if (projectsStore.items.length === 0) await projectsStore.fetchAll()
  await topicsStore.fetchAll(projectId)
})

function openEdit(topic: Topic) {
  editingTopic.value = topic
  showModal.value = false
}

function openCreate() {
  editingTopic.value = null
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingTopic.value = null
}

async function confirmDeleteTopic() {
  if (!confirmDelete.value) return
  confirmDelete.value.deleting = true
  await topicsStore.remove(projectId, confirmDelete.value.topic.id)
  confirmDelete.value = null
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
</script>

<template>
  <div class="min-h-screen bg-zinc-950">
    <!-- Sticky header -->
    <div class="sticky top-0 z-10 bg-zinc-950/95 px-6 pb-4 pt-6 backdrop-blur-sm sm:px-10 sm:pt-10">
      <div class="mb-4 flex items-center gap-3">
        <RouterLink
          to="/"
          class="flex cursor-pointer items-center justify-center rounded-lg p-1.5 text-zinc-500 transition hover:bg-zinc-800 hover:text-zinc-200"
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
        <AppButton variant="secondary" @click="openCreate">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          New topic
        </AppButton>
        <!-- View toggle -->
        <div class="flex rounded-xl bg-zinc-900 p-0.5 ring-1 ring-zinc-800 text-xs font-medium">
          <span class="rounded-lg px-3 py-1.5 bg-zinc-800 text-zinc-200">Topics</span>
          <RouterLink
            :to="`/projects/${projectId}`"
            class="cursor-pointer rounded-lg px-3 py-1.5 text-zinc-500 hover:text-zinc-200 transition"
          >
            Kanban
          </RouterLink>
        </div>
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
          class="shrink-0 cursor-pointer text-zinc-500 transition hover:text-zinc-300"
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
          <AppButton variant="outline" @click="openCreate">Add first topic</AppButton>
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
            @edit="openEdit(topic)"
            @delete="confirmDelete = { topic, deleting: false }"
          />
          <NewProjectCard v-if="!search" aspect="landscape" @click="openCreate" />
        </div>
      </template>
    </div>

    <NewTopicModal
      :open="showModal || !!editingTopic"
      :project-id="projectId"
      :next-index="nextIndex"
      :topic="editingTopic"
      @close="closeModal"
      @created="closeModal"
      @updated="closeModal"
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

    <!-- Delete topic confirmation -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="confirmDelete"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm"
          @click.self="confirmDelete = null"
        >
          <div class="w-full max-w-xs rounded-2xl bg-zinc-900 p-6 shadow-2xl ring-1 ring-white/10">
            <h2 class="mb-2 text-base font-semibold text-zinc-100">Delete topic?</h2>
            <p class="mb-6 text-sm text-zinc-400">
              "{{ confirmDelete.topic.title }}" will be permanently deleted.
            </p>
            <div class="flex justify-end gap-2">
              <AppButton variant="ghost" :disabled="confirmDelete.deleting" @click="confirmDelete = null">
                Cancel
              </AppButton>
              <AppButton
                variant="primary"
                color="#dc2626"
                :disabled="confirmDelete.deleting"
                @click="confirmDeleteTopic"
              >
                {{ confirmDelete.deleting ? 'Deleting…' : 'Delete' }}
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
