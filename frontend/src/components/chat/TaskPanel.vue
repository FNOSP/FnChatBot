<script setup lang="ts">
import { computed } from 'vue'
import { CheckCircle2, Circle, Loader2, ArrowRight } from 'lucide-vue-next'

interface Task {
  content: string
  status: 'pending' | 'in_progress' | 'completed'
  activeForm?: string
}

const props = defineProps<{
  tasks: Task[]
}>()

const completedCount = computed(() => props.tasks.filter(t => t.status === 'completed').length)
const totalCount = computed(() => props.tasks.length)
const progress = computed(() => totalCount.value ? (completedCount.value / totalCount.value) * 100 : 0)
</script>

<template>
  <div class="bg-bg-card border border-border rounded-lg overflow-hidden shadow-sm">
    <div class="p-3 bg-bg-secondary border-b border-border flex justify-between items-center">
      <h3 class="font-medium text-sm">Plan &amp; Progress</h3>
      <span class="text-xs text-text-secondary">{{ completedCount }}/{{ totalCount }}</span>
    </div>
    
    <!-- Progress Bar -->
    <div class="h-1.5 bg-muted w-full">
      <div class="h-full bg-primary transition-all duration-500" :style="{ width: `${progress}%` }"></div>
    </div>

    <div class="p-3 space-y-3 max-h-[320px] overflow-y-auto">
      <div v-if="tasks.length === 0" class="text-xs text-text-muted text-center py-4">
        No plan formulated yet.
      </div>
      
      <div 
        v-for="(task, index) in tasks" 
        :key="index"
        class="flex items-start gap-2 text-sm group"
        :class="{
          'opacity-50': task.status === 'pending',
          'text-primary': task.status === 'in_progress'
        }"
      >
        <!-- Icon -->
        <div class="mt-0.5 shrink-0">
          <CheckCircle2 v-if="task.status === 'completed'" class="w-4 h-4 text-green-500" />
          <Loader2 v-else-if="task.status === 'in_progress'" class="w-4 h-4 animate-spin text-primary" />
          <Circle v-else class="w-4 h-4 text-muted-foreground" />
        </div>
        
        <!-- Content -->
        <div class="flex-1 min-w-0">
          <div :class="{'line-through text-muted-foreground': task.status === 'completed'}">
            {{ task.content }}
          </div>
          
          <!-- Active Form (Current Action) -->
          <div v-if="task.status === 'in_progress' && task.activeForm" class="text-xs text-muted-foreground mt-1 flex items-center gap-1">
            <ArrowRight class="w-3 h-3" />
            {{ task.activeForm }}...
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
