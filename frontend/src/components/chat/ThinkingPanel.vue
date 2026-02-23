<script setup lang="ts">
import { ref } from 'vue'
import { ChevronDown, ChevronRight, Loader2 } from 'lucide-vue-next'

const props = defineProps<{
  thinking: string
  tasks: any[]
}>()

const isExpanded = ref(false)
</script>

<template>
  <div class="border rounded-lg bg-card mb-4 overflow-hidden">
    <div 
      @click="isExpanded = !isExpanded"
      class="flex items-center gap-2 p-3 bg-zinc-50 dark:bg-zinc-900 cursor-pointer text-sm font-medium"
    >
      <component :is="isExpanded ? ChevronDown : ChevronRight" class="w-4 h-4" />
      <span>Thinking Process</span>
      <span v-if="tasks.length > 0" class="ml-auto text-xs text-muted-foreground">
        {{ tasks.filter(t => t.status === 'completed').length }}/{{ tasks.length }} tasks
      </span>
    </div>

    <div v-if="isExpanded" class="p-4 border-t border-border bg-background text-sm space-y-4">
      <!-- Tasks List -->
      <div v-if="tasks.length > 0" class="space-y-2">
        <h4 class="font-semibold text-xs uppercase text-muted-foreground">Plan</h4>
        <div v-for="(task, index) in tasks" :key="index" class="flex items-center gap-2">
          <div 
            class="w-4 h-4 rounded-full flex items-center justify-center border"
            :class="{
              'bg-green-500 border-green-500': task.status === 'completed',
              'bg-yellow-500 border-yellow-500': task.status === 'running',
              'border-zinc-300': task.status === 'pending'
            }"
          >
            <Loader2 v-if="task.status === 'running'" class="w-3 h-3 animate-spin text-white" />
            <span v-else-if="task.status === 'completed'" class="text-white text-xs">✓</span>
          </div>
          <span :class="{'text-muted-foreground': task.status === 'pending'}">
            {{ task.description || task.name }}
          </span>
        </div>
      </div>

      <!-- Thinking Log -->
      <div v-if="thinking">
        <h4 class="font-semibold text-xs uppercase text-muted-foreground mb-1">Log</h4>
        <p class="font-mono text-xs text-zinc-600 dark:text-zinc-400 whitespace-pre-wrap">{{ thinking }}</p>
      </div>
    </div>
  </div>
</template>
