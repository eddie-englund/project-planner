<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useProjectsStore } from "@/stores/projects";
import { useTitle } from "@vueuse/core";

useTitle('Projects');
import ProjectCard from "@/components/ProjectCard.vue";
import NewProjectCard from "@/components/NewProjectCard.vue";
import NewProjectModal from "@/components/NewProjectModal.vue";
import AppButton from "@/components/AppButton.vue";

const router = useRouter();
const store = useProjectsStore();
const { filtered, search, isLoading, error } = storeToRefs(store);

const showModal = ref(false);

onMounted(() => {
  store.search = ''
  store.fetchAll()
})
</script>

<template>
  <div class="min-h-screen bg-zinc-950">
    <!-- Sticky header -->
    <div class="sticky top-0 z-10 bg-zinc-950/95 px-6 pb-4 pt-6 backdrop-blur-sm sm:px-10 sm:pt-10">
      <div class="mb-4 flex items-center justify-between">
        <h1 class="text-xl font-semibold text-zinc-200">Projects</h1>
        <AppButton variant="secondary" @click="showModal = true">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          New project
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
          placeholder="Filter projects…"
          class="w-full bg-transparent text-sm text-zinc-100 placeholder-zinc-500 outline-none"
        />
        <button
          v-if="search"
          class="shrink-0 text-zinc-500 transition hover:text-zinc-300"
          @click="store.search = ''"
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
      <p v-else-if="error" class="text-sm text-red-400">{{ error }}</p>

      <template v-else>
        <!-- Empty search result -->
        <p v-if="filtered.length === 0 && search" class="py-16 text-center text-sm text-zinc-500">
          No projects match "{{ search }}"
        </p>

        <!-- No projects at all -->
        <div v-else-if="filtered.length === 0" class="flex flex-col items-center justify-center py-24 text-center">
          <p class="mb-1 text-sm font-medium text-zinc-400">No projects yet</p>
          <p class="mb-6 text-xs text-zinc-600">Create your first project to get started</p>
          <AppButton variant="outline" @click="showModal = true">Add first project</AppButton>
        </div>

        <!-- Grid -->
        <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          <ProjectCard
            v-for="(project, i) in filtered"
            :key="project.id"
            :project="project"
            class="animate-fade-up"
            :style="{ '--delay': `${Math.min(i * 40, 300)}ms` }"
            @click="router.push(`/projects/${project.id}`)"
          />
          <NewProjectCard v-if="!search" @click="showModal = true" />
        </div>
      </template>
    </div>

    <NewProjectModal :open="showModal" @close="showModal = false" @created="showModal = false" />
  </div>
</template>
