<script setup lang="ts">
import { ref } from 'vue'
import { ChevronDownIcon, ChevronRightIcon, LoadingIcon, CheckIcon } from 'tdesign-icons-vue-next'

const props = defineProps<{
  thinking: string
  tasks: any[]
}>()

const isExpanded = ref(false)
</script>

<template>
  <div class="border border-border rounded-lg bg-bg-card mb-3 overflow-hidden shadow-sm">
    <div 
      @click="isExpanded = !isExpanded"
      class="flex items-center gap-2 px-3 py-2 bg-bg-secondary cursor-pointer text-xs md:text-sm font-medium hover:bg-bg-hover transition-colors"
    >
      <component :is="isExpanded ? ChevronDownIcon : ChevronRightIcon" class="w-4 h-4" />
      <span>Thinking Process</span>
      <span v-if="tasks.length > 0" class="ml-auto text-xs text-text-secondary">
        {{ tasks.filter(t => t.status === 'completed').length }}/{{ tasks.length }} tasks
      </span>
    </div>

    <div v-if="isExpanded" class="p-3 border-t border-border bg-bg-primary text-xs md:text-sm space-y-3">
      <!-- Tasks List -->
      <div v-if="tasks.length > 0" class="space-y-2">
        <h4 class="font-semibold text-xs uppercase text-text-muted">Plan</h4>
        <div v-for="(task, index) in tasks" :key="index" class="flex items-center gap-2">
          <div 
            class="w-4 h-4 rounded-full flex items-center justify-center border"
            :class="{
              'bg-success border-success': task.status === 'completed',
              'bg-warning border-warning': task.status === 'running',
              'border-border': task.status === 'pending'
            }"
          >
            <LoadingIcon v-if="task.status === 'running'" class="w-3 h-3 animate-spin text-white" />
            <CheckIcon v-else-if="task.status === 'completed'" class="text-white text-xs" />
          </div>
          <span :class="{'text-text-muted': task.status === 'pending', 'text-text-primary': task.status !== 'pending'}">
            {{ task.description || task.name }}
          </span>
        </div>
      </div>

      <!-- Thinking Log -->
      <div v-if="thinking">
        <h4 class="font-semibold text-xs uppercase text-text-muted mb-1">Log</h4>
        <p class="font-mono text-xs text-text-secondary whitespace-pre-wrap">{{ thinking }}</p>
      </div>
    </div>
  </div>
</template>

