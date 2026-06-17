import { ref, computed } from "vue";
import { defineStore } from "pinia";
import { useApi } from "@/composables/useApi";
import type { Project, CreateProjectPayload } from "@/types/project";

export const useProjectsStore = defineStore("projects", () => {
  const items = ref<Project[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const search = ref('')
  const filtered = computed(() =>
    items.value.filter(p => p.title.toLowerCase().includes(search.value.toLowerCase()))
  )

  async function fetchAll() {
    isLoading.value = true;
    error.value = null;
    const { data, error: fetchError } = await useApi("/projects").json<Project[]>();
    isLoading.value = false;
    if (fetchError.value) {
      error.value = fetchError.value.message;
      return;
    }
    if (data.value) {
      items.value = data.value;
    }
  }

  async function create(payload: CreateProjectPayload) {
    const { data, error: fetchError } = await useApi("/projects").post(payload).json<Project>();
    if (fetchError.value) {
      error.value = fetchError.value.message;
      return;
    }
    if (data.value) {
      items.value.push(data.value);
    }
  }

  return { items, filtered, search, isLoading, error, fetchAll, create };
});
